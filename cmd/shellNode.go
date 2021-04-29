package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"

	"github.com/get-bridge/truss-cli/truss"
)

var shellNodeCmd = &cobra.Command{
	Use:   "node [node-name-or-instance-id]",
	Short: "Launch a shell on a Truss node via SSH.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		region, err := getAWSRegionFromKubeconfig()
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve AWS region from Kubeconfig")
		}

		awsRole := viper.GetString("aws.assumeRole")
		if awsRole == "" {
			return errors.New("Unable to find aws.assumeRole in config file")
		}

		sess := truss.NewAWSSession(region, awsRole)

		var instanceID string

		if strings.HasPrefix(args[0], "ip-") {
			// Lookup instance-id from Kubernetes if given a node name
			kubeNode, err := describeKubernetesNode(args[0])
			if err != nil {
				return errors.Wrap(err, "Unable to describe Kubernetes node")
			}

			providerIDParts := strings.Split(kubeNode.Spec.ProviderID, "/")
			instanceID = providerIDParts[len(providerIDParts)-1]
		} else if strings.HasPrefix(args[0], "i-") {
			instanceID = args[0]
		}

		instance, err := describeInstance(instanceID, sess)
		if err != nil {
			return errors.Wrap(err, "Unable to describe EC2 instance")
		}

		availabilityZone := *instance.Placement.AvailabilityZone
		hostname := *instance.PrivateDnsName

		publicKey, err := getSSHPublicKey()
		if err != nil {
			return err
		}

		username, err := cmd.Flags().GetString("user")
		if err != nil {
			return err
		}

		sendPublicKey(availabilityZone, instanceID, publicKey, username, sess)

		jump, err := getJump()
		if err != nil {
			return errors.Wrap(err, "Unable to find jumpbox in config file")
		}

		var sshCmd = []string{}

		if len(args) > 1 {
			sshCmd = args[1:]
		}

		return execSSHCommand(hostname, username, jump, sshCmd)
	},
}

func init() {
	shellNodeCmd.Flags().StringP("user", "u", "ec2-user", "The SSH user to target")
	shellCmd.AddCommand(shellNodeCmd)
}

func getAWSRegionFromKubeconfig() (string, error) {
	kubeconfig, err := getKubeconfig()
	if err != nil {
		return "", err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return "", err
	}

	host := config.Host
	hostParts := strings.Split(host, ".")
	region := hostParts[len(hostParts)-4]

	return region, nil
}

func execSSHCommand(hostname string, username string, jump string, sshCmd []string) error {
	sshBinary, err := exec.LookPath("ssh")
	if err != nil {
		return errors.Wrap(err, "Unable to locate ssh binary")
	}

	target := fmt.Sprintf("%s@%s", username, hostname)
	proxyJump := fmt.Sprintf("ProxyJump=%s", jump)

	args := []string{"ssh", "-o", proxyJump, target}
	args = append(args, sshCmd...)

	return syscall.Exec(sshBinary, args, os.Environ())
}

func getJump() (string, error) {
	kubeconfigName, err := getKubeconfigName()
	if err != nil {
		return "", err
	}

	clusterEnv := strings.Replace(kubeconfigName, "kubeconfig-truss-", "", 1)

	jump := viper.GetString("jumps." + clusterEnv)
	if jump == "" {
		return "", errors.New("Could not find jump for " + clusterEnv)
	}

	return jump, nil
}

func describeKubernetesNode(nodeName string) (*v1.Node, error) {
	kubeconfig, err := getKubeconfig()
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
}

func describeInstance(instanceID string, sess *session.Session) (*ec2.Instance, error) {
	ec2svc := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{&instanceID},
	}

	resp, err := ec2svc.DescribeInstances(params)

	if err != nil {
		return nil, err
	}

	return resp.Reservations[0].Instances[0], nil
}

func sendPublicKey(availabilityZone string, instanceID string, publicKey string, instanceUser string, sess *session.Session) error {
	svc := ec2instanceconnect.New(sess)

	params := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(availabilityZone),
		InstanceId:       aws.String(instanceID),
		InstanceOSUser:   aws.String(instanceUser),
		SSHPublicKey:     aws.String(publicKey),
	}

	_, err := svc.SendSSHPublicKey(params)

	if err != nil {
		return errors.Wrap(err, "There was an error sending SSH public key")
	}

	return nil
}

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"

	"github.com/instructure-bridge/truss-cli/truss"
)

var shellNodeCmd = &cobra.Command{
	Use:   "node [node-name]",
	Short: "Launch a shell on a Truss node.",
	Long: `Launch a shell on a Truss node via SSH.

	Usage:

	truss -e [env] shell node [instance-id]
	truss -e [env] shell node [node-name]

	Examples:

  truss -e nonprod-cmh shell node i-0b35c43cd48ab85ee
	truss -e nonprod-cmh shell node ip-10-12-11-61.us-east-2.compute.internal`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		region, err := getAWSRegionFromKubeconfig()

		if err != nil {
			log.Fatal("Error retrieving AWS region from Kubeconfig:", err.Error())
		}

		awsRole := viper.GetString("aws.assumeRole")

		if awsRole == "" {
			log.Fatal("Error reading aws.assumeRole from config file:", err.Error())
		}

		sess := truss.NewAWSSession(region, awsRole)

		var instanceID string

		if strings.HasPrefix(args[0], "ip-") {
			// Lookup instance-id from Kubernetes if given a node name
			kubeNode, err := describeKubernetesNode(args[0])

			if err != nil {
				log.Fatal("Unable to describe Kubernetes node:", err.Error())
			}

			providerIDParts := strings.Split(kubeNode.Spec.ProviderID, "/")
			instanceID = providerIDParts[len(providerIDParts)-1]
		} else if strings.HasPrefix(args[0], "i-") {
			instanceID = args[0]
		}

		instance := describeInstance(instanceID, sess)
		availabilityZone := *instance.Placement.AvailabilityZone
		hostname := *instance.PrivateDnsName

		sendPublicKey(availabilityZone, instanceID, getSSHPublicKey(), sess)

		jump, err := getJump()

		if err != nil {
			log.Fatal("Error encountered when finding jumpbox:", err.Error())
		}

		execSSHCommand(hostname, jump)

		return nil
	},
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

func execSSHCommand(hostname string, jump string) {
	sshBinary, err := exec.LookPath("ssh")
	if err != nil {
		log.Fatal("Could not find ssh binary")
	}

	target := fmt.Sprintf("ec2-user@%s", hostname)
	proxyJump := fmt.Sprintf("ProxyJump=%s", jump)
	syscall.Exec(sshBinary, []string{"ssh", "-o", proxyJump, target}, os.Environ())
}

func getJump() (string, error) {
	kubeconfigName, err := getKubeconfigName()
	if err != nil {
		return "", err
	}

	clusterEnv := strings.Replace(kubeconfigName, "kubeconfig-truss-", "", 1)
	jumps := viper.GetStringMap("jumps")

	return jumps[clusterEnv].(string), nil
}

func getSSHPublicKey() string {
	homedir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal("Error locating user homedir:", err.Error())
	}

	// TODO: Make this key configurable
	publicKeyPath := fmt.Sprintf("%s/.ssh/id_rsa.pub", homedir)

	publicKeyFile, err := ioutil.ReadFile(publicKeyPath)

	if err != nil {
		log.Fatalf("Unable to read public key from %s:", publicKeyFile)
	}

	return string(publicKeyFile)
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

	node, err := clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return node, nil
}

func describeInstance(instanceID string, sess *session.Session) *ec2.Instance {
	ec2svc := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{&instanceID},
	}

	resp, err := ec2svc.DescribeInstances(params)

	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	return resp.Reservations[0].Instances[0]
}

func sendPublicKey(availabilityZone string, instanceID string, publicKey string, sess *session.Session) {
	svc := ec2instanceconnect.New(sess)

	params := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(availabilityZone),
		InstanceId:       aws.String(instanceID),
		InstanceOSUser:   aws.String("ec2-user"),
		SSHPublicKey:     aws.String(publicKey),
	}

	_, err := svc.SendSSHPublicKey(params)

	if err != nil {
		fmt.Println("There was an error sending SSH public key:", err.Error())
		log.Fatal(err.Error())
	}
}

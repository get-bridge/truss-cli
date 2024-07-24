package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type shellPodOverrides struct {
	APIVersion string           `json:"apiVersion"`
	Metadata   shellPodMetadata `json:"metadata"`
	Spec       shellPodSpec     `json:"spec"`
}

type shellPodMetadata struct {
	Annotations  map[string]string `json:"annotations"`
	GenerateName string            `json:"generateName"`
	Name         string            `json:"name"`
}

type shellPodSpec struct {
	ServiceAccountName string `json:"serviceAccountName"`
}

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Launch a shell container in a Truss cluster",
	Long: `Launches a shell Pod in a Truss cluster. Useful for debugging purposes.

Examples:
  # Run a shell in the current namespace with the default image
  truss shell -e nonprod-cmh

  # Run a shell in the 'vault' Namespace with the 'vault' ServiceAccount
  truss shell -e nonprod-cmh -n vault -s vault

  # Run a shell in the 'cdp-edge' Namespace with Istio disabled
  truss shell -e nonprod-cmh -n cdp-edge --istio=false

  # Run a Pod with a different image and an alternative command
  truss shell -e nonprod-cmh -i nicolaka/netshoot

  # Run a container with a non-default command
  truss shell -e nonprod-cmh -- ls -al /
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set KUBECONFIG and verify availability of kubectl
		kubeconfig, err := getKubeconfig()
		if err != nil {
			return err
		}

		err = os.Setenv("KUBECONFIG", kubeconfig)
		if err != nil {
			return err
		}

		kubectlBinary, err := exec.LookPath("kubectl")
		if err != nil {
			return errors.Wrap(err, "Unable to locate kubectl binary")
		}

		// Generate a user-specific name template for the pod
		user, err := user.Current()
		if err != nil {
			return nil
		}

		replacer := strings.NewReplacer(".", "-", "_", "-")
		podName := "shell-" + replacer.Replace(user.Username)

		// Fetch options from flags
		image, err := cmd.Flags().GetString("image")
		if err != nil {
			return err
		}

		istioEnabled, err := cmd.Flags().GetBool("istio")
		if err != nil {
			return err
		}

		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}

		serviceaccount, err := cmd.Flags().GetString("serviceaccount")
		if err != nil {
			return err
		}

		// Build and exec the kubectl command
		kubectlArgs, err := buildShellKubectlArgs(podName, image, istioEnabled, serviceaccount, namespace, args)
		if err != nil {
			return err
		}

		return syscall.Exec(kubectlBinary, kubectlArgs, os.Environ())
	},
}

func buildShellKubectlArgs(podName string, image string, istioEnabled bool, serviceaccount string, namespace string, cmdArgs []string) ([]string, error) {
	kubectlArgs := []string{"kubectl", "run", podName, "--image-pull-policy=Always", "--restart=Never", "--rm", "--stdin", "--tty", "--attach=true", "--image=" + image}

	// Adding a label so we can easily find and clean these up later
	kubectlArgs = append(kubectlArgs, "--labels=truss.bridgeops.sh/shell=true")

	if istioEnabled {
		kubectlArgs = append(kubectlArgs, "--env=ISTIO_ENABLED=true")
	}

	// Here we are overriding the following in the Pod spec:
	//   * name/generateName to give the pod a unique name
	//   * annotations to enable/disable Istio injection
	overrides := &shellPodOverrides{}
	overrides.APIVersion = "v1"
	overrides.Metadata.Name = ""
	overrides.Metadata.GenerateName = podName + "-"
	overrides.Metadata.Annotations = map[string]string{
		"sidecar.istio.io/inject": strconv.FormatBool(istioEnabled),
	}

	if serviceaccount != "" {
		overrides.Spec.ServiceAccountName = serviceaccount
	}

	overridesJSON, err := json.Marshal(overrides)
	if err != nil {
		return kubectlArgs, err
	}

	if namespace != "" {
		kubectlArgs = append(kubectlArgs, "--namespace="+namespace)
	}

	kubectlArgs = append(kubectlArgs, "--overrides="+string(overridesJSON))

	// Normally , we'll run the image with its default command, but this allows
	// it to be overriden.
	if len(cmdArgs) > 0 {
		kubectlArgs = append(kubectlArgs, "--")
		kubectlArgs = append(kubectlArgs, cmdArgs...)
	}

	return kubectlArgs, nil
}

func init() {
	shellCmd.Flags().StringP("image", "i", "jdharrington/toolbox:latest", "The Docker image to use")
	shellCmd.Flags().StringP("namespace", "n", "", "Namespace to run the shell Pod in")
	shellCmd.Flags().StringP("serviceaccount", "s", "", "ServiceAccount to run the shell Pod as")
	shellCmd.Flags().BoolP("istio", "", false, "Whether or not the shell Pod should be Istio-enabled. Only has effect when launching a shell Pod in an Istio-enabled Namespace")

	rootCmd.AddCommand(shellCmd)
}

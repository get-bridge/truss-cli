package truss

import (
	"errors"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

type vaultAuthAWS struct {
	vaultRole string
	awsRole   string
}

// VaultAuthAWS vault auth
func VaultAuthAWS(vaultRole, awsRole string) VaultAuth {
	return &vaultAuthAWS{
		vaultRole: vaultRole,
		awsRole:   awsRole,
	}
}

func (auth *vaultAuthAWS) LoadCreds() (interface{}, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return stscreds.NewCredentials(sess, auth.awsRole).Get()
}

// Login for VaultAuth interface
func (auth *vaultAuthAWS) Login(data interface{}, port string) error {
	creds, ok := data.(credentials.Value)
	if !ok {
		return errors.New("aws login needs creds")
	}

	cmd := exec.Command("vault", "login", "-method=aws", "role="+auth.vaultRole)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_ADDR=https://localhost:"+port, "VAULT_SKIP_VERIFY=true")
	cmd.Env = append(cmd.Env,
		"AWS_SECRET_ACCESS_KEY="+creds.SecretAccessKey,
		"AWS_ACCESS_KEY_ID="+creds.AccessKeyID,
		"AWS_SESSION_TOKEN="+creds.SessionToken,
	)
	if _, err := cmd.Output(); err != nil {
		return errors.New(string(err.(*exec.ExitError).Stderr))
	}
	return nil
}

package truss

import (
	"errors"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

type vaultAuthAWS struct {
	vaultRole string
	awsRole   string
}

// VaultAuthAWS vault auth
func VaultAuthAWS(vaultRole string, awsRole string) VaultAuth {
	return &vaultAuthAWS{
		vaultRole: vaultRole,
		awsRole:   awsRole,
	}
}

// Login for VaultAuth interface
func (auth *vaultAuthAWS) Login() error {
	// assume aws role
	sess := session.Must(session.NewSession())
	creds, err := stscreds.NewCredentials(sess, auth.awsRole).Get()
	if err != nil {
		return err
	}

	cmd := exec.Command("vault", "login", "-method=aws", "role="+auth.vaultRole)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_ADDR=https://localhost:8200", "VAULT_SKIP_VERIFY=true")
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

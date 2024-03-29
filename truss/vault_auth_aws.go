package truss

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
)

type vaultAuthAWS struct {
	vaultRole string
	awsRole   string
	awsRegion string
}

// VaultAuthAWS vault auth
func VaultAuthAWS(vaultRole, awsRole, awsRegion string) VaultAuth {
	return &vaultAuthAWS{
		vaultRole: vaultRole,
		awsRole:   awsRole,
		awsRegion: awsRegion,
	}
}

func (auth *vaultAuthAWS) LoadCreds() (interface{}, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	creds := stscreds.NewCredentials(sess, auth.awsRole)

	// check valid creds
	_, err = creds.Get()
	if err != nil {
		return nil, err
	}

	return awsutil.GenerateLoginData(creds, "", auth.awsRegion, hclog.Default())
}

// Login for VaultAuth interface
func (auth *vaultAuthAWS) Login(data interface{}, addr string) (string, error) {
	loginData, ok := data.(map[string]interface{})
	if !ok {
		return "", errors.New("aws login needs creds")
	}

	// create a vault client
	client, err := newVaultClient(addr)
	if err != nil {
		return "", err
	}

	loginData["role"] = auth.vaultRole
	secret, err := client.Logical().Write("auth/aws/login", loginData)
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}

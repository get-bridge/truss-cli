package truss

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/vault/api"
	awsauth "github.com/hashicorp/vault/builtin/credential/aws"
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
	creds := stscreds.NewCredentials(sess, auth.awsRole)

	// check valid creds
	_, err = creds.Get()
	if err != nil {
		return nil, err
	}

	return awsauth.GenerateLoginData(creds, "", "")
}

// Login for VaultAuth interface
func (auth *vaultAuthAWS) Login(data interface{}, port string) (string, error) {
	loginData, ok := data.(map[string]interface{})
	if !ok {
		return "", errors.New("aws login needs creds")
	}

	// create a vault client
	loginData["role"] = auth.vaultRole
	config := api.Config{Address: "https://localhost:" + port}
	config.ConfigureTLS(&api.TLSConfig{Insecure: true})
	client, err := api.NewClient(&config)
	if err != nil {
		return "", err
	}

	secret, err := client.Logical().Write("auth/aws/login", loginData)
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}

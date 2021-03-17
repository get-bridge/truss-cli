package truss

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// ObfuscationTarget represents a YAML document that should have all of its values encrypted or decrypted
type ObfuscationTarget struct {
	Secrets map[string]map[string]string `yaml:"secrets"`
}

// NewObfuscationTarget produces a new ObfuscationTarget
func NewObfuscationTarget(raw io.Reader) (t ObfuscationTarget, err error) {
	d := yaml.NewDecoder(raw)
	d.SetStrict(true)
	if err := d.Decode(&t); err != nil {
		return t, ErrSecretFileConfigInvalidYaml
	}
	return
}

// Encrypt encrypts the values of the ObfuscationTarget in-place
func (t *ObfuscationTarget) Encrypt(v *VaultCmd, key string) error {
	r, err := v.Write(fmt.Sprintf("/transit/encrypt/%s", key), map[string]interface{}{
		"batch_input": t.plaintextBatchInput(),
	})
	if err != nil {
		return err
	}

	i := 0
	for name, secret := range t.Secrets {
		for k, _ := range secret {
			br := r.Data["batch_results"].([]interface{})
			r := br[i].(map[string]interface{})
			t.Secrets[name][k] = r["ciphertext"].(string)
			i++
		}
	}

	return nil
}

// Decrypt decrypts the values of the ObfuscationTarget in-place
func (t ObfuscationTarget) Decrypt(v *VaultCmd, key string) error {
	r, err := v.Write(fmt.Sprintf("/transit/decrypt/%s", key), map[string]interface{}{
		"batch_input": t.ciphertextBatchInput(),
	})
	if err != nil {
		return err
	}

	i := 0
	for name, secret := range t.Secrets {
		for k, _ := range secret {
			br := r.Data["batch_results"].([]interface{})
			r := br[i].(map[string]interface{})
			v, err := base64.StdEncoding.DecodeString(r["plaintext"].(string))
			if err != nil {
				return err
			}
			t.Secrets[name][k] = string(v)
			i++
		}
	}

	return nil
}

// Bytes returns a byte slice containing a YAML representation of the current
// state of the ObfuscationTarget (either encrypted or not)
func (t ObfuscationTarget) Bytes() []byte {
	b := bytes.NewBuffer(nil)
	yaml.NewEncoder(b).Encode(t)
	return b.Bytes()
}

func (t ObfuscationTarget) plaintextBatchInput() (out []map[string]interface{}) {
	for _, secret := range t.Secrets {
		for _, v := range secret {
			out = append(out, map[string]interface{}{
				"plaintext": base64.StdEncoding.EncodeToString([]byte(v)),
			})
		}
	}
	return
}

func (t ObfuscationTarget) ciphertextBatchInput() (out []map[string]interface{}) {
	for _, secret := range t.Secrets {
		for _, v := range secret {
			out = append(out, map[string]interface{}{
				"ciphertext": string(v),
			})
		}
	}
	return
}

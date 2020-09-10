package truss

import "io/ioutil"

// encryptAndSaveToDisk encrypts and saves to disk
func encryptAndSaveToDisk(vault VaultCmd, transitKeyName string, filePath string, raw []byte) error {
	enc, err := vault.Encrypt(transitKeyName, raw)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, enc, 0644)
}

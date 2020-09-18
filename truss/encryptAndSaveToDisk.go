package truss

import (
	"io/ioutil"
	"os"
	"path"
)

// encryptAndSaveToDisk encrypts and saves to disk
func encryptAndSaveToDisk(vault *VaultCmd, transitKeyName string, filePath string, raw []byte) error {
	enc, err := vault.Encrypt(transitKeyName, raw)
	if err != nil {
		return err
	}

	// ensure dir exists
	if err := os.MkdirAll(path.Dir(filePath), 0744); err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, enc, 0644)
}

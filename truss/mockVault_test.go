package truss

import "strings"

type mockVault struct {
	// mock saved vault commands
	commands [][]string
	secrets  map[string]interface{}
}

func (*mockVault) PortForward() (string, error) {
	return "", nil
}

func (*mockVault) ClosePortForward() error {
	return nil
}

func (m *mockVault) Run(args []string) ([]byte, error) {
	m.commands = append(m.commands, args)
	return []byte{}, nil
}

func (*mockVault) Decrypt(transitKeyName string, encrypted []byte) ([]byte, error) {
	decrypted := strings.Replace(string(encrypted), "-encrypted", "", 1)
	return []byte(decrypted), nil
}

func (*mockVault) Encrypt(transitKeyName string, raw []byte) ([]byte, error) {
	return append(raw, []byte("-encrypted")...), nil
}

func (*mockVault) GetToken() (string, error) {
	return "", nil
}

func (m *mockVault) GetMap(vaultPath string) (map[string]string, error) {
	return m.secrets[vaultPath].(map[string]string), nil
}

func (m *mockVault) GetPath(vaultPath string) ([]byte, error) {
	return m.secrets[vaultPath].([]byte), nil
}

func (m *mockVault) ListPath(vaultPath string) ([]string, error) {
	keys := []string{}
	for k := range m.secrets {
		keys = append(keys, k)
	}
	return keys, nil
}

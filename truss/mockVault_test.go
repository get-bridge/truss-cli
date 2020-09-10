package truss

import "strings"

type mockVault struct {
	// mock saved vault commands
	commands [][]string
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

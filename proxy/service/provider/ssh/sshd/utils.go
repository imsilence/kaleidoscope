package sshd

import (
	"crypto/rand"
	"crypto/rsa"

	"golang.org/x/crypto/ssh"
)

func generateSigner() (ssh.Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return ssh.NewSignerFromKey(key)
}

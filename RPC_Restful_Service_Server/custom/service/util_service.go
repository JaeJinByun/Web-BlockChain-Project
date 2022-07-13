package service

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func Hash(UserID string) []byte {
	var hash [32]byte
	hash = sha256.Sum256([]byte(UserID))
	return hash[:]
}
func Sign(privKey ecdsa.PrivateKey, hashid []byte) []byte {
	// 개인키로 HashID를 서명한다.
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, hashid)
	if err != nil {
		fmt.Println(err)
	}
	sig := append(r.Bytes(), s.Bytes()...)
	return sig
}

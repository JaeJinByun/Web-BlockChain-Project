package rpccustom

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	UserId     string           `json:"UserId"`
	PrivateKey ecdsa.PrivateKey `json:"PrivateKey"`
	PublicKey  []byte           `json:"PublicKey"`
	Address    string           `json:"Address"`
}

func (w *Wallet) NewWallet(userId string) {
	w.UserId = userId
	w.PrivateKey = NewPrivateKey()
	w.PublicKey = w.NewPublicKey()
	w.Address = w.NewAddress()
}
func NewPrivateKey() ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	fmt.Println(elliptic.P256())
	fmt.Println(privateKey.Curve)
	if err != nil {
		fmt.Println(err)
	}
	return *privateKey
}
func (w *Wallet) NewPublicKey() []byte {
	publicKey := append(w.PrivateKey.PublicKey.X.Bytes(), w.PrivateKey.PublicKey.Y.Bytes()...)
	return publicKey
}
func (w *Wallet) NewAddress() string {
	publicHashed := HashPublicKey(w.PublicKey)
	afterVersion := append([]byte{byte(1)}, publicHashed...)
	checksum := Checksum(afterVersion)
	fullData := append(afterVersion, checksum...)
	address := base58.Encode(fullData)
	return address
}

//공개 키 해싱
func HashPublicKey(pubKey []byte) []byte {
	//공개 키를 sha256으로 해싱
	pubSHA256 := sha256.Sum256(pubKey)
	//sha256으로 해싱한 키를 RIPEND160으로 해싱
	RIPEMDHasher := ripemd160.New()
	_, err := RIPEMDHasher.Write(pubSHA256[:])
	if err != nil {
		fmt.Println(err)
	}
	pubRIPEMD160 := RIPEMDHasher.Sum(nil)
	return pubRIPEMD160
}
func Checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:4]
}

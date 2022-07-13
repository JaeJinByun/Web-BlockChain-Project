package rpccustom

import (
	"crypto/ecdsa"
)

type Arg string
type UserInfo struct {
	UserId   string
	WAddr    string
	PrivateK ecdsa.PrivateKey
	PublicK  []byte
}
type Bundle int

func (t *Bundle) GetAddress(arg *Arg, reply *UserInfo) error {
	wallet := &Wallet{}
	wallet.NewWallet(string(*arg))
	*reply = UserInfo{wallet.UserId, wallet.Address, wallet.PrivateKey, wallet.PublicKey}
	return nil
}

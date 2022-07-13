package rpcclient

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"
)

type Arg string
type UserInfo struct {
	UserId   string           `json:"UserId"`
	WAddr    string           `json:"WAddr"`
	PrivateK ecdsa.PrivateKey `json:"PrivateK"`
	PublicK  []byte           `json:"PublicK"`
}

func Send(userId string, port string) *UserInfo {
	gob.Register(elliptic.P256())
	client, err := rpc.Dial("tcp", "localhost"+port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()
	arg := userId
	reply := &UserInfo{} // == new(UserInfo) --> pointer 반환
	err = client.Call("Bundle.GetAddress", &arg, reply)
	if err != nil {
		log.Fatal("GetAddr error:", err)
	}
	fmt.Printf(`You've been just registered!
    id: %v
    addr: %v
    private: %x
    public: %x
    `, reply.UserId, reply.WAddr, reply.PrivateK, reply.PublicK)
	return reply
}
func ToJson(userinfo *UserInfo) []byte {
	result, _ := json.Marshal(userinfo)
	return result
}

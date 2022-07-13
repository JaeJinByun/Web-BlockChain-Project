package restfulapi

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
)

type UserInfoReq struct {
	UserId   string           `json:"UserId"`
	WAddr    string           `json:"WAddr"`
	PrivateK ecdsa.PrivateKey `json:"PrivateK"`
	PublicK  []byte           `json:"PublicK"`
}
type Result struct { // UserInfoReq 의 PrivateK에서 D값만 받는 스트럭쳐
	UserId string `json:"UserId"`
	// PrivateK ecdsa.PrivateKey `json:"PrivateK"`
	PrivateK []byte `json:"PrivateK"`
	//  WAddr    string `json:"WAddr"`
	//  PublicK  []byte `json:"PublicK"`
}
type CreateBlock struct {
	UserId   string
	LogDb    int
	RId      int
	Content  string
	PrivateK []byte
}
type CurveParams struct {
	P       *big.Int
	N       *big.Int
	B       *big.Int
	Gx      *big.Int
	Gy      *big.Int
	Name    string
	BitSize int
}

func (c CreateBlock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	block := &CreateBlock{}
	body, _ := ioutil.ReadAll(r.Body) // []byte return
	json.Unmarshal(body, block)
	//block.PrivateK.PublicKey.Curve = elliptic.P256()
	fmt.Printf(`
    ui: %s
    ld: %s
    ri: %s
    ct: %s
    pk: %x
    `, block.UserId, block.LogDb, block.RId, block.Content, block.PrivateK)
	res, _ := http.Post("http://localhost:8081/create_block", "application/json", r.Body)
	fmt.Println(res)
}
func (u UserInfoReq) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.Body)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// cors(Cross-Origin Resource Sharing)를 위한 preflight. OPTIONS 메서드 요청만 처리
	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
		return
	}
	// preflight이 해결되면 요청되는 POST를 위한 코드
	res, _ := http.Post("http://localhost:7777/create", "application/json", r.Body)
	body, _ := ioutil.ReadAll(res.Body)
	userinfo := &UserInfoReq{}
	json.Unmarshal(body, &userinfo)
	result := &Result{}
	result.UserId = userinfo.UserId
	userinfo.PrivateK.Curve = elliptic.P256()
	// result.PrivateK = userinfo.PrivateK
	privatek, _ := x509.MarshalECPrivateKey(&userinfo.PrivateK)
	result.PrivateK = privatek
	fmt.Println(userinfo)
	fmt.Printf(`
    USERID: %v
    WADDR : %v
    PRIVATE : %x
    PUBLIC : %x`, userinfo.UserId, userinfo.WAddr, userinfo.PrivateK, userinfo.PublicK)

	mresult, _ := json.Marshal(result)
	w.Write(mresult) // 지갑정보 java에 전달
	// SendKey(userinfo.UserId, userinfo.PublicK) // 퍼블릭키 블록체인에 전달
}
func Start() {
	mux := http.NewServeMux()
	mux.Handle("/create_wallet", UserInfoReq{})
	mux.Handle("/create_block", CreateBlock{})
	log.Print("Listening port 7779...")
	http.ListenAndServe(":7779", mux)
}

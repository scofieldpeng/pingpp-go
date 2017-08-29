package token

import (
	"fmt"
	"log"
	"net/url"
	"time"

	pingpp "github.com/scofieldpeng/pingpp-go/pingpp"
)

type Client struct {
	B          pingpp.Backend
	Key        string
	PrivateKey string
}

func getC(authKey ...pingpp.AuthKey) Client {
	if len(authKey) == 0 {
		authKey = make([]pingpp.AuthKey, 1)
		authKey[0].Key = pingpp.Key
		authKey[0].PrivateKey = pingpp.AccountPrivateKey
	}
	return Client{pingpp.GetBackend(pingpp.APIBackend), authKey[0].Key, authKey[0].PrivateKey}
}

// 发送 Token 请求
func New(params *pingpp.TokenParams, authKey ...pingpp.AuthKey) (*pingpp.Token, error) {
	return getC(authKey...).New(params)
}

func (c Client) New(params *pingpp.TokenParams) (*pingpp.Token, error) {
	start := time.Now()
	paramsString, errs := pingpp.JsonEncode(params)
	if errs != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("TokenParams Marshall Errors is : %q\n", errs)
		}
	}
	if pingpp.LogLevel > 2 {
		log.Printf("params of card request to pingpp is :\n %v\n ", string(paramsString))
	}

	token := &pingpp.Token{}
	errch := c.B.Call("POST", "/tokens", c.Key, c.PrivateKey, nil, paramsString, token)
	if errch != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("%v\n", errch)
		}
		return nil, errch
	}
	if pingpp.LogLevel > 2 {
		log.Println("Token completed in ", time.Since(start))
	}
	return token, errch

}

//查询指定 token 对象
func Get(tok_id string, authKey ...pingpp.AuthKey) (*pingpp.Token, error) {
	return getC(authKey...).Get(tok_id)
}

func (c Client) Get(tok_id string) (*pingpp.Token, error) {
	var body *url.Values
	body = &url.Values{}
	token := &pingpp.Token{}
	err := c.B.Call("GET", fmt.Sprintf("/tokens/%v", tok_id), c.Key, c.PrivateKey, body, nil, token)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Get Card error: %v\n", err)
		}
	}
	return token, err
}

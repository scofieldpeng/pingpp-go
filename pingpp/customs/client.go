package customs

import (
	"fmt"
	"log"

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

func New(params *pingpp.CustomsParams, authKey ...pingpp.AuthKey) (*pingpp.Customs, error) {
	return getC(authKey...).New(params)
}

func (c Client) New(params *pingpp.CustomsParams) (*pingpp.Customs, error) {
	paramsString, _ := pingpp.JsonEncode(params)
	customs := &pingpp.Customs{}
	err := c.B.Call("POST", "/customs", c.Key, c.PrivateKey, nil, paramsString, customs)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("New Customs error: %v\n", err)
		}
	}
	return customs, err
}

func Get(Id string, authKey ...pingpp.AuthKey) (*pingpp.Customs, error) {
	return getC(authKey...).Get(Id)
}

func (c Client) Get(Id string) (*pingpp.Customs, error) {
	customs := &pingpp.Customs{}
	err := c.B.Call("GET", fmt.Sprintf("/customs/%s", Id), c.Key, c.PrivateKey, nil, nil, customs)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Get Customs error: %v\n", err)
		}
	}
	return customs, err
}

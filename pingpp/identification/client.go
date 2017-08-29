package identification

import (
	"log"

	pingpp "github.com/scofieldpeng/pingpp-go/pingpp"
)

const (
	IDENTIFY_IDCARD   = "id_card"
	IDENTIFY_BANKCARD = "bank_card"
)

type Client struct {
	B          pingpp.Backend
	Key        string
	PrivateKey string
}

func New(params *pingpp.IdentificationParams, authKey ...pingpp.AuthKey) (*pingpp.IdentificationResult, error) {
	return getC(authKey...).New(params)
}

func (c Client) New(params *pingpp.IdentificationParams) (*pingpp.IdentificationResult, error) {
	paramsString, errs := pingpp.JsonEncode(params)
	if errs != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("IdentificationParams Marshall Errors is : %q/n", errs)
		}
		return nil, errs
	}
	if pingpp.LogLevel > 2 {
		log.Printf("params of identification request to pingpp is :\n %v\n ", string(paramsString))
	}
	identificationResult := &pingpp.IdentificationResult{}

	err := c.B.Call("POST", "/identification", c.Key, c.PrivateKey, nil, paramsString, identificationResult)
	return identificationResult, err
}

func getC(authKey ...pingpp.AuthKey) Client {
	if len(authKey) == 0 {
		authKey = make([]pingpp.AuthKey, 1)
		authKey[0].Key = pingpp.Key
		authKey[0].PrivateKey = pingpp.AccountPrivateKey
	}
	return Client{pingpp.GetBackend(pingpp.APIBackend), authKey[0].Key, authKey[0].PrivateKey}
}

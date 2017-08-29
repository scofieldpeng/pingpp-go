package event

import (
	"log"
	"net/url"

	pingpp "github.com/scofieldpeng/pingpp-go/pingpp"
)

type Client struct {
	B          pingpp.Backend
	Key        string
	PrivateKey string
}

func Get(id string, authKey ...pingpp.AuthKey) (*pingpp.Event, error) {
	return getC(authKey...).Get(id)
}

func (c Client) Get(id string) (*pingpp.Event, error) {
	var body *url.Values
	body = &url.Values{}
	eve := &pingpp.Event{}
	err := c.B.Call("GET", "/events/"+id, c.Key, c.PrivateKey, body, nil, eve)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Get Event error: %v\n", err)
		}
	}
	return eve, err
}

func getC(authKey ...pingpp.AuthKey) Client {
	if len(authKey) == 0 {
		authKey = make([]pingpp.AuthKey, 1)
		authKey[0].Key = pingpp.Key
		authKey[0].PrivateKey = pingpp.AccountPrivateKey
	}
	return Client{pingpp.GetBackend(pingpp.APIBackend), authKey[0].Key, authKey[0].PrivateKey}
}

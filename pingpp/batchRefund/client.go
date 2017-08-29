package batchRefund

import (
	"fmt"
	"log"
	"net/url"

	pingpp "github.com/scofieldpeng/pingpp-go/pingpp"
)

type Client struct {
	B          pingpp.Backend
	Key        string
	PrivateKey string
}

func getC(authKey ...pingpp.AuthKey) Client {
	if len(authKey) == 0 {
		authKey = make([]pingpp.AuthKey,1)
		authKey[0].Key = pingpp.Key
		authKey[0].PrivateKey = pingpp.AccountPrivateKey
	}
	return Client{pingpp.GetBackend(pingpp.APIBackend), authKey[0].Key, authKey[0].PrivateKey}
}

func New(params *pingpp.BatchRefundParams, key ...pingpp.AuthKey) (*pingpp.BatchRefund, error) {
	return getC(key...).New(params)
}

func (c Client) New(params *pingpp.BatchRefundParams) (*pingpp.BatchRefund, error) {
	paramsString, _ := pingpp.JsonEncode(params)
	batchRefund := &pingpp.BatchRefund{}
	err := c.B.Call("POST", "/batch_refunds", c.Key, c.PrivateKey, nil, paramsString, batchRefund)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("New BatchRefunds error: %v\n", err)
		}
	}
	return batchRefund, err
}

func Get(Id string,authKey ...pingpp.AuthKey) (*pingpp.BatchRefund, error) {
	return getC(authKey...).Get(Id)
}

func (c Client) Get(Id string ) (*pingpp.BatchRefund, error) {
	batchRefund := &pingpp.BatchRefund{}
	err := c.B.Call("GET", fmt.Sprintf("/batch_refunds/%s", Id), c.Key, c.PrivateKey, nil, nil, batchRefund)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Get BatchRefunds error: %v\n", err)
		}
	}
	return batchRefund, err
}

func List(params *pingpp.PagingParams,authKey ...pingpp.AuthKey) (*pingpp.BatchRefundlList, error) {
	return getC().List(params)
}

func (c Client) List(params *pingpp.PagingParams) (*pingpp.BatchRefundlList, error) {
	body := &url.Values{}
	params.Filters.AppendTo(body)

	batchRefundlList := &pingpp.BatchRefundlList{}
	err := c.B.Call("GET", "/batch_refunds", c.Key, c.PrivateKey, body, nil, batchRefundlList)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Get BatchRefunds List error: %v\n", err)
		}
	}
	return batchRefundlList, err
}

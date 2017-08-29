package batchTransfer

import (
	"fmt"
	"log"
	"net/url"

	pingpp "github.com/scofieldpeng/pingpp-go/pingpp"
)

type Client struct {
	B   pingpp.Backend
	Key string
	PrivateKey string
}

func getC(authKey ...pingpp.AuthKey) Client {
	if len(authKey) == 0 {
		authKey = make([]pingpp.AuthKey,1)
		authKey[0].Key = pingpp.Key
		authKey[0].PrivateKey = pingpp.AccountPrivateKey
	}
	return Client{pingpp.GetBackend(pingpp.APIBackend), authKey[0].Key,authKey[0].PrivateKey}
}

func New(params *pingpp.BatchTransferParams,authKey ...pingpp.AuthKey) (*pingpp.BatchTransfer, error) {
	return getC(authKey...).New(params)
}

func (c Client) New(params *pingpp.BatchTransferParams) (*pingpp.BatchTransfer, error) {
	paramsString, _ := pingpp.JsonEncode(params)
	batchTransfer := &pingpp.BatchTransfer{}
	err := c.B.Call("POST", "/batch_transfers", c.Key,c.PrivateKey, nil, paramsString, batchTransfer)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("New BatchTransfer error: %v\n", err)
		}
	}
	return batchTransfer, err
}

func Get(Id string,authKey ...pingpp.AuthKey) (*pingpp.BatchTransfer, error) {
	return getC(authKey...).Get(Id)
}

func (c Client) Get(Id string) (*pingpp.BatchTransfer, error) {
	batchTransfer := &pingpp.BatchTransfer{}
	err := c.B.Call("GET", fmt.Sprintf("/batch_transfers/%s", Id), c.Key,c.PrivateKey, nil, nil, batchTransfer)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Get Batchtransfer error: %v\n", err)
		}
	}
	return batchTransfer, err
}

func List(params *pingpp.PagingParams,authKey ...pingpp.AuthKey) (*pingpp.BatchTransferlList, error) {
	return getC(authKey...).List(params)
}

func (c Client) List(params *pingpp.PagingParams) (*pingpp.BatchTransferlList, error) {
	body := &url.Values{}
	params.Filters.AppendTo(body)

	batchTransferlList := &pingpp.BatchTransferlList{}
	err := c.B.Call("GET", "/batch_transfers", c.Key,c.PrivateKey, body, nil, batchTransferlList)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Get Batchtransfer List error: %v\n", err)
		}
	}
	return batchTransferlList, err
}

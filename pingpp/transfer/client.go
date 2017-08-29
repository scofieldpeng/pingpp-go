package transfer

import (
	"log"
	"net/url"
	"strconv"

	pingpp "github.com/scofieldpeng/pingpp-go/pingpp"
)

type Client struct {
	B          pingpp.Backend
	Key        string
	PrivateKey string
}

func New(params *pingpp.TransferParams, authKey ...pingpp.AuthKey) (*pingpp.Transfer, error) {
	return getC(authKey...).New(params)
}

func (c Client) New(params *pingpp.TransferParams) (*pingpp.Transfer, error) {
	paramsString, errs := pingpp.JsonEncode(params)
	if errs != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("ChargeParams Marshall Errors is : %q/n", errs)
		}
		return nil, errs
	}
	if pingpp.LogLevel > 2 {
		log.Printf("params of redEnvelope request to pingpp is :\n %v\n ", string(paramsString))
	}
	transfer := &pingpp.Transfer{}
	err := c.B.Call("POST", "/transfers", c.Key, c.PrivateKey, nil, paramsString, transfer)
	return transfer, err
}

func Update(id string, authKey ...pingpp.AuthKey) (*pingpp.Transfer, error) {
	return getC(authKey...).Update(id)
}

func (c Client) Update(id string) (*pingpp.Transfer, error) {
	cancelParams := struct {
		Status string `json:"status"`
	}{
		Status: "canceled",
	}

	paramsString, _ := pingpp.JsonEncode(cancelParams)
	transfer := &pingpp.Transfer{}
	err := c.B.Call("PUT", "/transfers/"+id, c.Key, c.PrivateKey, nil, paramsString, transfer)
	return transfer, err
}

// Get returns the details of a redenvelope.
func Get(id string, authKey ...pingpp.AuthKey) (*pingpp.Transfer, error) {
	return getC(authKey...).Get(id)
}

func (c Client) Get(id string) (*pingpp.Transfer, error) {
	var body *url.Values
	body = &url.Values{}
	transfer := &pingpp.Transfer{}
	err := c.B.Call("GET", "/transfers/"+id, c.Key, c.PrivateKey, body, nil, transfer)
	return transfer, err
}

// List returns a list of transfer.
func List(params *pingpp.TransferListParams, authKey ...pingpp.AuthKey) *Iter {
	return getC(authKey...).List(params)
}

func (c Client) List(params *pingpp.TransferListParams) *Iter {
	type transferList struct {
		pingpp.ListMeta
		Values []*pingpp.Transfer `json:"data"`
	}

	var body *url.Values
	var lp *pingpp.ListParams

	if params != nil {
		body = &url.Values{}

		if params.Created > 0 {
			body.Add("created", strconv.FormatInt(params.Created, 10))
		}
		params.AppendTo(body)
		lp = &params.ListParams
	}

	return &Iter{pingpp.GetIter(lp, body, func(b url.Values) ([]interface{}, pingpp.ListMeta, error) {
		list := &transferList{}
		err := c.B.Call("GET", "/transfers", c.Key, c.PrivateKey, &b, nil, list)

		ret := make([]interface{}, len(list.Values))
		for i, v := range list.Values {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

type Iter struct {
	*pingpp.Iter
}

func (i *Iter) Transfer() *pingpp.Transfer {
	return i.Current().(*pingpp.Transfer)
}

func getC(authKey ...pingpp.AuthKey) Client {
	if len(authKey) == 0 {
		authKey = make([]pingpp.AuthKey, 1)
		authKey[0].Key = pingpp.Key
		authKey[0].PrivateKey = pingpp.AccountPrivateKey
	}
	return Client{pingpp.GetBackend(pingpp.APIBackend), authKey[0].Key, authKey[0].PrivateKey}
}

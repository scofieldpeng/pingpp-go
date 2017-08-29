package redEnvelope

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

func New(params *pingpp.RedEnvelopeParams, authKey ...pingpp.AuthKey) (*pingpp.RedEnvelope, error) {
	return getC(authKey...).New(params)
}

func (c Client) New(params *pingpp.RedEnvelopeParams) (*pingpp.RedEnvelope, error) {
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
	redEnvelope := &pingpp.RedEnvelope{}
	err := c.B.Call("POST", "/red_envelopes", c.Key, c.PrivateKey, nil, paramsString, redEnvelope)
	return redEnvelope, err
}

func Get(id string, authKey ...pingpp.AuthKey) (*pingpp.RedEnvelope, error) {
	return getC(authKey...).Get(id)
}

func (c Client) Get(id string) (*pingpp.RedEnvelope, error) {
	var body *url.Values
	body = &url.Values{}
	redEnvelope := &pingpp.RedEnvelope{}
	err := c.B.Call("GET", "/red_envelopes/"+id, c.Key, c.PrivateKey, body, nil, redEnvelope)
	return redEnvelope, err
}

func List(params *pingpp.RedEnvelopeListParams, authKey ...pingpp.AuthKey) *Iter {
	return getC(authKey...).List(params)
}

func (c Client) List(params *pingpp.RedEnvelopeListParams) *Iter {
	type redEnvelopeList struct {
		pingpp.ListMeta
		Values []*pingpp.RedEnvelope `json:"data"`
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
		list := &redEnvelopeList{}
		err := c.B.Call("GET", "/red_envelopes", c.Key, c.PrivateKey, &b, nil, list)

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

func (i *Iter) RedEnvelope() *pingpp.RedEnvelope {
	return i.Current().(*pingpp.RedEnvelope)
}

func getC(authKey ...pingpp.AuthKey) Client {
	if len(authKey) == 0 {
		authKey = make([]pingpp.AuthKey, 1)
		authKey[0].Key = pingpp.Key
		authKey[0].PrivateKey = pingpp.AccountPrivateKey
	}
	return Client{pingpp.GetBackend(pingpp.APIBackend), authKey[0].Key, authKey[0].PrivateKey}
}

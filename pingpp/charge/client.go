package charge

import (
	"log"
	"net/url"
	"strconv"
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
	return Client{B: pingpp.GetBackend(pingpp.APIBackend), Key: authKey[0].Key, PrivateKey: authKey[0].PrivateKey}
}

// 发送 charge 请求
func New(params *pingpp.ChargeParams, authKey ...pingpp.AuthKey) (*pingpp.Charge, error) {
	return getC(authKey...).New(params)
}

func (c Client) New(params *pingpp.ChargeParams) (*pingpp.Charge, error) {
	start := time.Now()
	paramsString, errs := pingpp.JsonEncode(params)
	if errs != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("ChargeParams Marshall Errors is : %q\n", errs)
		}
	}
	if pingpp.LogLevel > 2 {
		log.Printf("params of charge request to pingpp is :\n %v\n ", string(paramsString))
	}

	charge := &pingpp.Charge{}
	errch := c.B.Call("POST", "/charges", c.Key, c.PrivateKey, nil, paramsString, charge)
	if errch != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("%v\n", errch)
		}
		return nil, errch
	}
	if pingpp.LogLevel > 2 {
		log.Println("Charge completed in ", time.Since(start))
	}
	return charge, errch

}

// 撤销charge，此接口仅接受线下 isv_scan、isv_wap、isv_qr 渠道的订单调用
func Reverse(id string, authKey ...pingpp.AuthKey) (*pingpp.Charge, error) {
	return getC(authKey...).Reverse(id)
}

func (c Client) Reverse(id string) (*pingpp.Charge, error) {
	var body *url.Values
	body = &url.Values{}
	charge := &pingpp.Charge{}
	err := c.B.Call("POST", "/charges/"+id+"/reverse", c.Key, c.PrivateKey, body, nil, charge)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Reverse Charge error: %v\n", err)
		}
	}
	return charge, err
}

//查询指定 charge 对象
func Get(id string, authKey ...pingpp.AuthKey) (*pingpp.Charge, error) {
	return getC(authKey...).Get(id)
}

func (c Client) Get(id string) (*pingpp.Charge, error) {
	var body *url.Values
	body = &url.Values{}
	charge := &pingpp.Charge{}
	err := c.B.Call("GET", "/charges/"+id, c.Key, c.PrivateKey, body, nil, charge)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Get Charge error: %v\n", err)
		}
	}
	return charge, err
}

// 查询 charge 列表
func List(appId string, params *pingpp.ChargeListParams, authKey ...pingpp.AuthKey) *Iter {
	return getC(authKey...).List(appId, params)
}

func (c Client) List(appId string, params *pingpp.ChargeListParams) *Iter {
	type chargeList struct {
		pingpp.ListMeta
		Values []*pingpp.Charge `json:"data"`
	}

	var body *url.Values
	var lp *pingpp.ListParams

	if params == nil {
		params = &pingpp.ChargeListParams{}
	}
	params.Filters.AddFilter("app[id]", "", appId)
	body = &url.Values{}
	if params.Created > 0 {
		body.Add("created", strconv.FormatInt(params.Created, 10))
	}
	params.AppendTo(body)
	lp = &params.ListParams

	return &Iter{pingpp.GetIter(lp, body, func(b url.Values) ([]interface{}, pingpp.ListMeta, error) {
		list := &chargeList{}
		err := c.B.Call("GET", "/charges", c.Key, c.PrivateKey, &b, nil, list)

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

func (i *Iter) Charge() *pingpp.Charge {
	return i.Current().(*pingpp.Charge)
}

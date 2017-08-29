package card

import (
	"fmt"
	pingpp "github.com/scofieldpeng/pingpp-go/pingpp"
	"log"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	B   pingpp.Backend
	Key string
	PrivateKey string
}

func getC(authKey...pingpp.AuthKey) Client {
	if len(authKey) == 0 {
		authKey = make([]pingpp.AuthKey, 1)
		authKey[0].Key = pingpp.Key
		authKey[0].PrivateKey = pingpp.AccountPrivateKey
	}
	return Client{pingpp.GetBackend(pingpp.APIBackend),authKey[0].Key,authKey[0].PrivateKey}
}

// 发送 card 请求
func New(cus_id string, params *pingpp.CardParams,authKey ...pingpp.AuthKey) (*pingpp.Card, error) {
	return getC(authKey...).New(cus_id, params)
}

func (c Client) New(cus_id string, params *pingpp.CardParams) (*pingpp.Card, error) {
	start := time.Now()
	paramsString, errs := pingpp.JsonEncode(params)
	if errs != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("CardParams Marshall Errors is : %q\n", errs)
		}
	}
	if pingpp.LogLevel > 2 {
		log.Printf("params of card request to pingpp is :\n %v\n ", string(paramsString))
	}

	card := &pingpp.Card{}
	errch := c.B.Call("POST", fmt.Sprintf("/customers/%v/sources", cus_id), c.Key, c.PrivateKey,nil, paramsString, card)
	if errch != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("%v\n", errch)
		}
		return nil, errch
	}
	if pingpp.LogLevel > 2 {
		log.Println("Card completed in ", time.Since(start))
	}
	return card, errch

}

//查询指定 card 对象
func Get(cus_id string, card_id string,authKey ...pingpp.AuthKey) (*pingpp.Card, error) {
	return getC(authKey...).Get(cus_id, card_id)
}

func (c Client) Get(cus_id string, card_id string) (*pingpp.Card, error) {
	var body *url.Values
	body = &url.Values{}
	card := &pingpp.Card{}
	err := c.B.Call("GET", fmt.Sprintf("/customers/%v/sources/%v", cus_id, card_id), c.Key, c.PrivateKey,body, nil, card)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Get Card error: %v\n", err)
		}
	}
	return card, err
}

//删除指定 card 对象
func Delete(cus_id string, card_id string,authKey ...pingpp.AuthKey) (map[string]interface{}, error) {
	return getC(authKey...).Delete(cus_id, card_id)
}

func (c Client) Delete(cus_id string, card_id string) (map[string]interface{}, error) {
	var body *url.Values
	body = &url.Values{}
	res := make(map[string]interface{})
	err := c.B.Call("DELETE", fmt.Sprintf("/customers/%v/sources/%v", cus_id, card_id), c.Key,c.PrivateKey, body, nil, &res)
	if err != nil {
		if pingpp.LogLevel > 0 {
			log.Printf("Get Card error: %v\n", err)
		}
	}
	return res, err
}

// 查询 card 列表
func List(cus_id string, params *pingpp.CardListParams,authKey ...pingpp.AuthKey) *Iter {
	return getC(authKey...).List(cus_id, params)
}

func (c Client) List(cus_id string, params *pingpp.CardListParams) *Iter {
	type CardList struct {
		pingpp.ListMeta
		Values []*pingpp.Card `json:"data"`
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
		list := &CardList{}
		err := c.B.Call("GET", fmt.Sprintf("/customers/%v/sources", cus_id), c.Key, c.PrivateKey,&b, nil, list)

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

func (i *Iter) Card() *pingpp.Card {
	return i.Current().(*pingpp.Card)
}

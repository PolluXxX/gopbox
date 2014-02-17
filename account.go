package gopbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/polluxxx/goauth2"
)

type Account struct {
	Link      string       `json:"referral_link"`
	Name      string       `json:"display_name"`
	Uid       int          `json:"uid"`
	Country   string       `json:"country"`
	QuotaInfo AccountQuota `json:"quota_info"`

	Token *goauth2.OAuthToken
}

type AccountQuota struct {
	Shared int `json:"shared"`
	Quota  int `json:"quota"`
	Normal int `json:"normal"`
}

func NewAccount(token *goauth2.OAuthToken) (account *Account, err error) {
	account = new(Account)
	account.Token = token

	body, err := account.Call("account/info", "GET", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (account *Account) Call(endpoint, method string, params map[string]string) (body []byte, err error) {

	// Building query
	u := url.Values{}

	for p, v := range params {
		u.Set(p, v)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", MetaUrl, endpoint), strings.NewReader(u.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", account.Token.Type, account.Token.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		// Catch response
		apiErr := new(ApiError)
		err = json.Unmarshal(body, apiErr)
		if err != nil {
			return nil, err
		}

		apiErr.HTTPCode = resp.StatusCode
		return nil, apiErr
	}

	return body, nil
}

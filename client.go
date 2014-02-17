package gopbox

import (
	"net/url"
	"unicode"

	"github.com/polluxxx/goauth2"
)

const (
	AuthUrl    = "https://www.dropbox.com/1/oauth2/authorize"
	TokenUrl   = "https://api.dropbox.com/1/oauth2/token"
	MetaUrl    = "https://api.dropbox.com/1/"
	ContentUrl = "https://api-content.dropbox.com/1/"
)

type DropboxApi struct {
	Client *goauth2.Client
	Scope  string
}

func NewDropboxApi(id, secret, redirectUrl string) (*DropboxApi, error) {
	a := goauth2.Api{
		AuthUrl:  AuthUrl,
		TokenUrl: TokenUrl,
	}

	c := goauth2.NewClient(id, secret, redirectUrl, &a)

	da := DropboxApi{
		Client: c,
	}

	return &da, nil
}

func (dropbox *DropboxApi) GetAuthUrl(state string) (uri *url.URL, err error) {
	m := make(map[string]string)

	m["state"] = state

	return dropbox.Client.GetAuthUrl(m)
}

func (dropbox *DropboxApi) FinalizeAuth(code string) (account *Account, err error) {
	token, err := dropbox.Client.Exchange(code)
	if err != nil {
		return nil, err
	}

	// Hi, I'm Dropbox and I lowercase the token_type
	a := []rune(token.Type)
	a[0] = unicode.ToUpper(a[0])
	token.Type = string(a)

	account, err = NewAccount(token)
	if err != nil {
		return nil, err
	}

	return account, nil
}

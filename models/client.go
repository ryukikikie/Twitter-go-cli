package models

import (
	"net/url"

	"github.com/gomodule/oauth1/oauth"
)

type Client interface {
	ReqGet(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error)
	ReqPost(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error)
}

package models

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/gomodule/oauth1/oauth"
)

type TwClient struct {
	Client oauth.Client
}

func NewTwClient() TwClient {
	return TwClient{
		Client: oauth.Client{
			TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
			ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
			TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
		}}
}

func (tc *TwClient) ReqGet(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error) {
	resp, err := tc.Client.Get(nil, credentials, urlStr, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Get request return %d", resp.StatusCode)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	return buf, err
}

func (tc *TwClient) ReqPost(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error) {
	resp, err := tc.Client.Post(nil, credentials, urlStr, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Get request return %d", resp.StatusCode)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	return buf, err
}

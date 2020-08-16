package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/gomodule/oauth1/oauth"
)

type MockClient struct {
	client oauth.Client
}

func (mc *MockClient) ReqGet(credentials *oauth.Credentials, urlStr string) ([]byte, error) {
	switch urlStr {
	case "https://api.twitter.com/1.1/account/verify_credentials.json":
		responseBody, err := ioutil.ReadFile("./test/responseData/getUser.json")
		if err != nil {
			return nil, err
		}
		return responseBody, nil
	case "https://api.twitter.com/1.1/statuses/home_timeline.json":
		return nil, errors.New("Not implimented")
	}
	return nil, errors.New("Not implimented")
}

func (mc *MockClient) ReqPost(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error) {
	//Return test data
	fmt.Println("Call mock Post function! but it's not implimented")
	return nil, errors.New("Not implimented")
}

var twitterMockClient MockClient = MockClient{
	client: oauth.Client{}}

func TestGetUser(t *testing.T) {
	var actual User
	expected := User{
		Name:       "Miki.masumomo",
		ScreenName: "m_miki0108",
	}

	GetUser(&twitterMockClient, nil, &actual) // Don't need credential

	if actual.Name != expected.Name || actual.ScreenName != expected.ScreenName {
		t.Fatalf("User must be %v", expected)
	}
}

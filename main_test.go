package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"

	"github.com/gomodule/oauth1/oauth"
	"github.com/ryukikikie/twitter-go-cli/test"
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
		responseBody, err := ioutil.ReadFile("./test/responseData/getTimeLine.json")
		if err != nil {
			return nil, err
		}
		return responseBody, nil
	default:
		return nil, errors.New("Not implimented")
	}
}

type TweetFormat struct {
	CreatedAt string
	Text      string
}

var NumberOfLinePerTweet = int(3)

func (mc *MockClient) ReqPost(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error) {
	switch urlStr {
	case "https://api.twitter.com/1.1/statuses/update.json":
		responseBody, err := ioutil.ReadFile("./test/responseData/createPost.json")
		if err != nil {
			return nil, err
		}
		return responseBody, nil
	default:
		return nil, errors.New("Not implimented")
	}
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

func TestGetTimeLine(t *testing.T) {

	expectedTweets := []TweetFormat{
		TweetFormat{
			CreatedAt: "(Created at Sun Aug 16 23:45:54 +0000 2020)",
			Text:      "@test_user congrats!〜🥳",
		},
		TweetFormat{
			CreatedAt: "(Created at Sun Aug 16 23:43:17 +0000 2020)",
			Text:      "Good morning!\nおはようございます〜！"},
	}
	outputs := test.CaptureOutput(func() {
		GetTimeLine(&twitterMockClient, nil, 2) // Don't need credential
	})
	output := strings.Split(outputs, "\n")
	if len(expectedTweets) != len(output)/NumberOfLinePerTweet {
		t.Fatalf("Number of tweet must be %v, result:%v", len(expectedTweets), len(output)/NumberOfLinePerTweet)
	}
	for i := 0; i < len(expectedTweets); i++ {
		if expectedTweets[i].CreatedAt != output[(i*NumberOfLinePerTweet)] {
			t.Fatalf("CreatedAt must be %v, result:%v", expectedTweets[i].CreatedAt, output[(i*NumberOfLinePerTweet)])
		}
		if "Tweet" != output[(i*NumberOfLinePerTweet)+1] {
			t.Fatalf("Seccond line must be Tweet, result:%v", output[(i*NumberOfLinePerTweet)+1])
		}
		expectedTweetText := strings.Split(expectedTweets[i].Text, "\n")
		for j := 0; j < len(expectedTweetText); j++ {
			if expectedTweetText[j] != output[(i*NumberOfLinePerTweet)+2+j] {
				t.Fatalf("Text must be  %v, result:%v", expectedTweetText[j], output[(i*NumberOfLinePerTweet)+2+j])
			}
		}
	}
}

func TestCreatePost(t *testing.T) {
	expectedTweet := TweetFormat{
		CreatedAt: "(Created at Mon Aug 17 11:16:12 +0000 2020)",
		Text:      "This is tweet is created by my twitter CLI client using Golang.. just test",
	}

	outputs := test.CaptureOutput(func() {
		CreatePost(&twitterMockClient, nil, expectedTweet.Text) // Don't need credential
	})
	output := strings.Split(outputs, "\n")
	if "Your tweet has been posted!" != output[0] {
		t.Fatalf("First line must be 'Your tweet has been posted!', result:%v", output[0])
	}
	if fmt.Sprintf("%v %v", expectedTweet.CreatedAt, expectedTweet.Text) != output[1] {
		t.Fatalf("CreatedAt and posted tweet must be %v %v, result:%v", expectedTweet.CreatedAt, expectedTweet.Text, output[1])
	}
}

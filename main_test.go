package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"

	"github.com/gomodule/oauth1/oauth"
	"github.com/ryukikikie/twitter-go-cli/controller"
	"github.com/ryukikikie/twitter-go-cli/models"
	"github.com/ryukikikie/twitter-go-cli/test"
)

type MockClient struct {
	client oauth.Client
}

func (mc *MockClient) ReqGet(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error) {
	switch urlStr {
	case "https://api.twitter.com/1.1/account/verify_credentials.json":
		if form != nil {
			return nil, errors.New("getUser cannot take url.Values as a argument")
		}
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
	CreatedAtText string
	Text          string
	Author        string
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
	var actual models.User
	expected := models.User{
		Name:       "Miki.masumomo",
		ScreenName: "m_miki0108",
	}

	controller.GetUser(&twitterMockClient, nil, &actual) // Don't need credential

	if actual.Name != expected.Name || actual.ScreenName != expected.ScreenName {
		t.Fatalf("User must be %v", expected)
	}
}

func TestGetTimeLine(t *testing.T) {

	expectedTweets := []TweetFormat{
		TweetFormat{
			CreatedAtText: "(at 08-16-2020 23:45:54 Sun)",
			Text:          "@test_user congrats!〜🥳",
			Author:        "地獄寺",
		},
		TweetFormat{
			CreatedAtText: "(at 08-16-2020 23:43:17 Sun)",
			Text:          "Good morning!\nおはようございます〜！",
			Author:        "ヤマダ",
		},
	}
	outputs := test.CaptureOutput(func() {
		controller.GetTimeLine(&twitterMockClient, nil, 2) // Don't need credential
	})
	output := strings.Split(outputs, "\n")
	if len(expectedTweets) != len(output)/NumberOfLinePerTweet {
		t.Fatalf("Number of tweets must be %v, actual:%v", len(expectedTweets), len(output)/NumberOfLinePerTweet)
	}
	for i := 0; i < len(expectedTweets); i++ {
		if "---Tweet---" != output[(i*NumberOfLinePerTweet)] {
			t.Fatalf("First line of output must be Tweet, actual:%v", output[(i*NumberOfLinePerTweet)])
		}
		expectedTweetText := strings.Split(expectedTweets[i].Text, "\n")
		for j := 0; j < len(expectedTweetText); j++ {
			if expectedTweetText[j] != output[(i*NumberOfLinePerTweet)+1+j] {
				t.Fatalf("Text must be  %v, actual:%v", expectedTweetText[j], output[(i*NumberOfLinePerTweet)+1+j])
			}
		}
		expectedString := "Published by " + expectedTweets[i].Author + " " + expectedTweets[i].CreatedAtText
		if expectedString != output[(i*NumberOfLinePerTweet)+len(expectedTweetText)+1] {
			t.Fatalf("Author and CreatedAt must be %v, actual:%v", expectedString, output[(i*NumberOfLinePerTweet)+len(expectedTweetText)+1])
		}
	}
}

func TestCreatePost(t *testing.T) {
	expectedTweet := TweetFormat{
		CreatedAtText: "(Created at 08-17-2020 11:16:12 Mon)",
		Text:          "This is tweet is created by my twitter CLI client using Golang.. just test",
	}

	outputs := test.CaptureOutput(func() {
		controller.CreatePost(&twitterMockClient, nil, expectedTweet.Text) // Don't need credential
	})
	output := strings.Split(outputs, "\n")
	if "Your tweet has been posted!" != output[0] {
		t.Fatalf("First line must be 'Your tweet has been posted!', actual:%v", output[0])
	}
	if fmt.Sprintf("%v %v", expectedTweet.CreatedAtText, expectedTweet.Text) != output[1] {
		t.Fatalf("CreatedAt and posted tweet must be %v %v, actual:%v", expectedTweet.CreatedAtText, expectedTweet.Text, output[1])
	}
}

func TestHelp(t *testing.T) {
	allCommand := []string{
		"timeline", "tweet", "exit", "clear", "exit",
	}

	outputs := test.CaptureOutput(func() {
		controller.Help()
	})
	for _, c := range allCommand {
		if !strings.Contains(outputs, c) {
			t.Fatalf("%v should be included in help", c)
		}
	}
}

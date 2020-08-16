package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gomodule/oauth1/oauth"
)

// Wrap original oauth client method to mock test easily
type Client interface {
	ReqGet(credentials *oauth.Credentials, urlStr string) ([]byte, error)
	ReqPost(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error)
}

type TwClient struct {
	client oauth.Client
}

func (tc *TwClient) ReqGet(credentials *oauth.Credentials, urlStr string) ([]byte, error) {
	resp, err := tc.client.Get(nil, credentials, urlStr, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//TODO If stats code is not 200 return err
	buf, err := ioutil.ReadAll(resp.Body)
	return buf, err
}

func (tc *TwClient) ReqPost(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error) {
	resp, err := tc.client.Post(nil, credentials, urlStr, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//TODO If stats code is not 200 return err
	buf, err := ioutil.ReadAll(resp.Body)
	return buf, err
}

var twitterClient TwClient = TwClient{
	client: oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}}

var credPath = flag.String("config", "config.json", "Path to configuration file containing the application's credentials.")

type User struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

func readCredentials() error {
	b, err := ioutil.ReadFile(*credPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &twitterClient.client.Credentials)
}

func GetTimeLine(c Client, tokenCred *oauth.Credentials) {
	urlStr := "https://api.twitter.com/1.1/statuses/home_timeline.json"
	buf, err := c.ReqGet(tokenCred, urlStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
}
func CreatePost(c Client, tokenCred *oauth.Credentials, tweet string) {
	v := url.Values{}
	v.Set("status", tweet)
	urlStr := "https://api.twitter.com/1.1/statuses/update.json"
	buf, err := c.ReqPost(tokenCred, urlStr, v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
}

func GetUser(c Client, tokenCred *oauth.Credentials, user *User) {

	urlStr := "https://api.twitter.com/1.1/account/verify_credentials.json"
	buf, err := c.ReqGet(tokenCred, urlStr)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(buf, &user)
	fmt.Println(string(buf))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	var user = User{}
	// Obtaining a request token
	if err := readCredentials(); err != nil {
		log.Fatal(err)
	}

	tempCred, err := twitterClient.client.RequestTemporaryCredentials(nil, "oob", nil)
	if err != nil {
		log.Fatal(err)
	}

	u := twitterClient.client.AuthorizationURL(tempCred, nil)

	fmt.Println("Enter verification code:")
	openbrowser(u)

	var verificationCode string
	fmt.Scanln(&verificationCode)

	tokenCred, _, err := twitterClient.client.RequestToken(nil, tempCred, verificationCode)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Welcome to TWITTER-GOCLI-APP!")
	GetUser(&twitterClient, tokenCred, &user)
	for {
		fmt.Printf("@%v=>", user.ScreenName)
		var command string
		fmt.Scanln(&command)
		switch command {
		case "timeline":
			GetTimeLine(&twitterClient, tokenCred)
		case "tweet":
			fmt.Println("Make a tweet through CLI🧊")
			inputReader := bufio.NewReader(os.Stdin)
			input, _ := inputReader.ReadString('\n')
			CreatePost(&twitterClient, tokenCred, input)
		case "exit":
			fmt.Print("CLI terminating")
			//insert settimeout & loop below
			for i := 0; i < 3; i++ {
				time.Sleep(500 * time.Millisecond)
				fmt.Print(".")
			}
			fmt.Println()
			return
		}
		fmt.Println()
	}
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

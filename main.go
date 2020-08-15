package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/gomodule/oauth1/oauth"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

var credPath = flag.String("config", "config.json", "Path to configuration file containing the application's credentials.")

type User struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

var user = User{}

func readCredentials() error {
	b, err := ioutil.ReadFile(*credPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &oauthClient.Credentials)
}

func getTimeLine(tokenCred *oauth.Credentials) {
	resp, err := oauthClient.Get(nil, tokenCred,
		"https://api.twitter.com/1.1/statuses/home_timeline.json", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Fatal(err)
	}
}
func createPost(tokenCred *oauth.Credentials, tweet string) {
	v := url.Values{}
	v.Set("status", tweet)
	urlStr := "https://api.twitter.com/1.1/statuses/update.json"
	resp, err := oauthClient.Post(nil, tokenCred, urlStr, v)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Fatal(err)
	}
}

func getUser(tokenCred *oauth.Credentials) {

	resp, err := oauthClient.Get(nil, tokenCred,
		"https://api.twitter.com/1.1/account/verify_credentials.json", nil)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(buf, &user)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	// Obtaining a request token
	if err := readCredentials(); err != nil {
		log.Fatal(err)
	}

	tempCred, err := oauthClient.RequestTemporaryCredentials(nil, "oob", nil)
	if err != nil {
		log.Fatal(err)
	}

	u := oauthClient.AuthorizationURL(tempCred, nil)

	fmt.Println("Enter verification code:")
	openbrowser(u)

	var verificationCode string
	fmt.Scanln(&verificationCode)

	tokenCred, _, err := oauthClient.RequestToken(nil, tempCred, verificationCode)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Welcome to TWITTER-GOCLI-APP!")

	getUser(tokenCred)
	for {
		fmt.Printf("@%v=>", user.ScreenName)
		var command string
		fmt.Scanln(&command)
		switch command {
		case "timeline":
			getTimeLine(tokenCred)
			fmt.Print("\n")
		case "tweet":
			var tweet string
			fmt.Scanln(&tweet)
			createPost(tokenCred, tweet)
			fmt.Print("\n")
		case "exit":
			fmt.Println("CLI terminating")
			return
		}
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

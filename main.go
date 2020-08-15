package main

import (
	"bufio"
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
	"time"

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
var client TestClient = 

type TestClient interface {
	multiply(a int, b int) int
}

type data string

//mocked by main_test
func (data data) multiply(a int, b int) int {
	return a * b
}

func TestedFunction() int {
	fmt.Println("called TestedFunction")
	return client.multiply(1, 3)
}

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
		case "tweet":
			fmt.Println("Make a tweet through CLI🧊")
			inputReader := bufio.NewReader(os.Stdin)
			input, _ := inputReader.ReadString('\n')
			createPost(tokenCred, input)
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

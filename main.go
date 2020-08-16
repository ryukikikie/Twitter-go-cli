package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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

type Tweet struct {
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

type TweetArr []Tweet

var user = User{}
var tweet = Tweet{}

var iconArr = []string{"🐉", "🐍", "🐲"}

func readCredentials() error {
	b, err := ioutil.ReadFile(*credPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &oauthClient.Credentials)
}

func getTimeLine(tokenCred *oauth.Credentials) {
	v := url.Values{}
	v.Set("count", "1")
	resp, err := oauthClient.Get(nil, tokenCred,
		"https://api.twitter.com/1.1/statuses/home_timeline.json", v)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var tweetArr TweetArr
	json.Unmarshal(buf, &tweetArr)
	for _, v := range tweetArr {
		fmt.Println("(Created at " + v.CreatedAt + ")")
		fmt.Println("Tweet")
		fmt.Println(v.Text)
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
	dt := time.Now()
	fmt.Println(dt.Format("01-02-2006 15:04:05 Mon"))

	getUser(tokenCred)
	for {
		randomIndex := rand.Intn(len(iconArr))
		avatar := iconArr[randomIndex]
		fmt.Printf("[%v%v]", avatar, user.ScreenName)
		var command string
		fmt.Scanln(&command)
		switch command {
		case "timeline":
			getTimeLine(tokenCred)
		case "tweet":
			fmt.Println("Tweet through CLI🧊")
			inputReader := bufio.NewReader(os.Stdin)
			input, _ := inputReader.ReadString('\n')
			createPost(tokenCred, input)
		case "clear":
			fmt.Print("\033[H\033[2J")
		case "exit":
			fmt.Print("CLI terminating")
			//insert settimeout & loop below
			for i := 0; i < 3; i++ {
				time.Sleep(500 * time.Millisecond)
				fmt.Print(".")
			}
			fmt.Println()
			return
		default:
			fmt.Println("Input command doesn't exit 😂, or some typo")
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

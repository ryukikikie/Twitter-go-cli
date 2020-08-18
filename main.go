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
	"strings"
	"time"

	"github.com/gomodule/oauth1/oauth"
)

type Command struct {
	description string
	option      map[string]string // it's not used now
}

var commands = map[string]Command{
	"timeline": Command{
		description: "Get timeline",
		option:      make(map[string]string),
	},
	"tweet": Command{
		description: "Post tweet",
		option:      make(map[string]string),
	},
	"clear": Command{
		description: "Clear console",
		option:      make(map[string]string),
	},
	"exit": Command{
		description: "Terminate this cli client",
		option:      make(map[string]string),
	},
	"help": Command{
		description: "Show how to use",
		option:      make(map[string]string),
	},
}

// Wrap original oauth client method to mock test easily
type Client interface {
	ReqGet(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error)
	ReqPost(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error)
}

type TwClient struct {
	client oauth.Client
}

func (tc *TwClient) ReqGet(credentials *oauth.Credentials, urlStr string, form url.Values) ([]byte, error) {
	resp, err := tc.client.Get(nil, credentials, urlStr, form)
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
	resp, err := tc.client.Post(nil, credentials, urlStr, form)
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

var twitterClient TwClient = TwClient{
	client: oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}}

var credPath = flag.String("config", "config.json", "Path to configuration file containing the application's credentials.")

var CustomTimeFormatter = "01-02-2006 15:04:05 Mon"

type CustomTime struct {
	time.Time
}

func NowCustomTime() CustomTime {
	return CustomTime{time.Now()}
}
func (t *CustomTime) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse(time.RubyDate, strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

func (t CustomTime) Format() string {
	return time.Time(t.Time).Format(CustomTimeFormatter)
}

type User struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type Tweet struct {
	CreatedAt CustomTime `json:"created_at"`
	Text      string     `json:"text"`
}

var iconArr = []string{"🐉", "🐍", "🐲"}

func readCredentials() error {
	b, err := ioutil.ReadFile(*credPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &twitterClient.client.Credentials)
}

func GetTimeLine(c Client, tokenCred *oauth.Credentials, limit int) {
	values := url.Values{}
	values.Set("count", string(limit))
	urlStr := "https://api.twitter.com/1.1/statuses/home_timeline.json"
	buf, err := c.ReqGet(tokenCred, urlStr, values)
	if err != nil {
		log.Fatal(err)
	}
	var tweets []Tweet
	err = json.Unmarshal(buf, &tweets)
	if err != nil {
		log.Fatal(err)
	}
	for _, tweet := range tweets {
		fmt.Printf("(Created at %s)\n", tweet.CreatedAt.Format())
		fmt.Println("---Tweet---")
		fmt.Println(tweet.Text)
	}
}
func CreatePost(c Client, tokenCred *oauth.Credentials, tweet string) {
	values := url.Values{}
	values.Set("status", tweet)
	urlStr := "https://api.twitter.com/1.1/statuses/update.json"
	buf, err := c.ReqPost(tokenCred, urlStr, values)
	if err != nil {
		log.Fatal(err)
	}
	var createdTweet Tweet
	json.Unmarshal(buf, &createdTweet)
	fmt.Println("Your tweet has been posted!")
	fmt.Printf("(Created at %s) %s\n", createdTweet.CreatedAt.Format(), createdTweet.Text)
}

func GetUser(c Client, tokenCred *oauth.Credentials, user *User) {

	urlStr := "https://api.twitter.com/1.1/account/verify_credentials.json"
	buf, err := c.ReqGet(tokenCred, urlStr, nil)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(buf, &user)
	if err != nil {
		log.Fatal(err)
	}
}

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func Exit() {
	fmt.Print("CLI terminating")
	//insert settimeout & loop below
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Println()
	os.Exit(1)
}

func Help() {
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("        <command> [arguments]")
	fmt.Println()
	fmt.Println("The commands are:")
	fmt.Println()
	for name, command := range commands {
		fmt.Printf("        %s:%s\n", name, command.description)
		if len(command.option) > 0 {
			fmt.Println("        options")
		}
		for o, description := range command.option {
			fmt.Printf("%s : %s\n", o, description)
		}
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

	dt := NowCustomTime()
	fmt.Println(dt.Format())

	GetUser(&twitterClient, tokenCred, &user)

	for {
		randomIndex := rand.Intn(len(iconArr))
		avatar := iconArr[randomIndex]
		fmt.Printf("[%v%v]", avatar, user.ScreenName)
		var command string
		fmt.Scanln(&command)

		switch command {
		case "timeline":
			GetTimeLine(&twitterClient, tokenCred, 2)
		case "tweet":
			fmt.Println("Tweet through CLI🧊")
			inputReader := bufio.NewReader(os.Stdin)
			input, _ := inputReader.ReadString('\n')
			CreatePost(&twitterClient, tokenCred, input)
		case "clear":
			Clear()
		case "exit":
			Exit()
		case "help":
			Help()
		default:
			fmt.Println("Input command doesn't exit 😂, or some typo")
			fmt.Println("If you need any help, see 'help'")
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

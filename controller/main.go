package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/gomodule/oauth1/oauth"
	"github.com/ryukikikie/twitter-go-cli/models"
)

func NewTwClient() models.TwClient {
	return models.NewTwClient()
}

func GetTimeLine(c models.Client, tokenCred *oauth.Credentials, limit int) {
	values := url.Values{}
	values.Set("count", strconv.Itoa(limit))
	urlStr := "https://api.twitter.com/1.1/statuses/home_timeline.json"
	buf, err := c.ReqGet(tokenCred, urlStr, values)
	if err != nil {
		log.Fatal(err)
	}
	var tweets []models.Tweet
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

func CreatePost(c models.Client, tokenCred *oauth.Credentials, tweet string) {
	values := url.Values{}
	values.Set("status", tweet)
	urlStr := "https://api.twitter.com/1.1/statuses/update.json"
	buf, err := c.ReqPost(tokenCred, urlStr, values)
	if err != nil {
		log.Fatal(err)
	}
	var createdTweet models.Tweet
	json.Unmarshal(buf, &createdTweet)
	fmt.Println("Your tweet has been posted!")
	fmt.Printf("(Created at %s) %s\n", createdTweet.CreatedAt.Format(), createdTweet.Text)
}

func GetUser(c models.Client, tokenCred *oauth.Credentials, user *models.User) {

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
func main() {
	NewTwClient()
}
func NowCustomTime() models.CustomTime {
	return models.CustomTime{time.Now()}
}
func NewUser() models.User {
	return models.User{}
}

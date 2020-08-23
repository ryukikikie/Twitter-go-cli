package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/ryukikikie/twitter-go-cli/controller"
)

var twitterClient = controller.NewTwClient()

var credPath = flag.String("config", "config.json", "Path to configuration file containing the application's credentials.")

var iconArr = []string{"🐉", "🐍", "🐲"}

func readCredentials() error {
	b, err := ioutil.ReadFile(*credPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &twitterClient.Client.Credentials)
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

func main() {
	var user = controller.NewUser()
	// Obtaining a request token
	if err := readCredentials(); err != nil {
		log.Fatal(err)
	}

	tempCred, err := twitterClient.Client.RequestTemporaryCredentials(nil, "oob", nil)
	if err != nil {
		log.Fatal(err)
	}

	u := twitterClient.Client.AuthorizationURL(tempCred, nil)

	fmt.Println("Enter verification code:")
	openbrowser(u)

	var verificationCode string
	fmt.Scanln(&verificationCode)

	tokenCred, _, err := twitterClient.Client.RequestToken(nil, tempCred, verificationCode)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Welcome to TWITTER-GOCLI-APP!")

	dt := controller.NowCustomTime()
	fmt.Println(dt.Format())

	controller.GetUser(&twitterClient, tokenCred, &user)

	for {
		randomIndex := rand.Intn(len(iconArr))
		avatar := iconArr[randomIndex]
		fmt.Printf("[%v%v]", avatar, user.ScreenName)
		var command string
		fmt.Scanln(&command)

		switch command {
		case "timeline":
			controller.GetTimeLine(&twitterClient, tokenCred, 2)
		case "tweet":
			fmt.Println("Tweet through CLI🧊")
			inputReader := bufio.NewReader(os.Stdin)
			input, _ := inputReader.ReadString('\n')
			controller.CreatePost(&twitterClient, tokenCred, input)
		case "clear":
			Clear()
		case "exit":
			Exit()
		case "help":
			controller.Help()
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

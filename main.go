package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	a "github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	REPO       = os.Getenv("REPO")
	REPO_URL   = "https://github.com/" + REPO
	LATEST_URL = REPO_URL + "/releases/latest"

	TWITTER_ACCESS_TOKEN                = os.Getenv("TWITTER_ACCESS_TOKEN")
	TWITTER_ACCESS_TOKEN_SECRET         = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	TWITTER_CONSUMER_KEY                = os.Getenv("TWITTER_CONSUMER_KEY")
	TWITTER_CONSUMER_SECRET             = os.Getenv("TWITTER_CONSUMER_SECRET")
	twitterClient               tweeter = api(*a.NewTwitterApiWithCredentials(
		TWITTER_ACCESS_TOKEN,
		TWITTER_ACCESS_TOKEN_SECRET,
		TWITTER_CONSUMER_KEY,
		TWITTER_CONSUMER_SECRET,
	))
)

type tweeter interface {
	tweet(message string) (string, error)
}

type api a.TwitterApi

func (api api) tweet(message string) (string, error) {
	tweet, err := (a.TwitterApi(api)).PostTweet(message, nil)
	return tweet.Text, err
}

func main() {
	lambda.Start(Handler)
}

func Handler() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", LATEST_URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	tag := getLatestTag(resp.Request.URL.Path)

	message := fmt.Sprintf("%s %s released! check the new features on GitHub.\n%s", REPO, tag, REPO_URL)
	msg, err := twitterClient.tweet(message)
	if !strings.Contains(msg, tag) {
		log.Fatalf("failed to tweet: %s", msg)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func getLatestTag(url string) string {
	s := strings.Split(url, "/")
	return s[len(s)-1]
}

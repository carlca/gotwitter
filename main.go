package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/ChimeraCoder/anaconda"
)

// Credentials struct is read from ~/gotwitter/config.json config file
type Credentials struct {
	ConsumerKey    string `json:"consumerkey"`
	ConsumerSecret string `json:"consumersecret"`
	AccessToken    string `json:"accesstoken"`
	AccessSecret   string `json:"accesssecret"`
}

func (c *Credentials) String() string {
	return fmt.Sprintf("ConsumerKey: '%v'\nConsumerSecret:'%v'\nAccessToken: '%v'\nAccessSecret:'%v'\n",
		c.ConsumerKey, c.ConsumerSecret, c.AccessToken, c.AccessSecret)
}

func main() {
	// get current user
	usr, err := user.Current()
	if err != nil {
		reportError(err)
	}
	// get config file location
	credFile := path.Join(usr.HomeDir, ".gotwitter/config.json")
	file, err := os.Open(credFile)
	if err != nil {
		reportError(err)
	}

	creds, err := readConfig(file)

	fmt.Println(creds)

	api := getTwitterAPI(creds)

	searchResult, _ := api.GetSearch("golang", nil)
	for _, tweet := range searchResult.Statuses {
		fmt.Println(tweet.Text)
		fmt.Println("")
	}
}

func readConfig(file *os.File) (*Credentials, error) {
	decoder := json.NewDecoder(file)
	creds := &Credentials{}
	err := decoder.Decode(&creds)
	return creds, err
}

func getTwitterAPI(creds *Credentials) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(creds.ConsumerKey)
	anaconda.SetConsumerSecret(creds.ConsumerSecret)
	return anaconda.NewTwitterApi(creds.AccessToken, creds.AccessSecret)
}

func reportError(err error) {
	fmt.Println("error:", err)
}

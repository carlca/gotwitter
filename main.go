package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/ChimeraCoder/anaconda"
	"github.com/pkg/errors"
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

	usr, err := getCurrentUser()
	configPath := path.Join(usr.HomeDir, ".gotwitter/config.json")
	file, err := openConfigFile(configPath)
	creds, err := readConfig(file)
	fmt.Println(creds)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	api := getTwitterAPI(creds)

	searchResult, _ := api.GetSearch("golang", nil)
	for _, tweet := range searchResult.Statuses {
		fmt.Println(tweet.Text)
		fmt.Println("")
	}
}

func getCurrentUser() (*user.User, error) {
	usr, err := user.Current()
	if err != nil {
		errors.Wrap(err, "error retriving user.Current()")
	}
	return usr, err
}

func openConfigFile(configPath string) (*os.File, error) {
	file, err := os.Open(configPath)
	if err != nil {
		errors.Wrap(err, "error reading config path")
	}
	return file, err
}

func readConfig(file *os.File) (*Credentials, error) {
	decoder := json.NewDecoder(file)
	creds := &Credentials{}
	err := decoder.Decode(&creds)
	if err != nil {
		errors.Wrap(err, "error in decoder.Decode(&creds)")
	}
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

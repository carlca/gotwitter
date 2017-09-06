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

func terminateOnError(err error) {
	if err != nil {
		fmt.Printf("%+v", errors.WithStack(err))
		os.Exit(1) // or anything else ...
	}
}

func main() {

	var err error

	// get the current os user
	usr, err := user.Current()
	terminateOnError(err)

	// get the current users's configuration path for the gotwitter application

	var config string
	config = path.Join(usr.HomeDir, ".gotwitter/config.json")
	_, err = os.Stat(config)
	terminateOnError(err)

	// open a file based on the specified config path
	var file *os.File
	file, err = os.Open(config)
	terminateOnError(err)
	defer file.Close()

	// read the config json file into the Credentials struct
	var creds *Credentials
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&creds)
	terminateOnError(err)

	// get TwitterAPI based on stored credentials
	var api *anaconda.TwitterApi
	anaconda.SetConsumerKey(creds.ConsumerKey)
	anaconda.SetConsumerSecret(creds.ConsumerSecret)
	api = anaconda.NewTwitterApi(creds.AccessToken, creds.AccessSecret)
	// I'm thinking of adding a context field to the methodWrapper struct...
	// err := errors.New("Error creating TwitterAPI")

	// search current Twitter timeline for golang content
	tweets, err := api.GetSearch("golang", nil)
	terminateOnError(err)

	for _, tweet := range tweets.Statuses {
		fmt.Println(tweet.Text)
		fmt.Println("")
	}
}

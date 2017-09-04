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

// methodWrapper struct provides pre and post function call wrapping
type methodWrapper struct {
	Err error
}

func (m *methodWrapper) do(f func(m *methodWrapper)) {
	if m.Err != nil {
		return
	}
	f(m)
	if m.Err != nil {
		errors.WithStack(m.Err)
	}
}

func main() {

	// construct an instance of methodWrapper
	m := &methodWrapper{}

	// get the current os user
	var usr *user.User
	m.do((func(m *methodWrapper) {
		usr, m.Err = user.Current()
	}))

	// get the current users's configuration path for the gotwitter application
	var config string
	m.do((func(m *methodWrapper) {
		config = path.Join(usr.HomeDir, ".gotwitter/config.json")
		_, m.Err = os.Stat(config)
	}))

	// open a file based on the specified config path
	var file *os.File
	m.do((func(m *methodWrapper) {
		file, m.Err = os.Open(config)
	}))

	// read the config json file into the Credentials struct
	var creds *Credentials
	m.do((func(m *methodWrapper) {
		decoder := json.NewDecoder(file)
		defer file.Close()
		m.Err = decoder.Decode(&creds)
	}))

	// get TwitterAPI based on stored credentials
	var api *anaconda.TwitterApi
	m.do((func(m *methodWrapper) {
		anaconda.SetConsumerKey(creds.ConsumerKey)
		anaconda.SetConsumerSecret(creds.ConsumerSecret)
		api = anaconda.NewTwitterApi(creds.AccessToken, creds.AccessSecret)
		m.Err = nil
		// I'm thinking of adding a context field to the methodWrapper struct...
		// err := errors.New("Error creating TwitterAPI")
	}))

	// search current Twitter timeline for golang content
	var tweets anaconda.SearchResponse
	m.do((func(m *methodWrapper) {
		tweets, m.Err = api.GetSearch("golang", nil)
		if m.Err == nil {
			for _, tweet := range tweets.Statuses {
				fmt.Println(tweet.Text)
				fmt.Println("")
			}
		}
	}))

	// final check on any errors which may have occurred
	if m.Err != nil {
		fmt.Printf("%+v", m.Err)
	}
}

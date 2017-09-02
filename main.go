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

// the methodWrapper struct provides pre and post function call wrapping
type methodWrapper struct {
	Result interface{}
	Err    error
}

func (mw *methodWrapper) do(f func() (interface{}, error)) {
	if mw.Err != nil {
		return
	}
	result, err := f()
	if err != nil {
		errors.WithStack(err)
	}
	mw.Err = err
	mw.Result = result
}

func main() {

	// construct an instance of methodWrapper
	m := &methodWrapper{}

	// get the current os user
	m.do((func() (interface{}, error) {
		result, err := user.Current()
		return result, err
	}))

	// get the current users's configuration path for the gotwitter application
	m.do((func() (interface{}, error) {
		usr := m.Result.(*user.User)
		configPath := path.Join(usr.HomeDir, ".gotwitter/config.json")
		_, err := os.Stat(configPath)
		return configPath, err
	}))

	// open a file based on the specified config path
	m.do((func() (interface{}, error) {
		config := m.Result.(string)
		file := &os.File{}
		file, err := os.Open(config)
		return file, err
	}))

	// read the config json file into the Credentials struct
	m.do((func() (interface{}, error) {
		file := m.Result.(*os.File)
		decoder := json.NewDecoder(file)
		creds := &Credentials{}
		err := decoder.Decode(&creds)
		return creds, err
	}))

	// get TwitterAPI based on stored credentials
	m.do((func() (interface{}, error) {
		creds := m.Result.(*Credentials)
		anaconda.SetConsumerKey(creds.ConsumerKey)
		anaconda.SetConsumerSecret(creds.ConsumerSecret)
		api := anaconda.NewTwitterApi(creds.AccessToken, creds.AccessSecret)
		if api == nil {
			err := errors.New("Error creating TwitterAPI")
			return api, err
		}
		return api, nil
	}))

	// search current Twitter timeline for golang content
	m.do((func() (interface{}, error) {
		api := m.Result.(*anaconda.TwitterApi)
		searchResult, err := api.GetSearch("golang", nil)
		if err == nil {
			for _, tweet := range searchResult.Statuses {
				fmt.Println(tweet.Text)
				fmt.Println("")
			}
		}
		return api, err
	}))

	// final check on any errors which may have occurred
	if m.Err != nil {
		fmt.Printf("%+v", m.Err)
	}
}

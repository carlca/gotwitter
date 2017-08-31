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

type methods struct {
	Err error
}

func (m *methods) checkWrap(errMsg string) {
	if m.Err != nil {
		errors.Wrap(m.Err, errMsg)
	}
}

func (m *methods) getCurrentUser() *user.User {
	if m.Err != nil {
		return nil
	}
	usr := &user.User{}
	usr, m.Err = user.Current()
	m.checkWrap("getCurrentUser: error retriving user.Current()")
	return usr
}

func (m *methods) getConfigPath(usr *user.User) string {
	if m.Err != nil {
		return ""
	}
	config := path.Join(usr.HomeDir, ".gotwitter/config.json")
	_, m.Err = os.Stat(config)
	m.checkWrap("getConfigPath: invald config path")
	return config
}

func (m *methods) openConfigFile(configPath string) *os.File {
	if m.Err != nil {
		return nil
	}
	file := &os.File{}
	file, m.Err = os.Open(configPath)
	m.checkWrap("openConfigFile: error reading config path")
	return file
}

func (m *methods) readConfig(file *os.File) *Credentials {
	if m.Err != nil {
		return nil
	}
	decoder := json.NewDecoder(file)
	creds := &Credentials{}
	m.Err = decoder.Decode(&creds)
	m.checkWrap("readConfig: error in decoder.Decode(&creds)")
	return creds
}

func (m *methods) getTwitterAPI(creds *Credentials) *anaconda.TwitterApi {
	if m.Err != nil {
		return nil
	}
	anaconda.SetConsumerKey(creds.ConsumerKey)
	anaconda.SetConsumerSecret(creds.ConsumerSecret)
	api := anaconda.NewTwitterApi(creds.AccessToken, creds.AccessSecret)
	if api == nil {
		errors.Wrap(m.Err, "getTwitterAPI: error in anaconda.NewTwitterApi(...)")
	}
	return api
}

func (m *methods) searchTimeline(api anaconda.TwitterApi, key string) {
	if m.Err != nil {
		return
	}
	var searchResult anaconda.SearchResponse
	searchResult, m.Err = api.GetSearch("golang", nil)
	errors.Wrap(m.Err, "searchTimeline: error in api.GetSearch(...)")
	for _, tweet := range searchResult.Statuses {
		fmt.Println(tweet.Text)
		fmt.Println("")
	}
}

func (m *methods) checkForErrors() {
	if m.Err != nil {
		fmt.Printf("%+v", m.Err)
	}
}

func main() {
	m := &methods{}
	usr := m.getCurrentUser()
	config := m.getConfigPath(usr)
	file := m.openConfigFile(config)
	creds := m.readConfig(file)
	api := m.getTwitterAPI(creds)
	m.searchTimeline(*api, "golang")
	m.checkForErrors()
}

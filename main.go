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

type wrapper struct {
	Err error
}

func (w *wrapper) checkWrap(errMsg string) {
	if w.Err != nil {
		errors.Wrap(w.Err, errMsg)
	}
}

func (w *wrapper) getCurrentUser() *user.User {
	if w.Err != nil {
		return nil
	}
	usr := &user.User{}
	usr, w.Err = user.Current()
	w.checkWrap("getCurrentUser: error retriving user.Current()")
	return usr
}

func (w *wrapper) getConfigPath(usr *user.User) string {
	if w.Err != nil {
		return ""
	}
	config := path.Join(usr.HomeDir, ".gotwitter/config.json")
	_, w.Err = os.Stat(config)
	w.checkWrap("getConfigPath: invald config path")
	return config
}

func (w *wrapper) openConfigFile(configPath string) *os.File {
	if w.Err != nil {
		return nil
	}
	file := &os.File{}
	file, w.Err = os.Open(configPath)
	w.checkWrap("openConfigFile: error reading config path")
	return file
}

func (w *wrapper) readConfig(file *os.File) *Credentials {
	if w.Err != nil {
		return nil
	}
	decoder := json.NewDecoder(file)
	creds := &Credentials{}
	w.Err = decoder.Decode(&creds)
	w.checkWrap("readConfig: error in decoder.Decode(&creds)")
	return creds
}

func (w *wrapper) checkForErrors() {
	if w.Err != nil {
		fmt.Printf("%+v", w.Err)
	}
}

func main() {
	w := &wrapper{}
	usr := w.getCurrentUser()
	config := w.getConfigPath(usr)
	file := w.openConfigFile(config)
	creds := w.readConfig(file)
	fmt.Println(creds)
	w.checkForErrors()

	api := getTwitterAPI(creds)

	searchResult, _ := api.GetSearch("golang", nil)
	for _, tweet := range searchResult.Statuses {
		fmt.Println(tweet.Text)
		fmt.Println("")
	}
}

// func getCurrentUser() (*user.User, error) {
// 	usr, err := user.Current()
// 	if err != nil {
// 		errors.Wrap(err, "error retriving user.Current()")
// 	}
// 	return usr, err
// }

// func openConfigFile(configPath string) (*os.File, error) {
// 	file, err := os.Open(configPath)
// 	if err != nil {
// 		errors.Wrap(err, "error reading config path")
// 	}
// 	return file, err
// }

// func readConfig(file *os.File) (*Credentials, error) {
// 	decoder := json.NewDecoder(file)
// 	creds := &Credentials{}
// 	err := decoder.Decode(&creds)
// 	if err != nil {
// 		errors.Wrap(err, "error in decoder.Decode(&creds)")
// 	}
// 	return creds, err
// }

func getTwitterAPI(creds *Credentials) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(creds.ConsumerKey)
	anaconda.SetConsumerSecret(creds.ConsumerSecret)
	return anaconda.NewTwitterApi(creds.AccessToken, creds.AccessSecret)
}

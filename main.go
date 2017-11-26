package main

import (
	"encoding/json"
	"flag"
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

// get the current os user
func getCurrentUser() *user.User {
	usr, err := user.Current()
	terminateOnError(err)
	return usr
}

// get the current users's configuration path for the gotwitter application
func getUserConfig(usr *user.User) string {
	config := path.Join(usr.HomeDir, ".gotwitter/config.json")
	_, err := os.Stat(config)
	terminateOnError(err)
	return config
}

// open a file based on the specified config path
func openFile(configPath string) *os.File {
	file, err := os.Open(configPath)
	terminateOnError(err)
	return file
}

// read the config json file into the Credentials struct
func readCredentials(file *os.File) *Credentials {
	var creds *Credentials
	decoder := json.NewDecoder(file)
	defer file.Close()
	err := decoder.Decode(&creds)
	terminateOnError(err)
	return creds
}

// get TwitterAPI based on stored credentials
func getTwitterAPI(creds *Credentials) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(creds.ConsumerKey)
	anaconda.SetConsumerSecret(creds.ConsumerSecret)
	api := anaconda.NewTwitterApi(creds.AccessToken, creds.AccessSecret)
	return api
	// I'm thinking of adding a context field to the methodWrapper struct...
	// err := errors.New("Error creating TwitterAPI")
}

// API extends the *anaconda.TwitterApi struct
type API struct {
	*anaconda.TwitterApi
}

// createTwitterAPI groups functions
func createTwitterAPI() *API {
	currentUser := getCurrentUser()
	configPath := getUserConfig(currentUser)
	configFile := openFile(configPath)
	credentials := readCredentials(configFile)
	twitterAPI := getTwitterAPI(credentials)
	return &API{twitterAPI}
}

func (api *API) searchTweets(searchKey string) {
	tweets, err := api.GetSearch(searchKey, nil)
	terminateOnError(err)
	for _, tweet := range tweets.Statuses {
		fmt.Println(tweet.Text)
		fmt.Println("")
	}
}

func (api *API) listFollowers() {
	pages := api.GetFollowersListAll(nil)
	for page := range pages {
		fmt.Println(page.Followers)
		fmt.Println("------------------------------")
	}
}

func main() {
	opt := flag.String("opt", "", "")
	opt := os.Args[1]
	param := os.Args[2]
	api := createTwitterAPI()
	switch opt {
	case "s":
		api.searchTweets(param)
	case "f":
		api.listFollowers()
	}
}

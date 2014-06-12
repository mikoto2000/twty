package main

import (
	"./twitter"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func getConfig() (string, map[string]string) {
	home := os.Getenv("HOME")
	dir := filepath.Join(home, ".config")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(home, "Application Data")
		}
	} else if runtime.GOOS == "plan9" {
		home = os.Getenv("home")
		dir = filepath.Join(home, ".config")
	}
	_, err := os.Stat(dir)
	if err != nil {
		if os.Mkdir(dir, 0700) != nil {
			log.Fatal("failed to create directory:", err)
		}
	}
	dir = filepath.Join(dir, "twty")
	_, err = os.Stat(dir)
	if err != nil {
		if os.Mkdir(dir, 0700) != nil {
			log.Fatal("failed to create directory:", err)
		}
	}
	file := filepath.Join(dir, "settings.json")
	config := map[string]string{}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		config["ClientToken"] = "eAw30zFQxWg7tb8NSEmOGyR44"
		config["ClientSecret"] = "Qpt7o3kMgp0Ca8fDevbUHtrBWTMpdKdFCkWTZf1Zeu9KFdFPBW"
	} else {
		err = json.Unmarshal(b, &config)
		if err != nil {
			log.Fatal("could not unmarhal settings.json:", err)
		}
	}
	return file, config
}

func main() {
	//	reply := flag.Bool("r", false, "show replies")
	//	list := flag.String("l", "", "show tweets")
	//	user := flag.String("u", "", "show user timeline")
	//	favorite := flag.String("f", "", "specify favorite ID")
	//	search := flag.String("s", "", "search word")
	//	inreply := flag.String("i", "", "specify in-reply ID, if not specify text, it will be RT.")
	verbose := flag.Bool("v", false, "detail display")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage of twty:
  -f ID: specify favorite ID
  -i ID: specify in-reply ID, if not specify text, it will be RT.
  -l USER/LIST: show list's timeline (ex: mattn_jp/subtech)
  -u USER: show user's timeline
  -s WORD: search timeline
  -r: show replies
  -v: detail display
`)
	}
	flag.Parse()

	file, config := getConfig()
	clientToken := config["ClientToken"]
	clientSecret := config["ClientSecret"]
	accessToken, foundAccessToken := config["AccessToken"]
	accessSecret, foundAccessSecret := config["AccessSecret"]

	var twClient *twitter.Twitter
	if foundAccessToken && foundAccessSecret {
		twClientTmp := twitter.NewTwitterFromAccessInfo(accessToken, accessSecret, clientToken, clientSecret)
		twClient = twClientTmp
	} else {
		twClientTmp, authorized, err := twitter.NewTwitterFromClientInfo(clientToken, clientSecret)
		if err != nil {
			log.Fatal("faild to create twitter:", err)
		}
		if authorized {
			config["AccessToken"] = twClientTmp.Token.Token
			config["AccessSecret"] = twClientTmp.Token.Secret
			b, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				log.Fatal("failed to store file:", err)
			}
			err = ioutil.WriteFile(file, b, 0700)
			if err != nil {
				log.Fatal("failed to store file:", err)
			}
		}
		twClient = twClientTmp
	}

	//	if len(*search) > 0 {
	//		tweets, err := twClient.GetStatuses(token, "https://api.twitter.com/1.1/search/tweets.json", map[string]string{"q": *search})
	//		if err != nil {
	//			log.Fatal("failed to get tweets:", err)
	//		}
	//		showTweets(tweets, *verbose)
	//	} else if *reply {
	//		tweets, err := twClient.GetTweets(token, "https://api.twitter.com/1.1/statuses/mentions_timeline.json", map[string]string{})
	//		if err != nil {
	//			log.Fatal("failed to get tweets:", err)
	//		}
	//		showTweets(tweets, *verbose)
	//	} else if len(*list) > 0 {
	//		part := strings.SplitN(*list, "/", 2)
	//		if len(part) == 2 {
	//			tweets, err := twClient.GetTweets(token, "https://api.twitter.com/1.1/lists/statuses.json", map[string]string{"owner_screen_name": part[0], "slug": part[1]})
	//			if err != nil {
	//				log.Fatal("failed to get tweets:", err)
	//			}
	//			showTweets(tweets, *verbose)
	//		}
	//	} else if len(*user) > 0 {
	//		tweets, err := twClient.GetTweets(token, "https://api.twitter.com/1.1/statuses/user_timeline.json", map[string]string{"screen_name": *user})
	//		if err != nil {
	//			log.Fatal("failed to get tweets:", err)
	//		}
	//		showTweets(tweets, *verbose)
	//	} else if len(*favorite) > 0 {
	//		twClient.PostTweet(token, "https://api.twitter.com/1.1/favorites/create.json", map[string]string{"id": *favorite})
	//	} else if flag.NArg() == 0 {
	//		if len(*inreply) > 0 {
	//			twClient.PostTweet(token, "https://api.twitter.com/1.1/statuses/retweet/"+*inreply+".json", map[string]string{})
	//		} else {
	tweets, err := twClient.GetHomeTimeline()
	if err != nil {
		log.Fatal("failed to get tweets:", err)
	}
	showTweets(tweets, *verbose)
	//		}
	//	} else {
	//		twClient.PostTweet(token, "https://api.twitter.com/1.1/statuses/update.json", map[string]string{"status": strings.Join(flag.Args(), " "), "in_reply_to_status_id": *inreply})
	//	}
}

func showTweets(tweets []twitter.Tweet, verbose bool) {
	if verbose {
		for i := len(tweets) - 1; i >= 0; i-- {
			name := tweets[i].User.Name
			user := tweets[i].User.ScreenName
			text := tweets[i].Text
			text = strings.Replace(text, "\r", "", -1)
			text = strings.Replace(text, "\n", " ", -1)
			text = strings.Replace(text, "\t", " ", -1)
			fmt.Println(user + ": " + name)
			fmt.Println("  " + text)
			fmt.Println("  " + tweets[i].Identifier)
			fmt.Println("  " + tweets[i].CreatedAt)
			fmt.Println()
		}
	} else {
		for i := len(tweets) - 1; i >= 0; i-- {
			user := tweets[i].User.ScreenName
			text := tweets[i].Text
			fmt.Println(user + ": " + text)
		}
	}
}

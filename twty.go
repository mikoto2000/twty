package main

import (
	"./myoauth"
	"./tweet"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	reply := flag.Bool("r", false, "show replies")
	list := flag.String("l", "", "show tweets")
	user := flag.String("u", "", "show user timeline")
	favorite := flag.String("f", "", "specify favorite ID")
	search := flag.String("s", "", "search word")
	stream := flag.Bool("S", false, "stream timeline")
	inreply := flag.String("i", "", "specify in-reply ID, if not specify text, it will be RT.")
	verbose := flag.Bool("v", false, "detail display")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage of twty:
  -f ID: specify favorite ID
  -i ID: specify in-reply ID, if not specify text, it will be RT.
  -l USER/LIST: show list's timeline (ex: mattn_jp/subtech)
  -u USER: show user's timeline
  -s WORD: search timeline
  -S: stream timeline
  -r: show replies
  -v: detail display
`)
	}
	flag.Parse()

	file, config := getConfig()
	token, authorized, err := myoauth.GetAccessToken(config)
	if err != nil {
		log.Fatal("faild to get access token:", err)
	}
	if authorized {
		b, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal("failed to store file:", err)
		}
		err = ioutil.WriteFile(file, b, 0700)
		if err != nil {
			log.Fatal("failed to store file:", err)
		}
	}

	if len(*search) > 0 {
		tweets, err := tweet.GetStatuses(token, "https://api.twitter.com/1.1/search/tweets.json", map[string]string{"q": *search})
		if err != nil {
			log.Fatal("failed to get tweets:", err)
		}
		tweet.ShowTweets(tweets, *verbose)
	} else if *reply {
		tweets, err := tweet.GetTweets(token, "https://api.twitter.com/1.1/statuses/mentions_timeline.json", map[string]string{})
		if err != nil {
			log.Fatal("failed to get tweets:", err)
		}
		tweet.ShowTweets(tweets, *verbose)
	} else if len(*list) > 0 {
		part := strings.SplitN(*list, "/", 2)
		if len(part) == 2 {
			tweets, err := tweet.GetTweets(token, "https://api.twitter.com/1.1/lists/statuses.json", map[string]string{"owner_screen_name": part[0], "slug": part[1]})
			if err != nil {
				log.Fatal("failed to get tweets:", err)
			}
			tweet.ShowTweets(tweets, *verbose)
		}
	} else if len(*user) > 0 {
		tweets, err := tweet.GetTweets(token, "https://api.twitter.com/1.1/statuses/user_timeline.json", map[string]string{"screen_name": *user})
		if err != nil {
			log.Fatal("failed to get tweets:", err)
		}
		tweet.ShowTweets(tweets, *verbose)
	} else if len(*favorite) > 0 {
		tweet.PostTweet(token, "https://api.twitter.com/1.1/favorites/create.json", map[string]string{"id": *favorite})
	} else if *stream {
		url_ := "https://userstream.twitter.com/1.1/user.json"
		param := make(url.Values)
		myoauth.OauthClient.SignParam(token, "GET", url_, param)
		url_ = url_ + "?" + param.Encode()
		res, err := http.Get(url_)
		if err != nil {
			log.Fatal("failed to get tweets:", err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatal("failed to get tweets:", err)
		}
		buf := bufio.NewReader(res.Body)
		var last []byte
		for {
			b, _, err := buf.ReadLine()
			last = append(last, b...)
			var tweets [1]tweet.Tweet
			err = json.Unmarshal(last, &tweets[0])
			if err != nil {
				continue
			}
			last = []byte{}
			tweet.ShowTweets(tweets[:], *verbose)
		}
	} else if flag.NArg() == 0 {
		if len(*inreply) > 0 {
			tweet.PostTweet(token, "https://api.twitter.com/1.1/statuses/retweet/"+*inreply+".json", map[string]string{})
		} else {
			tweets, err := tweet.GetTweets(token, "https://api.twitter.com/1.1/statuses/home_timeline.json", map[string]string{})
			if err != nil {
				log.Fatal("failed to get tweets:", err)
			}
			tweet.ShowTweets(tweets, *verbose)
		}
	} else {
		tweet.PostTweet(token, "https://api.twitter.com/1.1/statuses/update.json", map[string]string{"status": strings.Join(flag.Args(), " "), "in_reply_to_status_id": *inreply})
	}
}

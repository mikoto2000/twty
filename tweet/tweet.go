package tweet

import (
	"../myoauth"
	"github.com/garyburd/go-oauth/oauth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Tweet struct {
	Text       string
	Identifier string `json:"id_str"`
	Source     string
	CreatedAt  string `json:"created_at"`
	User       struct {
		Name            string
		ScreenName      string `json:"screen_name"`
		FollowersCount  int    `json:"followers_count"`
		ProfileImageURL string `json:"profile_image_url"`
	}
	Place *struct {
		Id       string
		FullName string `json:"full_name"`
	}
	Entities struct {
		HashTags []struct {
			Indices [2]int
			Text    string
		}
		UserMentions []struct {
			Indices    [2]int
			ScreenName string `json:"screen_name"`
		} `json:"user_mentions"`
		Urls []struct {
			Indices [2]int
			Url     string
		}
	}
}

type RSS struct {
	Channel struct {
		Title       string
		Description string
		Link        string
		Item        []struct {
			Title       string
			Description string
			PubDate     string
			Link        []string
			Guid        string
			Author      string
		}
	}
}

func GetTweets(token *oauth.Credentials, url_ string, opt map[string]string) ([]Tweet, error) {
	param := make(url.Values)
	for k, v := range opt {
		param.Set(k, v)
	}
	myoauth.OauthClient.SignParam(token, "GET", url_, param)
	url_ = url_ + "?" + param.Encode()
	res, err := http.Get(url_)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	var tweets []Tweet
	err = json.NewDecoder(res.Body).Decode(&tweets)
	if err != nil {
		return nil, err
	}
	return tweets, nil
}

func GetStatuses(token *oauth.Credentials, url_ string, opt map[string]string) ([]Tweet, error) {
	param := make(url.Values)
	for k, v := range opt {
		param.Set(k, v)
	}
	myoauth.OauthClient.SignParam(token, "GET", url_, param)
	url_ = url_ + "?" + param.Encode()
	res, err := http.Get(url_)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	var statuses struct {
		Statuses []Tweet
	}
	err = json.NewDecoder(res.Body).Decode(&statuses)
	if err != nil {
		return nil, err
	}
	return statuses.Statuses, nil
}

func ShowTweets(tweets []Tweet, verbose bool) {
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

func PostTweet(token *oauth.Credentials, url_ string, opt map[string]string) error {
	param := make(url.Values)
	for k, v := range opt {
		param.Set(k, v)
	}
	myoauth.OauthClient.SignParam(token, "POST", url_, param)
	res, err := http.PostForm(url_, url.Values(param))
	if err != nil {
		log.Println("failed to post tweet:", err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("failed to get timeline:", err)
		return err
	}
	return nil
}


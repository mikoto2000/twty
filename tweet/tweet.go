package tweet

import (
	"../myoauth"
	"encoding/json"
	"github.com/garyburd/go-oauth/oauth"
	"log"
	"net/http"
	"net/url"
)

type Twitter struct {
	Token *oauth.Credentials
}

func NewTwitterFromClientInfo(clientToken string, clientSecret string) (*Twitter, bool, error) {

	token, authorized, err := myoauth.GetAccessToken(clientToken, clientSecret)

	var funcError error

	if err != nil {
		funcError = err
	}

	var twitter *Twitter
	if funcError == nil && authorized {
		twitter = new(Twitter)
		twitter.Token = token
	} else {
		twitter = nil
	}

	return twitter, authorized, funcError
}

func (twitter *Twitter) GetHomeTimeline() ([]Tweet, error) {
	tweets, err := getTweets(twitter.Token, "https://api.twitter.com/1.1/statuses/home_timeline.json", map[string]string{})

	return tweets, err
}

func NewTwitterFromAccessInfo(accessToken string, accessSecret string, clientToken string, clientSecret string) *Twitter {

	token := myoauth.NewAccessToken(clientToken, clientSecret, accessToken, accessSecret)

	twitter := new(Twitter)
	twitter.Token = token
	return twitter
}

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

func getTweets(token *oauth.Credentials, url_ string, opt map[string]string) ([]Tweet, error) {
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

package myoauth

import (
	"bufio"
	"github.com/garyburd/go-oauth/oauth"
	"log"
	"net/http"
	"os"
    "fmt"
)

var OauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

func clientAuth(requestToken *oauth.Credentials) (*oauth.Credentials, error) {
	url_ := OauthClient.AuthorizationURL(requestToken, nil)

	fmt.Println("url ", url_)

	print("PIN: ")
	stdin := bufio.NewReader(os.Stdin)
	b, err := stdin.ReadBytes('\n')
	if err != nil {
		log.Fatal("canceled")
	}

	if b[len(b)-2] == '\r' {
		b = b[0 : len(b)-2]
	} else {
		b = b[0 : len(b)-1]
	}
	accessToken, _, err := OauthClient.RequestToken(http.DefaultClient, requestToken, string(b))
	if err != nil {
		log.Fatal("failed to request token:", err)
	}
	return accessToken, nil
}

func NewAccessToken(clientToken string, clientSecret string, accessToken string, accessSecret string)  (*oauth.Credentials) {
	OauthClient.Credentials.Token = clientToken
	OauthClient.Credentials.Secret = clientSecret

	return &oauth.Credentials{accessToken, accessSecret}
}

func GetAccessToken(clientToken string, clientSecret string) (*oauth.Credentials, bool, error) {
	OauthClient.Credentials.Token = clientToken
	OauthClient.Credentials.Secret = clientSecret

	authorized := false
	var token *oauth.Credentials
	requestToken, err := OauthClient.RequestTemporaryCredentials(http.DefaultClient, "", nil)
	if err != nil {
		log.Print("failed to request temporary credentials:", err)
		return nil, false, err
	}
	token, err = clientAuth(requestToken)
	if err != nil {
		log.Print("failed to request temporary credentials:", err)
		return nil, false, err
	}

	authorized = true

	return token, authorized, nil
}


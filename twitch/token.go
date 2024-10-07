package twitch

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type response struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

func accessToken() string {
	var form = url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{os.Getenv("TWITCH_REFRESH")},
		"client_id":     []string{os.Getenv("TWITCH_CLIENT")},
		"client_secret": []string{os.Getenv("TWITCH_SECRET")},
	}
	if resp, e := http.PostForm("https://id.twitch.tv/oauth2/token", form); e != nil {
		log.Printf("ERROR: %q\n", e)
		return ""
	} else {
		defer resp.Body.Close()
		var r = response{}
		var d, e = io.ReadAll(resp.Body)
		if e != nil {
			log.Printf("ERROR: %q\n", e)
			return ""
		}
		if e := json.Unmarshal(d, &r); e != nil {
			log.Printf("ERROR: %q\n", e)
			return ""
		}
		return r.AccessToken
	}
}

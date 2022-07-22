// Package twitter provides the client and interface to interact with the twitter api.
package twitter

import (
	"encoding/json"
	"fmt"
	"github.com/jhidalgoesp/example-services/internal/oauth"
	"go.uber.org/zap"
	"io/ioutil"
)

var timelineUrl = "https://api.twitter.com/1.1/statuses/user_timeline.json"

// Tweet represents a twitter tweet.
type Tweet struct {
	Id        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
	User      User   `json:"user"`
}

// User represents an user.
type User struct {
	Name         string `json:"name"`
	ScreenName   string `json:"screen_name"`
	ProfileImage string `json:"profile_image_url"`
}

// Api describes twitter api integration.
type Api interface {
	GetLastTweets(username, count string, tweets *[]Tweet) error
}

// Client manages the set of APIs for twitter access.
type Client struct {
	Log *zap.SugaredLogger
}

// GetLastTweets retrieves the last n getTweets from an user.
func (c Client) GetLastTweets(username, count string, tweets *[]Tweet) error {
	httpClient := oauth.InitOauthClient()
	resp, _ := httpClient.Get(fmt.Sprintf("%s?screen_name=%s&count=%s", timelineUrl, username, count))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(body, &tweets)
	if err != nil {
		return err
	}

	return nil
}

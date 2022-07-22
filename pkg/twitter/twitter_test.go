package twitter

import (
	"fmt"
	"github.com/jhidalgoesp/example-services/tests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTwitter_GetLastTweets(t *testing.T) {
	t.Run("GetLastTweets should return an array of tweets", func(t *testing.T) {
		expected := "[{\"id\":123,\"created_at\":\"Sun Nov 08 23:12:45 +0000 2020\",\"text\":\"this is my test from Java\",\"user\":{\"name\":\"Wilsonnn\",\"screen_name\":\"backend_test\",\"profile_image_url\":\"http://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png\"}}]"
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expected)
		}))
		defer svr.Close()
		timelineUrl = svr.URL
		client := Client{nil}
		var tweets []Tweet

		err := client.GetLastTweets("test", "1", &tweets)

		tests.AssertNilError(t, err)
		tests.AssertSameInt(t, tweets[0].Id, 123)
		tests.AssertStrings(t, tweets[0].Text, "this is my test from Java")
		tests.AssertStrings(t, tweets[0].User.Name, "Wilsonnn")
	})
}

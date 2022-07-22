package oauth

import (
	"github.com/dghubble/oauth1"
	"net/http"
	"os"
)

func InitOauthClient() *http.Client {
	config := oauth1.NewConfig(os.Getenv("consumerKey"), os.Getenv("consumerSecret"))
	token := oauth1.NewToken(os.Getenv("token"), os.Getenv("tokenSecret"))

	return config.Client(oauth1.NoContext, token)
}

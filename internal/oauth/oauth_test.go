package oauth

import (
	"github.com/jhidalgoesp/example-services/tests"
	"testing"
)

func TestInitOauthClient(t *testing.T) {
	t.Run("InitOauthClient should return a pointer", func(t *testing.T) {
		client := InitOauthClient()
		tests.AssertKindIsPointer(t, client)
	})
}

package awsProvider

import (
	"github.com/jhidalgoesp/example-services/tests"
	"testing"
)

func TestInitDynamoClient(t *testing.T) {
	t.Run("InitDynamoClient returns a pointer", func(t *testing.T) {
		dynamoClient := NewProvider().NewDynamoClient()
		tests.AssertKindIsPointer(t, dynamoClient)
	})
}

func TestInitSSMClient(t *testing.T) {
	t.Run("InitSSMClient returns a pointer", func(t *testing.T) {
		ssmClient := NewProvider().NewSSMClient()
		tests.AssertKindIsPointer(t, ssmClient)
	})
}

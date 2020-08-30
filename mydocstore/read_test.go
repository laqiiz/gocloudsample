package mydocstore

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gocloud.dev/docstore"
	"gocloud.dev/docstore/awsdynamodb"
	"gocloud.dev/gcerrors"
	"log"
	"testing"
)

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "local",
	}))

	// Create DynamoDB client
	db = dynamodb.New(sess, &aws.Config{Endpoint: aws.String("http://localhost:8000")})

	// docstore
	myTableColl *docstore.Collection
)

func init() {
	// Open Collection via Go CDK
	coll, err := awsdynamodb.OpenCollection(db, "MyTable", "MyHashKey", "MyRangeKey", nil)
	if err != nil {
		log.Fatal(err)
	}
	myTableColl = coll
}

func TestBatchGetNotFound(t *testing.T) {

	t.Cleanup(func() {
		_ = myTableColl.Close()
	})

	gets := []Item{
		{
			MyHashKey:  "00001",
			MyRangeKey: 1,
		},
		{
			MyHashKey:  "00001",
			MyRangeKey: 2,
		},
	}

	actions := myTableColl.Actions()
	for _, v := range gets {
		v := v
		actions = actions.Get(&v)
	}

	if err := actions.Do(context.Background()); err != nil {
		if aerrs, ok := err.(docstore.ActionListError); ok {
			for _, aerr := range aerrs {
				if gcerrors.Code(aerr.Err) == gcerrors.NotFound {
					//read_test.go:55: get NotFound err: item {0xc00041c270 map[] {0xc002a0 0xc00041c270 409} [{MyHashKey true 0xb64ac0 [0] {false} [77 121 72 97 115 104 75 101 121] 0x9e3ce0} {MyRangeKey true 0xb640c0 [1] {false} [77 121 82 97 110 103 101 75 101 121] 0x9e3ce0} {MyText true 0xb64ac0 [2] {false} [77 121 84 101 120 116] 0x9e3f40}]} not found (code=NotFound)
					//read_test.go:55: get NotFound err: item {0xc00041c2a0 map[] {0xc002a0 0xc00041c2a0 409} [{MyHashKey true 0xb64ac0 [0] {false} [77 121 72 97 115 104 75 101 121] 0x9e3ce0} {MyRangeKey true 0xb640c0 [1] {false} [77 121 82 97 110 103 101 75 101 121] 0x9e3ce0} {MyText true 0xb64ac0 [2] {false} [77 121 84 101 120 116] 0x9e3f40}]} not found (code=NotFound)
					t.Logf("get NotFound err: %v", aerr.Err)
				} else {
					t.Errorf("unknown err: %v", aerr.Err)
				}
			}
		}
	}

}

func BatchGet(keys ...interface{}) ([]bool, error) {
	resp := make([]bool, len(keys))

	actions := myTableColl.Actions()
	for _, v := range keys {
		v := v
		actions = actions.Get(&v)
	}
	if err := actions.Do(context.Background()); err != nil {
		if aerrs, ok := err.(docstore.ActionListError); ok {
			for _, aerr := range aerrs {
				if gcerrors.Code(aerr.Err) == gcerrors.NotFound {
					resp[aerr.Index] = false
					continue
				}

			}
		}
	}

	return nil

}

package pmongo

import (
	"context"
	"fmt"
	"testing"
)

func TestMongoCommand(t *testing.T) {
	MyClient()
}

func MyClient() {
	cfg := Config{}
	client, _ := cfg.GetMongoClient()
	DBOps(client)
}

func DBOps(db DB) {
	ctx := context.Background()
	rows := make([]map[string]interface{}, 0)
	err := db.Table("test").Find(map[string]interface{}{}).Sort("create_time").Limit(20).All(ctx, &rows)
	if err != nil {
		fmt.Println(err)
	}
}

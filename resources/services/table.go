package services

import (
	"context"

	"github.com/aqeel/cq-source-duo/client"
	"github.com/aqeel/cq-source-duo/client/admin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
)

func UsersTable() *schema.Table {
	return &schema.Table{
		Name:      "duo_users_table",
		Resolver:  fetchUsersTable,
		Transform: transformers.TransformWithStruct(&admin.User{}),
	}
}

func fetchUsersTable(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	maxCounter := 100
	offset := 1
	limit := 100
	cl := meta.(*client.Client)
	for i := maxCounter; i >= 0; i-- {
		result, err := cl.GetUsers(ctx, offset, limit)
		if err != nil {
			return err
		}
		res <- result
		if len(result) == 0 {
			break
		}
	}
	return nil
}

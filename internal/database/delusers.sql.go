// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: delusers.sql

package database

import (
	"context"
)

const delUsers = `-- name: DelUsers :exec
DELETE FROM users
`

func (q *Queries) DelUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, delUsers)
	return err
}

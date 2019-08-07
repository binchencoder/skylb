package db

import (
	"upper.io/db.v3/lib/sqlbuilder"
)

func withTx(tx sqlbuilder.Tx, fn func() error) error {
	if err := fn(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

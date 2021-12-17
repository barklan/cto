package db

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/barklan/cto/pkg/postgres/models"
	_ "github.com/lib/pq"
)

func TestClient(t *testing.T) {
	tx := dbx.MustBegin()
	tx.MustExec("INSERT INTO client (tg_nick) VALUES ($1)", "barklan")
	tx.MustExec("INSERT INTO client (tg_nick) VALUES ($1)", "johndoe")

	err := tx.Commit()
	if err != nil {
		panic("failed to commit transaction")
	}

	clients := []models.Client{}

	statement, _, err := sq.Select("*").From("client").ToSql()
	if err != nil {
		panic("sql generation failed")
	}

	err = dbx.Select(&clients, statement)
	if err != nil {
		panic(err)
	}

	barklan, _ := clients[0], clients[1]
	if barklan.TGNick != "barklan" {
		t.Errorf("Got %v want %v", barklan.TGNick, "barklan")
	}
}

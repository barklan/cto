package models

import "database/sql"

type Client struct {
	ID     string `db:"id"`
	Active bool   `db:"active"`
	TGNick string `db:"tg_nick"`
}

type Project struct {
	ID          string         `db:"id"`
	Active      bool           `db:"active"`
	ClientID    string         `db:"client_id"`
	PrettyTitle sql.NullString `db:"pretty_title"`
	SecretKey   string         `db:"secret_key"`
}

type Chat struct {
	ID        string `db:"id"`
	Active    bool   `db:"active"`
	ProjectID string `db:"project_id"`
}

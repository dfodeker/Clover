package database

import (
	"database/sql"
	"fmt"

	errorsx "github.com/dfodeker/clover/errors"
	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	db *sql.DB
}

func NewClient(pathToDB string) (Client, error) {
	db, err := sql.Open("sqlite3", pathToDB)
	if err != nil {
		return Client{}, err
	}
	c := Client{db}
	err = c.autoMigrate()
	if err != nil {
		return Client{}, err
	}
	return c, nil

}

// type Note struct {
// 	ID        uuid.UUID `db:"id"         json:"id"`
// 	Key       string    `db:"key"        json:"key"`
// 	Title     string    `db:"title"      json:"title"`
// 	Body      string    `db:"body"       json:"body"`
// 	Tags      []string  `db:"-"          json:"tags"` // slice in memory
// 	CreatedAt time.Time `db:"created_at" json:"created_at"`
// 	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
// 	Checksum  string    `db:"checksum"   json:"checksum"`
// }

func (c *Client) autoMigrate() error {

	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		password TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
);`
	_, err := c.db.Exec(userTable)
	if err != nil {

		return errorsx.Wrap(err, "error occured creating user table")
	}
	refreshTokenTable := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		token TEXT PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		revoked_at TIMESTAMP,
		user_id TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	_, err = c.db.Exec(refreshTokenTable)
	if err != nil {

		return errorsx.Wrap(err, "error occured creating refresh_tokens table")
	}

	notesTable := `
	CREATE TABLE IF NOT EXISTS notes (
		id TEXT PRIMARY KEY,
		slug NOT NULL UNIQUE,
		title TEXT NOT NULL,
		body TEXT, 
		tags TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		checksum TEXT NOT NULL
	);
	`
	_, err = c.db.Exec(notesTable)
	if err != nil {
		return errorsx.Wrap(err, "error occured creating notes table")
	}
	tagsTable := `
	CREATE TABLE IF NOT EXISTS tags(
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = c.db.Exec(tagsTable)
	if err != nil {
		return errorsx.Wrap(err, " error occured creating tags table")
	}
	noteTagTable := `
	CREATE TABLE IF NOT EXISTS note_tags (
  note_id INTEGER NOT NULL,
  tag_id  INTEGER NOT NULL,
  PRIMARY KEY (note_id, tag_id),            
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id)  REFERENCES tags(id)  ON DELETE CASCADE
);
	`
	_, err = c.db.Exec(noteTagTable)
	if err != nil {
		return errorsx.Wrap(err, "error occured cretign noteTag table")
	}

	noteIndex := `CREATE INDEX IF NOT EXISTS idx_note_tags_tag ON note_tags(tag_id);`
	_, err = c.db.Exec(noteIndex)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) Reset() error {
	if _, err := c.db.Exec("DELETE FROM refresh_tokens"); err != nil {
		return fmt.Errorf("failed to reset table refresh_tokens: %w", err)
	}
	if _, err := c.db.Exec("DELETE FROM users"); err != nil {
		return fmt.Errorf("failed to reset table users: %w", err)
	}
	if _, err := c.db.Exec("DELETE FROM notes"); err != nil {
		return fmt.Errorf("failed to reset table notes: %w", err)
	}
	if _, err := c.db.Exec("DELETE FROM tags"); err != nil {
		return fmt.Errorf("failed to reset table tags: %w", err)
	}
	if _, err := c.db.Exec("DELETE FROM note_tags"); err != nil {
		return fmt.Errorf("failed to reset table note_tags: %w", err)
	}
	if _, err := c.db.Exec("DROP INDEX IF EXISTS idx_note_tags_tag;"); err != nil {
		return fmt.Errorf("failed to reset note_tags index: %w", err)
	}
	return nil
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	Key       string    `db:"key"        json:"key"`
	Title     string    `db:"title"      json:"title"`
	Body      string    `db:"body"       json:"body"`
	Tags      []string  `db:"-"          json:"tags"` // slice in memory
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Checksum  string    `db:"checksum"   json:"checksum"`
}

func NewNote(title, body string, tags []string) (Note, error) {
	id := uuid.New()
	tags = normalizeTags(tags)

	n := Note{
		ID:        id,
		Title:     strings.TrimSpace(title),
		Body:      body,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	n.Checksum = checksumNote(n.Title, n.Body)
	return n, nil
}

func normalizeTags(ts []string) []string {
	seen := make(map[string]struct{}, len(ts))
	out := make([]string, 0, len(ts))
	for _, t := range ts {
		t = strings.TrimSpace(strings.ToLower(t))
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	return out
}

func checksumNote(title, body string) string {
	// normalize body minimally: LF newlines, trim trailing spaces
	b := strings.ReplaceAll(body, "\r\n", "\n")
	lines := strings.Split(b, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	b = strings.Join(lines, "\n")
	h := sha256.Sum256([]byte(title + "\n\n" + b))
	return hex.EncodeToString(h[:])
}

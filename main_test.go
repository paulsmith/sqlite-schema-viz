package main

import (
	"bytes"
	"database/sql"
	"flag"
	"os"
	"path"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestRender(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	schema := `
CREATE TABLE users ( id integer, name text );
CREATE TABLE posts ( id integer, user_id integer, title text, body text, foreign key (user_id) references users(id) );
CREATE TABLE comments ( id integer, post_id integer, body text, foreign key (post_id) references posts(id) );
`

	if _, err := db.Exec(schema); err != nil {
		t.Fatal(err)
	}

	actual, err := render(db, SVG)
	if err != nil {
		t.Fatal(err)
	}

	goldenPath := path.Join("testdata", t.Name()+".golden.svg")

	if *update {
		if err := os.WriteFile(goldenPath, actual, 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Compare to golden image
	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		// Produce a test error, but don't display the raw image data
		t.Error("rendered diagram does not match golden image")
	}
}

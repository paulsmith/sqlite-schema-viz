/*
 * This program generates a graph visualization of a SQLite database schema.
 * It uses the graphviz library to generate the graph.
 *
 * Copyright 2024 Paul Smith <paulsmith@pobox.com>
 */
package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"log"
	"strings"

	"github.com/goccy/go-graphviz"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed sqlite-schema-diagram/sqlite-schema-diagram.sql
var schemaDiagramSQL string

func main() {
	// The required argument is the path to the SQLite database file.
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("missing required argument: path to SQLite database file")
	}
	dbPath := flag.Arg(0)

	// Open the database.
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Execute the SQL query to get the schema diagram as a Graphviz DOT string
	rows, err := db.Query(schemaDiagramSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Get all the rows combined into a single source string
	var source strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			log.Fatal(err)
		}
		source.WriteString(line)
	}

	// Using the go-graphviz library, parse the DOT string into a graph and render it to a PNG file.
	graph, err := graphviz.ParseBytes([]byte(source.String()))
	if err != nil {
		log.Fatal(err)
	}

	g := graphviz.New()
	if err := g.RenderFilename(graph, graphviz.PNG, "schema.png"); err != nil {
		log.Fatal(err)
	}
}

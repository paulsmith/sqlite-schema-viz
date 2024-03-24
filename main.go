/*
 * This program generates a graph visualization of a SQLite database schema.
 * It uses the graphviz library to generate the graph.
 *
 * Copyright 2024 Paul Smith <paulsmith@pobox.com>
 * Credit to @Screwtapello on Gitlab for the sqlite-schema-diagram.sql file.
 */
package main

import (
	"bytes"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-graphviz"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed sqlite-schema-diagram/sqlite-schema-diagram.sql
var schemaDiagramSQL string

func usage() {
	log.Fatalf("Usage: %s <path-to-sqlite-database>\n", os.Args[0])
}

func main() {
	var outputImgPath string
	flag.StringVar(&outputImgPath, "o", "schema.png", "Output image file path")

	// The required argument is the path to the SQLite database file.
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}
	dbPath := flag.Arg(0)

	// Open the database.
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	img, err := render(db, PNG)
	if err != nil {
		log.Fatal(err)
	}

	// Write the image to the output file.
	if err := os.WriteFile(outputImgPath, img, 0644); err != nil {
		log.Fatal(err)
	}
}

type format int

const (
	PNG format = iota
	SVG
)

func (f format) GV() graphviz.Format {
	switch f {
	case PNG:
		return graphviz.PNG
	case SVG:
		return graphviz.SVG
	default:
		panic(fmt.Sprintf("unknown format: %v", f))
	}
}

// render produces an in-memory image of the SQLite database schema graph
func render(db *sql.DB, format format) ([]byte, error) {
	// Execute the SQL query to get the schema diagram as a Graphviz DOT string
	rows, err := db.Query(schemaDiagramSQL)
	if err != nil {
		return nil, fmt.Errorf("querying schema diagram: %w", err)
	}
	defer rows.Close()

	// Get all the rows combined into a single source string
	var source strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			log.Fatal(err)
		}
		source.WriteString(line + "\n")
	}

	// Using the go-graphviz library, parse the DOT string into a graph and render it to a PNG file.
	graph, err := graphviz.ParseBytes([]byte(source.String()))
	if err != nil {
		return nil, fmt.Errorf("parsing DOT: %w", err)
	}

	g := graphviz.New()
	var img bytes.Buffer
	if err := g.Render(graph, format.GV(), &img); err != nil {
		return nil, fmt.Errorf("rendering graph: %w", err)
	}

	return img.Bytes(), nil
}

sqlite-schema-viz:
	go build .

.PHONY: sqlite-schema-viz

example.png: sqlite-schema-viz
	grep CREATE README.md | sqlite3 example.db
	./sqlite-schema-viz -o example.png example.db
	rm example.db

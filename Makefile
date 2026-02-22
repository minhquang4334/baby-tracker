.PHONY: all frontend backend run frontend-watch clean deps

BINARY := baby-care
FRONTEND_DIR := frontend
STATIC_DIR := static

all: deps frontend backend

deps:
	cd $(FRONTEND_DIR) && npm install
	go mod download

frontend:
	cd $(FRONTEND_DIR) && npx esbuild src/main.ts \
		--bundle \
		--outfile=../$(STATIC_DIR)/app.js \
		--minify \
		--loader:.css=css \
		--external:*.ttf \
		--external:*.woff \
		--external:*.woff2
	cp $(FRONTEND_DIR)/index.html $(STATIC_DIR)/index.html

frontend-watch:
	cd $(FRONTEND_DIR) && npx esbuild src/main.ts \
		--bundle \
		--outfile=../$(STATIC_DIR)/app.js \
		--watch \
		--loader:.css=css \
		--external:*.ttf \
		--external:*.woff \
		--external:*.woff2

backend:
	go build -o $(BINARY) .

run:
	go run . --port 8080

clean:
	rm -f $(BINARY)
	rm -f $(STATIC_DIR)/app.js $(STATIC_DIR)/app.css $(STATIC_DIR)/index.html

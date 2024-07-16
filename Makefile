release: testAll clean init gen_tags build

init:
	go mod tidy
	go mod download

clean:
	go clean
	go clean -testcache
	rm -rf bin/
	rm -rf public/
	rm -rf tmp/
	rm tags || true

build:
	mkdir -p bin/
	go build -o ./bin/solver ./cmd/solver/main.go

testAll:
	go test -v ./...

test_with_coverage:
	mkdir -p tmp/
	go test -coverprofile ./tmp/cover.out ./...
	go tool cover -html=./tmp/cover.out -o ./tmp/cover.html
	go-cover-treemap -coverprofile tmp/cover.out > tmp/out.svg
	open ./tmp/out.svg
	open tmp/cover.html

gen_tags:
	ctags -R --exclude=.git --exclude=cache

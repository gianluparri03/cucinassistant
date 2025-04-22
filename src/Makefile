run:
	CA_ENV=development go run .

test:
	CA_ENV=testing go test -v cucinassistant/database cucinassistant/langs

cover:
	CA_ENV=testing go test -coverprofile=cover.out -covermode atomic cucinassistant/database
	go tool cover -html=cover.out

fmt:
	go fmt ./...

gen:
	go generate ./...
	templ generate

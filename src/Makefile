# Install all the dependencies
install:
	go mod download
	go install golang.org/x/tools/cmd/stringer@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go get -u github.com/a-h/templ/cmd/templ@latest

# Runs code-generation tools
gen:
	go generate ./...
	templ generate


# Starts CucinAssistant
run:
	CA_ENV=development go run .


# Runs the tests
test:
	CA_ENV=testing go test -v cucinassistant/database cucinassistant/langs

# Runs the tests and shows the coverage
cover:
	CA_ENV=testing go test -coverprofile=cover.out -covermode atomic cucinassistant/database
	go tool cover -html=cover.out

# Runs the tests from the ci
test_ci:
	CA_ENV=testing_ci go test -v cucinassistant/database cucinassistant/langs


# Formats the source code
fmt:
	go fmt ./...
	templ fmt .

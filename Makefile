db_container = ca_db
db_user = cucinassistant
db_pass = cucinassistant
db_name = cucinassistant
db_test = test


# Makes sure that the database is running
start_db:
	@if [ "$(shell docker ps -a -q -f name=$(db_container))" ]; then \
		if [ ! "$(shell docker ps -aq -f status='running' -f name=$(db_container))" ]; then \
			docker start $(db_container); \
		fi \
	else \
		docker run -d --name $(db_container) \
		-e POSTGRES_USER=$(db_user) \
		-e POSTGRES_PASSWORD=$(db_pass) \
		-e POSTGRES_DATABASE=$(db_name) \
		-p 5432:5432 postgres; \
	fi


# Drops the database entirely
drop_db:
	@docker stop $(db_container)
	@docker rm $(db_container)


# Creates a test table in the database
create_test_db: start_db
	@-docker exec -it $(db_container) dropdb -U $(db_user) "$(db_test)"
	@docker exec -it $(db_container) createdb -U $(db_user) "$(db_test)"


# Opens a shell with the database
dbsh: start_db
	@docker exec -it $(db_container) psql -U $(db_user) -d $(db_name)


# Runs the webserver
run: start_db
	@go run main.go config/debug.yml


# Runs the tests
test: create_test_db
	@go test -v cucinassistant/database -args config/test.yml || true


# Runs the tests and shows a coverage report
cover: create_test_db
	@go test -coverprofile=/tmp/cover.out -covermode atomic cucinassistant/database -args config/test.yml || true
	@go tool cover -html=/tmp/cover.out


# Runs go fmt
fmt:
	@go fmt ./...

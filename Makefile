DB_CONTAINER_NAME = ca_db
DB_USERNAME = ca
DB_PASSWORD = ca
DB_DEV_DB = ca_dev
DB_TEST_DB = ca_test


# Makes sure that the database is running
start_db:
	@if [ "$(shell docker ps -a -q -f name=$(DB_CONTAINER_NAME))" ]; then \
		if [ ! "$(shell docker ps -aq -f status='running' -f name=$(DB_CONTAINER_NAME))" ]; then \
			docker start $(DB_CONTAINER_NAME); \
		fi \
	else \
		docker run -d --name $(DB_CONTAINER_NAME) \
		-e POSTGRES_USER=$(DB_USERNAME) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_DEV_DB) \
		-p 5432:5432 postgres; \
	fi


# Drops the database entirely
drop_db:
	@docker stop $(DB_CONTAINER_NAME)
	@docker rm $(DB_CONTAINER_NAME)

# Creates the test database
create_test_db: start_db
	@-docker exec -it $(DB_CONTAINER_NAME) dropdb -U $(DB_USERNAME) "$(DB_TEST_DB)"
	@docker exec -it $(DB_CONTAINER_NAME) createdb -U $(DB_USERNAME) "$(DB_TEST_DB)"

# Opens a shell with the database
dbsh: start_db
	@docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USERNAME) -d $(DB_DEV_DB)

# Runs the webserver
run: start_db
	@cd src; CA_ENV=development go run .

# Runs the tests
test: create_test_db
	@-cd src; CA_ENV=testing go test -v cucinassistant/database

# Runs the tests and shows a coverage report
cover: create_test_db
	@cd src; CA_ENV=testing go test -coverprofile=/tmp/cover.out -covermode atomic cucinassistant/database
	@cd src; go tool cover -html=/tmp/cover.out

# Runs go fmt
fmt:
	@cd src; go fmt ./...

# Generates code
gen:
	@cd src; go generate ./...
	@cd src; templ generate

# Generates the translate.*.toml files
lang_gen:
	@cd src/langs; goi18n merge active.*.toml
    
# Merges the translate.*.toml into the active.*.toml
lang_save:
	@cd src/langs; goi18n merge active.*.toml translate.*.toml; rm translate.*.toml

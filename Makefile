start_db:
	@if [ "$(shell docker ps -a -q -f name='ca-db')" ]; then \
		if [ ! "$(shell docker ps -aq -f status='running' -f name='ca-db')" ]; then \
			docker start ca; \
		fi \
	else \
		docker run -d --name ca-db -e MARIADB_USER=ca-user -e MARIADB_PASSWORD=ca-pass -e MARIADB_DATABASE=cucinassistant -e MARIADB_ROOT_PASSWORD=rpass -p 3306:3306 mariadb:10.6; \
	fi

drop_db:
	@docker stop ca-db
	@docker rm ca-db

run: start_db
	@go run main.go config_debug.yml

fmt:
	@go fmt ./...
	@go mod tidy

test: start_db
	@docker exec -it ca-db mariadb -u root -prpass -e "DROP DATABASE IF EXISTS test; CREATE DATABASE test; GRANT ALL PRIVILEGES ON test.* TO 'ca-user'@'%';"
	@go test -v cucinassistant/database -args ../config_test.yml || true
	@docker exec -it ca-db mariadb -u root -prpass -e "DROP DATABASE test;"

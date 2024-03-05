.SILENT:run
.SILENT:build
.SILENT:dbsh
.SILENT:test
.SILENT:push


run:
	flask -A cucinassistant:app run -h 0.0.0.0 -p 8080 --debug

build:
	if [[ -z "$(version)" ]]; then \
		echo "You need to specify a version. The correct syntax is 'make build version=...'"; \
		exit 1; \
	fi

	docker build -t git.gianlucaparri.me/gianluparri03/cucinassistant:$(version) .

dbsh:
	docker exec -it cucinassistant-db mariadb -u cucinassistant -pcucinassistant cucinassistant

test:
	echo "Waiting for the environment to boot..."
	
	docker stop cucinassistant-test 2> /dev/null || true
	docker run --name cucinassistant-test -d --rm -e MARIADB_USER=test -e MARIADB_PASSWORD=test -e MARIADB_DATABASE=test -e MARIADB_RANDOM_ROOT_PASSWORD=1 -p 3307\:3306 mariadb\:10.6 > /dev/null
	echo "- Database created"
	
	while [[ 1 ]]; do docker exec cucinassistant-test mariadb -u test -ptest test -e '' 2> /dev/null && break; sleep .5; done
	echo "- Database ready"
	
	python3 -m unittest cucinassistant.database.tests || true
	docker stop cucinassistant-test > /dev/null

push:
	make build $(version)
	docker push git.gianlucaparri.me/gianluparri03/cucinassistant:$(version)

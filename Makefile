.SILENT:run
.SILENT:build
.SILENT:test
.SILENT:push

run:
	flask -A application:app run -h 0.0.0.0 -p 8080 --debug

build:
	docker build -t cucinassistant .

test:
	make build
	docker run --rm -p 8080:80 -v ./config.cfg:/cucinassistant/config.cfg cucinassistant
	docker rmi cucinassistant

push:
	if [[ -z "$(version)" ]]; then \
		echo "You need to specify a version. The correct syntax is 'make push version=...'"; \
		exit 1; \
	fi

	make build
	docker tag cucinassistant git.gianlucaparri.me/gianluparri03/cucinassistant:$(version)
	docker push git.gianlucaparri.me/gianluparri03/cucinassistant:$(version)

.SILENT:run
.SILENT:build
.SILENT:test
.SILENT:push


run:
	flask -A application:app run -h 0.0.0.0 -p 8080 --debug

build:
	if [[ -z "$(version)" ]]; then \
		echo "You need to specify a version. The correct syntax is 'make push version=...'"; \
		exit 1; \
	fi

	docker build -t git.gianlucaparri.me/gianluparri03/cucinassistant:$(version) .

test:
	make build version=test
	docker run --rm --name cucinassistant-test --network host -v ./config.cfg:/cucinassistant/config.cfg git.gianlucaparri.me/gianluparri03/cucinassistant:test

push:
	make build $(version)
	docker push git.gianlucaparri.me/gianluparri03/cucinassistant:$(version)

run:
	flask -A application:app run -h 0.0.0.0 -p 8080 --debug

build:
	docker build -t cucinassistant .

test:
	make build
	docker run --rm -p 8080:80 -v ./config.cfg:/cucinassistant/config.cfg cucinassistant
	docker rmi cucinassistant

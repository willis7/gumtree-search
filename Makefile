NAME = gumtree
PWD := $(MKPATH:%/MAKEFILE=%)

clean:
	cd "$(PWD)"
	rm -rf vendor

setup:
	glide install

build:
	go build -race -o $(NAME)

docker-build:
	docker build --rm -t willis7/gumtree-searcher .

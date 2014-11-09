## 
## Makefile
## 
## Created by Frederic DELBOS - fred@hyperboloide.com on Nov  9 2014.
## 

NAME = fe

all:
	go build

clean:
	rm -fr $(NAME)

cross:
	goxc -pv=$(VERSION) -bc="linux,amd64"
	rm -fr debian

debug: re
	./$(NAME) \
	--tmp ./tmp \
	--private 7031 \
	--public 7030 \
	--mongo bubble \
	--db test-fe \
	--coll files \
	--redis bubble:6379 \
	buckets.json

fmt:
	go fmt ./...

re: clean all

test:
	go test

.PHONY: all clean cross debug fmt re test

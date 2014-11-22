## 
## Makefile
## 
## Created by Frederic DELBOS - fred@hyperboloide.com on Nov  9 2014.
## 

NAME = fe
VERSION = 0.1.1

all:
	go build -v

clean:
	rm -fr $(NAME)

cross:
	goxc -pv=$(VERSION) -bc="linux,amd64"
	rm -fr debian

debug: re
	./$(NAME) --private 7031 --public 7030 config.json

fmt:
	go fmt ./...

re: clean all

test:
	go test ./...

.PHONY: all clean cross debug fmt re test

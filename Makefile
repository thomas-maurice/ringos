.PHONY: all
all: bin

.PHONY: bin
bin:
	if ! [ -d bin ]; then mkdir bin; fi;
	(cd ringo-cli; go build -o ../bin/ringo-cli)
	(cd randomizer; go build -o ../bin/randomizer)


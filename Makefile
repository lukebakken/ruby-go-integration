.PHONY: all install-deps clean

all: libsum.so
	bundle exec ./app.rb

libsum.so:
	go build -buildmode=c-shared -o libsum.so libsum.go

install-deps:
	bundler install --binstubs --path vendor

clean:
	rm -f libsum.h libsum.so

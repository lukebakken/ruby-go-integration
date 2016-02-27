.PHONY: all install-deps clean

all: riak-client.so
	bundle exec ./app.rb

riak-client.so:
	go build -buildmode=c-shared -o riak-client.so riak-client.go

install-deps:
	bundler install --binstubs --path vendor

clean:
	rm -f riak-client.h riak-client.so

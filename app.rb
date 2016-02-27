#!/usr/bin/env ruby

require 'ffi'
require 'benchmark'

module Riak
  extend FFI::Library
  ffi_lib './riak-client.so'
  attach_function :Start, [], :void
  attach_function :Stop, [], :void
  attach_function :Ping, [], :bool
end

Riak.Start()
puts "Ping result: #{Riak.Ping()}"
Riak.Stop()

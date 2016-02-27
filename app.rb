#!/usr/bin/env ruby

require 'ffi'
require 'benchmark'

module Riak
  extend FFI::Library
  ffi_lib './riak-client.so'
  attach_function :RiakClusterStart, [], :void
  attach_function :RiakClusterStop, [], :void
  attach_function :RiakClusterPing, [], :bool
end

Riak.RiakClusterStart()
puts "Ping result: #{Riak.RiakClusterPing()}"
Riak.RiakClusterStop()

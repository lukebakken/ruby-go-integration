#!/usr/bin/env ruby

require 'ffi'

module Riak
    # https://github.com/ffi/ffi/wiki/Pointers
    class FetchArgs < FFI::Struct
        layout :bucketType, :pointer,
            :bucket, :pointer,
            :key, :pointer

        def initialize(bt, b, key)
            self[:bucketType] = FFI::MemoryPointer.from_string(bt)
            self[:bucket] = FFI::MemoryPointer.from_string(b)
            self[:key] = FFI::MemoryPointer.from_string(key)
        end
    end

    extend FFI::Library
    ffi_lib './riak-client.so'

    callback :tcb_cb, [], :void
    attach_function :TestCallback, [ :tcb_cb ], :void

    attach_function :TestStruct, [ FetchArgs.by_value ], :void
    attach_function :Start, [], :void
    attach_function :Stop, [], :void
    attach_function :Ping, [], :bool

    TestCB = FFI::Function.new(:void, []) do |ignored|
        puts "IN CALLBACK"
    end
end

a = Riak::FetchArgs.new('rb type', 'rb bucket', 'rb key')
Riak.TestStruct(a)

Riak.TestCallback(Riak::TestCB)
# puts "Riak.TestCallback tcb_rv: #{tcb_rv}"

# Riak.Start()
# puts "Ping result: #{Riak.Ping()}"
# Riak.Stop()

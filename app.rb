#!/usr/bin/env ruby

require 'ffi'
require 'benchmark'

module Sum
  extend FFI::Library
  ffi_lib './libsum.so'
  attach_function :add, [:int, :int], :int
end

def ffi_bm
    x = rand(1000)
    y = rand(1000)
    for i in 1..50000
        Sum.add(x, y)
    end
end

def ruby_bm
    x = rand(1000)
    y = rand(1000)
    for i in 1..50000
        x + y
    end
end

Benchmark.bm(11) do |x|
  x.report('FFI first:')  { ffi_bm }
  x.report('FFI second:') { ffi_bm }
  x.report('FFI third:')  { ffi_bm }
end

Benchmark.bm(10) do |x|
  x.report('RB first:')  { ruby_bm }
  x.report('RB second:') { ruby_bm }
  x.report('RB third:')  { ruby_bm }
end

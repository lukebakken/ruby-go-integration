#!/usr/bin/env ruby

require 'benchmark'
require 'riak'

all_keys = []
File.open('/home/lbakken/Projects/basho/riak/dev/keys.dat', 'r') do |f|
    f.each_line do |key|
        all_keys.push(key.chomp)
    end
end

puts "Key count: #{all_keys.length}"

h = 'riak-test'
type = 'default'
bucket = 'tweets'
batchsz = 128

client = Riak::Client.new(:nodes => [
  {:host => h, :pb_port => 10017 },
  {:host => h, :pb_port => 10027 },
  {:host => h, :pb_port => 10037 },
  {:host => h, :pb_port => 10047 }
])
bucket = client.bucket('tweets')

def fetch_keys_multiget(batchsz, bucket, keys)
    keys.each_slice(batchsz) do |batch|
        bucket.get_many(batch)
    end
end

def fetch_keys(bucket, keys)
    keys.each {|k| bucket[k]}
end

Benchmark.bm(16) do |x|
    (1..3).each do |r|
        x.report("run #{r} (mult):") { fetch_keys_multiget(batchsz, bucket, all_keys) }
        x.report("run #{r}:") { fetch_keys(bucket, all_keys) }
    end
end

# Benchmark.bm do |x|
#     x.report { fetch_keys_multiget(batchsz, bucket, all_keys) }
# end

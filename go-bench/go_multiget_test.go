package riak_client

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	util "github.com/lukebakken/goutil"
	riak "github.com/basho/riak-go-client"
)

var c *riak.Cluster
var keys []string

func init() {
	keys = make([]string, 131072)
	file, err := os.Open("/home/lbakken/Projects/basho/riak/dev/keys.dat")
	if err != nil {
		util.ErrExit(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		keys[i] = scanner.Text()
		i++
	}

	if err := scanner.Err(); err != nil {
		util.ErrExit(err)
	}

	addr_fmt := "riak-test:%d"
	nodes := make([]*riak.Node, 4)
	for i = 0; i < 4; i++ {
		port := 10017 + (i * 10)
		var node *riak.Node
		nodeOpts := &riak.NodeOptions{
			RemoteAddress: fmt.Sprintf(addr_fmt, port),
		}
		node, err = riak.NewNode(nodeOpts)
		if err != nil {
			util.ErrExit(err)
		}
		if node == nil {
			util.ErrExit(errors.New("node was nil"))
		}
		nodes[i] = node
	}

	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}
	c, err = riak.NewCluster(opts)
	if err != nil {
		util.ErrExit(err)
	}

	if err = c.Start(); err != nil {
		util.ErrExit(err)
	}
	// util.LogDebug("[BenchmarkMultiget/init]", "Key count: %v", len(keys))
}

func BenchmarkMultiget(b *testing.B) {
	batchsz := 128
	count := len(keys)
	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i += batchsz {
			wg := &sync.WaitGroup{}
			s := i
			e := s + (batchsz - 1)
			for _, key := range keys[s:e] {
				cmd, fcerr := riak.NewFetchValueCommandBuilder().
					WithBucket("tweets").
					WithKey(key).
					Build()
				if fcerr != nil {
					util.ErrExit(fcerr)
				}
				async := &riak.Async{
					Command: cmd,
					Wait: wg,
				}
				if err := c.ExecuteAsync(async); err != nil {
					util.ErrExit(err)
				}
			}
			wg.Wait()
		}
	}
}

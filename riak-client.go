package main

import (
	"C"
	"errors"
	"fmt"
	"log"
	"os"

	riak "github.com/basho/riak-go-client"
)

var slog = log.New(os.Stdout, "", log.LstdFlags)
var elog = log.New(os.Stderr, "", log.LstdFlags)

func LogInfo(source, format string, v ...interface{}) {
	slog.Printf(fmt.Sprintf("[INFO] %s %s", source, format), v...)
}

func LogDebug(source, format string, v ...interface{}) {
	slog.Printf(fmt.Sprintf("[DEBUG] %s %s", source, format), v...)
}

func LogError(source, format string, v ...interface{}) {
	elog.Printf(fmt.Sprintf("[DEBUG] %s %s", source, format), v...)
}

func LogErr(source string, err error) {
	elog.Println("[ERROR]", source, err)
}

func ErrExit(err error) {
	LogErr("[APP]", err)
	os.Exit(1)
}

var cluster *riak.Cluster

//export RiakClusterStart
func RiakClusterStart() {
	var err error

	var node *riak.Node
	nodeOpts := &riak.NodeOptions{
		RemoteAddress: "127.0.0.1:10017",
	}
	node, err = riak.NewNode(nodeOpts)
	if err != nil {
		ErrExit(err)
	}
	if node == nil {
		ErrExit(errors.New("node was nil"))
	}

	nodes := []*riak.Node{node}
	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}
	cluster, err = riak.NewCluster(opts)
	if err != nil {
		ErrExit(err)
	}

	if err = cluster.Start(); err != nil {
		ErrExit(err)
	}
}

//export RiakClusterStop
func RiakClusterStop() {
	if cluster != nil {
		if err := cluster.Stop(); err != nil {
			ErrExit(err)
		}
	}
}

//export RiakClusterPing
func RiakClusterPing() bool {
	cmd := &riak.PingCommand{}
	if err := cluster.Execute(cmd); err != nil {
		ErrExit(err)
	}
	return cmd.Success()
}

//export RiakClusterGet
func RiakClusterGet(btype, bucket, key string) *C.char {
	var err error
	var cmd riak.Command

	builder := riak.NewFetchValueCommandBuilder()
	cmd, err = builder.WithBucketType(btype).
		WithBucket(bucket).
		WithKey(key).
		Build()
	if err != nil {
		ErrExit(err)
	}
	if err := cluster.Execute(cmd); err != nil {
		ErrExit(err)
	}

	fvc := cmd.(*riak.FetchValueCommand)
	rsp := fvc.Response
	if rsp.IsNotFound {
		return C.CString("")
	}
	if len(rsp.Values) == 0 {
		return C.CString("")
	}
	object := rsp.Values[0]
	return C.CString(string(object.Value))
}

func main() {
}

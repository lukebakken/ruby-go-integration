package main

/*
#include "riak-types.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"log"
	"os"

	riak "github.com/basho/riak-go-client"
)

type FetchArgs C.struct_fetchArgs

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

//export TestStruct
func TestStruct(a FetchArgs) {
	LogDebug("[TestStruct]", "bucketType: %v", C.GoString(a.bucketType))
	LogDebug("[TestStruct]", "bucket: %v", C.GoString(a.bucket))
	LogDebug("[TestStruct]", "key: %v", C.GoString(a.key))
}

//export Start
func Start() {
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

//export Stop
func Stop() {
	if cluster != nil {
		if err := cluster.Stop(); err != nil {
			ErrExit(err)
		}
	}
}

//export Ping
func Ping() bool {
	cmd := &riak.PingCommand{}
	if err := cluster.Execute(cmd); err != nil {
		ErrExit(err)
	}
	return cmd.Success()
}

//export Fetch
func Fetch(a FetchArgs) *C.char {
	var err error
	var cmd riak.Command

	bt := C.GoString(a.bucketType)
	b := C.GoString(a.bucket)
	k := C.GoString(a.key)

	builder := riak.NewFetchValueCommandBuilder()
	cmd, err = builder.WithBucketType(bt).
		WithBucket(b).
		WithKey(k).
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
	}
	object := rsp.Values[0]
	return C.CString(string(object.Value))
}

func main() {}

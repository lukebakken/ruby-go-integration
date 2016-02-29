package main

/*
#cgo LDFLAGS: libriak.a
#include "riak-types.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"

	util "github.com/lukebakken/goutil"
	riak "github.com/basho/riak-go-client"
)

type FetchArgs C.struct_fetchArgs

var cluster *riak.Cluster

//export TestStruct
func TestStruct(a FetchArgs) {
	util.LogDebug("[TestStruct]", "bucketType: %v", C.GoString(a.bucketType))
	util.LogDebug("[TestStruct]", "bucket: %v", C.GoString(a.bucket))
	util.LogDebug("[TestStruct]", "key: %v", C.GoString(a.key))
}

//export TestCallback
func TestCallback(pcb unsafe.Pointer) {
	util.LogDebug("[TestCallback]", "before calling pcb")
	C.call_tcb_cb((C.cb_fn)(pcb))
	util.LogDebug("[TestCallback]", "after calling pcb")
}

//export Start
func Start() {
	var err error

	addr_fmt := "riak-test:%d"
	nodes := make([]*riak.Node, 4)
	for i := 0; i < 4; i++ {
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
	cluster, err = riak.NewCluster(opts)
	if err != nil {
		util.ErrExit(err)
	}

	if err = cluster.Start(); err != nil {
		util.ErrExit(err)
	}
}

//export Stop
func Stop() {
	if cluster != nil {
		if err := cluster.Stop(); err != nil {
			util.ErrExit(err)
		}
	}
}

//export Ping
func Ping() bool {
	cmd := &riak.PingCommand{}
	if err := cluster.Execute(cmd); err != nil {
		util.ErrExit(err)
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
		util.ErrExit(err)
	}
	if err := cluster.Execute(cmd); err != nil {
		util.ErrExit(err)
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

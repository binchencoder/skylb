package main

import (
	etcd "github.com/coreos/etcd/client"
)

type nodeSlice []*etcd.Node

func (ns nodeSlice) Len() int {
	return len(ns)
}

func (ns nodeSlice) Less(i, j int) bool {
	return ns[i].Key < ns[j].Key
}

func (ns nodeSlice) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

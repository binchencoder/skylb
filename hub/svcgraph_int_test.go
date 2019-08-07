package hub

import (
	"testing"
)

func TestSkypbServiceGraph(t *testing.T) {
	etcdCli := GetTestEtcdClient()
	if err := BuildDependencies(etcdCli); nil != err {
		t.Fatal(err)
	}
	traverseCallRoots()
}

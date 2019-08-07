package hub

import (
	"fmt"
	"testing"

	"github.com/golang/glog"
)

func traverseCalls(caller string, level int) {
	sep := ""
	for i := 0; i < level; i++ {
		sep += "  "
	}
	fmt.Printf("%s %s\n", sep, caller)

	// Inefficient lookup.
	//for c, cps := range calls {
	//	if c.Caller != caller {
	//		continue
	//	}
	//	for _, cp := range cps {
	//		traverseCalls(cp.Callee, level + 1)
	//	}
	//}

	// Efficient lookup.
	cps := SimpleCallingMap[caller]
	for _, cp := range cps {
		if cp.Callee == caller {
			glog.Warningf("Loop detected: %s", caller)
			continue
		}
		traverseCalls(cp.Callee, level+1)
	}
}

func traverseCallRoots() {
	// Use fmt instead of glog so as to exclude the prefix timestamp etc.
	fmt.Println("------------")
	rs := FindRoots()
	fmt.Printf("roots: %#v\n", rs)
	for _, r := range rs {
		traverseCalls(SimpleCallingMap[r][0].Caller, 0)
	}
	fmt.Println("------------")
}

func printLatestCallers(t *testing.T) {
	// Here uses glog so as to print logs along with production code.
	glog.Infoln("LatestCallers:")
	for _, ci := range latestCallers {
		glog.Infof("  %v", ci)
	}
	glog.Infof("callingMap: %v", callingMap)
}

func TestSimpleCallingMap(t *testing.T) {
	fmt.Println(">>> SimpleCallingMap:")
	for caller, cps := range SimpleCallingMap {
		fmt.Printf("%v\n", caller)
		for _, cp := range cps {
			fmt.Printf("  %v\n", cp)
		}
	}
}

func TestCalledMap(t *testing.T) {
	fmt.Println(">>> CalledMap:")
	calledMap := GenCalledMap()
	for callee, callers := range calledMap {
		fmt.Printf("%v\n", callee)
		for _, cr := range callers {
			fmt.Printf("  %v\n", cr)
		}
	}
}

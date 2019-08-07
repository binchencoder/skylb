package hub

import (
	"testing"
)

type CallPairTest struct {
	p      string
	callee string
	caller string
}

var cctests = []CallPairTest{
	// Abnormal case
	{"/", INVALID_VALUE, INVALID_VALUE},
	{"/a", INVALID_VALUE, INVALID_VALUE},

	// Normal case
	{"/skylb/graph/default/idm-service/clients/windows-client", "idm-service", "windows-client"},
	{"/skylb/graph/default/recency-service/clients/mac-client", "recency-service", "mac-client"},

	// Simplified normal case
	{"/idm-service/clients/windows-client", "idm-service", "windows-client"},
	{"/recency-service/clients/mac-client", "recency-service", "mac-client"},
}

func TestCallPair(t *testing.T) {
	for _, tt := range cctests {
		cc := parseCallPair(tt.p)
		if tt.callee != cc.Callee || tt.caller != cc.Caller {
			t.Errorf("%s --> %v; want %s,%s", tt.p, cc, tt.callee, tt.caller)
		}
	}
}

func TestRecordCallPairs(t *testing.T) {
	callingMap = make(map[CallerInfo]CallerPairs)

	// c1 -> s1 (t1)
	recordCallPairs("/s1/clients/c1", 100)

	finishRecording()
	printLatestCallers(t)
	if 1 != len(callingMap) {
		t.Errorf("len err %v", callingMap)
	}
	if cp, _ := callingMap[CallerInfo{"c1", 100}]; 1 != len(cp) || "s1" != cp[0].Callee {
		t.Errorf("data err %v\n%v", cp, callingMap)
	}

	// c1 -> s1 (t2)
	// c1 -> s2 (t2)
	recordCallPairs("/s1/clients/c1", 200)
	recordCallPairs("/s2/clients/c1", 200)

	finishRecording()
	printLatestCallers(t)
	if 1 != len(callingMap) {
		t.Errorf("len err %v", callingMap)
	}
	if cp, _ := callingMap[CallerInfo{"c1", 200}]; 2 != len(cp) ||
		"s1" != cp[0].Callee || "s2" != cp[1].Callee {
		t.Errorf("data err %v\n%v", cp, callingMap)
	}

	// c1 -> s2 (t3)
	recordCallPairs("/s2/clients/c1", 300)

	finishRecording()
	printLatestCallers(t)
	if 1 != len(callingMap) {
		t.Errorf("len err %v", callingMap)
	}
	if cp, _ := callingMap[CallerInfo{"c1", 300}]; 1 != len(cp) || "s2" != cp[0].Callee {
		t.Errorf("data err %v\n%v", cp, callingMap)
	}

	// Now more complicated mixes.
	recordCallPairs("/sX/clients/s2", 400)
	recordCallPairs("/sY/clients/s2", 400)
	recordCallPairs("/sZ/clients/s3", 400)

	finishRecording()
	printLatestCallers(t)
	if 3 != len(callingMap) {
		t.Errorf("len err %v", callingMap)
	}
	if cp, _ := callingMap[CallerInfo{"c1", 300}]; 1 != len(cp) || "s2" != cp[0].Callee {
		t.Errorf("data err %v\n%v", cp, callingMap)
	}
	if cp, _ := callingMap[CallerInfo{"s2", 400}]; 2 != len(cp) ||
		"sX" != cp[0].Callee || "sY" != cp[1].Callee {
		t.Errorf("data err %v\n%v", cp, callingMap)
	}
	if cp, _ := callingMap[CallerInfo{"s3", 400}]; 1 != len(cp) || "sZ" != cp[0].Callee {
		t.Errorf("data err %v\n%v", cp, callingMap)
	}

	traverseCallRoots()
	/*
		Result:
		c1
		  s2
		    sX
		    sY
		s3
		  sZ
	*/

	// Make up a connected graph for the following test.
	recordCallPairs("/s3/clients/sX", 400)
	finishRecording()

	// Print the graph.
	root := latestCallers[0]
	traverseCalls(root.Caller, 0)
	/*
		Result:
		c1
		  s2
		    sX
		      s3
		        sZ
		    sY
	*/

	traverseCallRoots()
	/*
		Result:
		c1
		  s2
		    sX
		      s3
		        sZ
		    sY
	*/

	recordCallPairs("/sZ/clients/c2", 400)
	finishRecording()
	traverseCallRoots()
	/*
		Result:
		c1
		  s2
		    sX
		      s3
		        sZ
		    sY
		c2
		  sZ
	*/
}

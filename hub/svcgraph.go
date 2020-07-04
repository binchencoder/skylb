package hub

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"
	"golang.org/x/net/context"

	"github.com/binchencoder/skylb-api/prefix"
)

const (
	INVALID_VALUE = "INVALID"
)

type CallerInfo struct {
	Caller    string
	Timestamp int64
}

type CallPair struct {
	Callee string
	CallerInfo
}

type CallerPairs []*CallPair

type CallerNames []string

var (
	latestCallers = []*CallerInfo{}
	callingMap    = make(map[CallerInfo]CallerPairs)
	// SimpleCallingMap holds call map for fast traversing. caller --> callees.
	SimpleCallingMap = make(map[string]CallerPairs)

	sgLock sync.RWMutex
)

func (ci *CallerInfo) newerThan(other *CallerInfo) bool {
	if ci.Caller != other.Caller {
		// Not the same caller, so don't even bother to compare.
		return false
	}
	return ci.Timestamp > other.Timestamp
}

func (ci *CallerInfo) eq(other *CallerInfo) bool {
	return ci.Caller == other.Caller && ci.Timestamp == other.Timestamp
}

func (ci *CallerInfo) String() string {
	return fmt.Sprintf("Caller{%s,%d}", ci.Caller, ci.Timestamp)
}

func (cp *CallPair) String() string {
	return fmt.Sprintf("CallPair{%s->%s}", cp.Caller, cp.Callee)
}

func parseCallPair(path string) *CallPair {
	ss := strings.Split(path, "/")
	ln := len(ss)
	if ln < 4 {
		return &CallPair{
			Callee: INVALID_VALUE,
			CallerInfo: CallerInfo{
				Caller:    INVALID_VALUE,
				Timestamp: 0,
			},
		}
	}

	return &CallPair{
		Callee: ss[ln-3],
		CallerInfo: CallerInfo{
			Caller:    ss[ln-1],
			Timestamp: 0, // The Timestamp will be filled elsewhere.
		},
	}
}

func recordCallPairs(path string, timestamp int64) {
	sgLock.Lock()
	defer sgLock.Unlock()

	glog.V(4).Infof("path %s ts %d", path, timestamp)
	cp := parseCallPair(path)
	cp.Timestamp = timestamp
	// Record CallerInfo
	ci := cp.CallerInfo
	had := false // whether already had this CallerInfo.
	for i, c := range latestCallers {
		if c.eq(&ci) {
			had = true
			break
		}
		if ci.newerThan(c) {
			latestCallers[i] = &CallerInfo{} // Mark to remove later
			delete(callingMap, *c)
		}
	}
	if !had {
		latestCallers = append(latestCallers, &ci)
	}
	// Record CallPair
	cps, ok := callingMap[ci]
	if !ok {
		cps = CallerPairs{}
	}
	cps = append(cps, cp)
	callingMap[ci] = cps
}

func finishRecording() {
	sgLock.Lock()
	defer sgLock.Unlock()

	// Filter the records marked for removal.
	kept := []*CallerInfo{}
	for _, ci := range latestCallers {
		if ci.Timestamp > 0 {
			kept = append(kept, ci)
		}
	}
	latestCallers = kept
	// Generate simplified call pairs
	for ci, cps := range callingMap {
		SimpleCallingMap[ci.Caller] = cps
	}
}

func BuildDependencies(etcdCli etcd.KeysAPI) error {
	callingMap = make(map[CallerInfo]CallerPairs)
	SimpleCallingMap = make(map[string]CallerPairs)

	getOpts := etcd.GetOptions{
		Recursive: true,
	}
	resp, err := etcdCli.Get(context.Background(), prefix.GraphKey, &getOpts)
	if nil != err {
		return err
	}
	logLevel := glog.Level(5)
	for _, ns := range resp.Node.Nodes {
		glog.V(logLevel).Infof("namespace %#v", ns.Key)
		for _, svc := range ns.Nodes {
			glog.V(logLevel).Infof("> svc %#v", svc.Key)
			for _, cs := range svc.Nodes {
				glog.V(logLevel).Infof(">> 'clients' %#v", cs.Key)
				if !strings.HasSuffix(cs.Key, "/clients") {
					continue
				}
				for _, client := range cs.Nodes {
					if strings.HasSuffix(client.Key, AddrKey) {
						continue
					}
					glog.V(logLevel).Infof(">>> client %#v", client.Key)
					for _, tsNode := range client.Nodes {
						glog.V(logLevel).Infof(">>>> ts %#v", tsNode.Key)
						if !strings.HasSuffix(tsNode.Key, TimestampKey) {
							continue
						}
						ts, err := strconv.Atoi(tsNode.Value)
						if nil != err {
							glog.Warningf("%v err:%v", tsNode.Value, err)
							ts = 1
						}
						recordCallPairs(client.Key, int64(ts))
					}
				}
			}
		}
	}
	finishRecording()
	return nil
}

// FindRoots finds the root nodes in the service graph.
func FindRoots() []string {
	sgLock.RLock()
	defer sgLock.RUnlock()

	called := map[string]bool{}
	for _, cps := range SimpleCallingMap {
		for _, cp := range cps {
			called[cp.Callee] = true
		}
	}
	rootNodes := []string{}
	for c := range SimpleCallingMap {
		if !called[c] {
			rootNodes = append(rootNodes, c)
		}
	}
	return rootNodes
}

// GenCalledMap generates mapping of callee and callers.
func GenCalledMap() map[string]CallerNames {
	sgLock.RLock()
	defer sgLock.RUnlock()

	calledMap := map[string]CallerNames{}
	for caller, cps := range SimpleCallingMap {
		for _, cp := range cps {
			callee := cp.Callee
			callers, ok := calledMap[callee]
			if !ok {
				callers = CallerNames{}
			}
			callers = append(callers, caller)
			calledMap[callee] = callers
		}
	}

	return calledMap
}

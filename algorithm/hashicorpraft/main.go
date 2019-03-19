package main

import (
	"fmt"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type cacheManager struct {
	data map[string]string
	sync.RWMutex
}

var cm cacheManager
var host string = "8080"

func main() {
	setupcacheManager()
	ServerStart(host)
}

func ServerStart(port string) {
	http.HandleFunc("/set", sethandler)
	http.HandleFunc("/get", gethandler)
	log.Printf("server listen on :%s \n", host)
	http.ListenAndServe(fmt.Sprintf(":%s", host), nil)
}

func setupRaftConfig() {
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(host)
	raftConfig.Logger = log.New(os.Stderr, "raft: ", log.Ldate|log.Ltime)

	logStore, err := raftboltdb.NewBoltStore(filepath.Join("./",
		"raft-log.bolt"))
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join("./",
		"raft-stable.bolt"))

	snapshotStore, err := raft.NewFileSnapshotStore("./", 1, os.Stderr)

}
func setupcacheManager() {
	cm = cacheManager{
		data: make(map[string]string),
	}
}
func sethandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")
	cm.Lock()
	cm.data[key] = value
	log.Printf("set key %s of value: %s\n", key, value)
	cm.Unlock()
}

func gethandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	cm.Lock()
	if v, ok := cm.data[key]; ok {
		w.Write([]byte(v))
		log.Printf("get key %s of value: %s\n", key, v)
	}
	cm.Unlock()
}

func newRaftTransport() (*raft.NetworkTransport, error) {
	address, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(address.String(), address, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}
	return transport, nil
}

func (f *FSM) Apply(logEntry *raft.Log) interface{} {
	e := logEntryData{}
	if err := json.Unmarshal(logEntry.Data, &e); err != nil {
		panic("Failed unmarshaling Raft log entry.")
	}
	ret := f.ctx.st.cm.Set(e.Key, e.Value)
	return ret
}

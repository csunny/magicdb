// Copyright 2019 The magicdb Authors
//
// Licensed under the Apache Licence, Version 2.0(the "License");
// You may not use the file except in compliance with the Licence.
// You may obtain a copy of the Licence at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distrubuted under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package raft

import (
	"errors"
	"io/ioutil"
	"time"

	praft "github.com/hashicorp/raft"
	consensus "github.com/libp2p/go-libp2p-consensus"
	"github.com/libp2p/go-libp2p-core/peer"
	host "github.com/libp2p/go-libp2p-host"
	libp2praft "github.com/libp2p/go-libp2p-raft"
)

// TODO config this use config file
const (
	raftTmpFolder = "/tmp/magicdb"
)

var (
	// ErrExists
	ErrExists = errors.New("bootstrap already exists")
)

func NewRaftNode(peer host.Host, pids []peer.ID, op consensus.Op,
	raftQuiet bool) (*praft.Raft, *libp2praft.Consensus, *praft.NetworkTransport, error) {

	// consensus
	var cns *libp2praft.Consensus
	if op != nil {
		cns = libp2praft.NewOpLog(&raftState{}, op)
	} else {
		cns = libp2praft.NewConsensus(&raftState{3})
	}

	// Create Raft servers configuration
	servers := make([]praft.Server, len(pids))
	for i, pid := range pids {
		servers[i] = praft.Server{
			Suffrage: praft.Voter,
			ID:       praft.ServerID(pid.Pretty()),
			Address:  praft.ServerAddress(pid.Pretty()),
		}
	}

	serverConfig := praft.Configuration{Servers: servers}

	// transport
	transport, err := libp2praft.NewLibp2pTransport(peer, time.Minute)
	if err != nil {
		return nil, nil, nil, err
	}

	// config
	config := praft.DefaultConfig()
	if raftQuiet {
		config.LogOutput = ioutil.Discard
		config.Logger = nil
	}
	config.LocalID = praft.ServerID(peer.ID().Pretty())

	// Snapshot store, we can use disk, mem, and file.
	// There we use file for snapshot store
	snapshots, err := praft.NewFileSnapshotStore(raftTmpFolder, 3, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	// logstore & stable store.  There we use mem store.
	logStore := praft.NewInmemStore()

	// bootstrap  This should only be called at the begging of time for the
	// cluster with an identical configuration listing all Voter servers.
	bootstrapped, err := praft.HasExistingState(logStore, logStore, snapshots)
	if err != nil {
		return nil, nil, nil, err
	}
	if !bootstrapped {
		// Bootstrap cluster.
		praft.BootstrapCluster(config, logStore, logStore, snapshots, transport, serverConfig)
	}

	raftNode, err := praft.NewRaft(config, cns.FSM(), logStore, logStore, snapshots, transport)
	if err != nil {
		return nil, nil, nil, err
	}
	return raftNode, cns, transport, nil
}

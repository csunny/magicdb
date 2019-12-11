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

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/hashicorp/raft"
	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	libp2praft "github.com/libp2p/go-libp2p-raft"
)

func main() {

	// This example shows how to use go-libp2p-raft to create a cluster
	// which agrees on a State. In order to do it, it defines a state,
	// creates three Raft nodes and launches them. We call a function which
	// lets the cluster leader repeteadly update . At the
	// end of the execution we verify that all members have agreed on the same state.

	// Some error handling has been excluded for simplicity

	// Declare an object which represents the State.
	// Note that State objects should have public/exported fields.
	// as they are [de]serialized
	type raftState struct {
		Value int
	}

	// error handling ommitted
	newPeer := func(listenPort int) host.Host {
		h, _ := libp2p.New(
			context.Background(),
			libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		)

		return h
	}

	// Create peers and make sure they known about each others.
	peer1 := newPeer(9997)
	peer2 := newPeer(9998)
	peer3 := newPeer(9999)

	defer peer1.Close()
	defer peer2.Close()
	defer peer3.Close()

	peer1.Peerstore().AddAddrs(peer2.ID(), peer2.Addrs(), peerstore.PermanentAddrTTL)
	peer1.Peerstore().AddAddrs(peer3.ID(), peer3.Addrs(), peerstore.PermanentAddrTTL)

	peer2.Peerstore().AddAddrs(peer1.ID(), peer1.Addrs(), peerstore.PermanentAddrTTL)
	peer2.Peerstore().AddAddrs(peer3.ID(), peer3.Addrs(), peerstore.PermanentAddrTTL)

	peer3.Peerstore().AddAddrs(peer1.ID(), peer1.Addrs(), peerstore.PermanentAddrTTL)
	peer3.Peerstore().AddAddrs(peer2.ID(), peer2.Addrs(), peerstore.PermanentAddrTTL)

	// Create the consensus instances and initialize them with a state.
	// Note that state is just used for local initialization, and that,
	// only states submitted via CommitState() alters the state of the
	// cluster.
	consensus1 := libp2praft.NewConsensus(&raftState{3})
	consensus2 := libp2praft.NewConsensus(&raftState{3})
	consensus3 := libp2praft.NewConsensus(&raftState{3})

	// Create Libp2p transport raft
	transport1, err := libp2praft.NewLibp2pTransport(peer1, time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}

	transport2, err := libp2praft.NewLibp2pTransport(peer2, time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}

	transport3, err := libp2praft.NewLibp2pTransport(peer3, time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer transport1.Close()
	defer transport2.Close()
	defer transport3.Close()

	// Create Raft servers configuration for bootstrapping the cluster
	// Note that both IDs and Address are set to the Peer ID.
	servers := make([]raft.Server, 0)
	for _, h := range []host.Host{peer1, peer2, peer3} {
		servers = append(servers, raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(h.ID().Pretty()),
			Address:  raft.ServerAddress(h.ID().Pretty()),
		})
	}

	serversCfg := raft.Configuration{Servers: servers}

	// Create Raft Configs. The Local ID is PeerOID
	config1 := raft.DefaultConfig()
	config1.LogOutput = ioutil.Discard
	config1.Logger = nil
	config1.LocalID = raft.ServerID(peer1.ID().Pretty())

	config2 := raft.DefaultConfig()
	config2.LogOutput = ioutil.Discard
	config2.Logger = nil
	config2.LocalID = raft.ServerID(peer2.ID().Pretty())

	config3 := raft.DefaultConfig()
	config3.LogOutput = ioutil.Discard
	config3.Logger = nil
	config3.LocalID = raft.ServerID(peer2.ID().Pretty())

	// Create snapshotStores. Use FileSnapshotStore in production.
	snapshot1 := raft.NewInmemSnapshotStore()
	snapshot2 := raft.NewInmemSnapshotStore()
	snapshot3 := raft.NewInmemSnapshotStore()

	// Create the InmemStores for use as long store and stable store
	logStore1 := raft.NewInmemStore()
	logStore2 := raft.NewInmemStore()
	logStore3 := raft.NewInmemStore()

	// Bootstrap the stores with the serverConfigs
	raft.BootstrapCluster(config1, logStore1, logStore1, snapshot1, transport1, serversCfg.Clone())
	raft.BootstrapCluster(config2, logStore2, logStore2, snapshot2, transport2, serversCfg.Clone())
	raft.BootstrapCluster(config3, logStore3, logStore3, snapshot3, transport3, serversCfg.Clone())

	// Create Raft objects. Our consensus provides an implementation of
	// Raft.FSM

	raft1, err := raft.NewRaft(config1, consensus1.FSM(), logStore1, logStore1, snapshot1, transport1)
	if err != nil {
		log.Fatal(err)
	}

	raft2, err := raft.NewRaft(config2, consensus2.FSM(), logStore2, logStore2, snapshot2, transport2)
	if err != nil {
		log.Fatal(err)
	}

	raft3, err := raft.NewRaft(config3, consensus3.FSM(), logStore3, logStore3, snapshot3, transport3)
	if err != nil {
		log.Fatal(err)
	}

	// Create the actors using the Raft nodes
	actor1 := libp2praft.NewActor(raft1)
	actor2 := libp2praft.NewActor(raft2)
	actor3 := libp2praft.NewActor(raft3)

	// Set the actors so that we can CommitState() and GetCureentState()
	consensus1.SetActor(actor1)
	consensus2.SetActor(actor2)
	consensus3.SetActor(actor3)

	// This function updates the cluster state commiting 1000 updates.
	updateState := func(c *libp2praft.Consensus) {
		nUpdates := 0
		for {
			if nUpdates >= 1000 {
				break
			}

			newState := &raftState{nUpdates * 2}

			// CommitState() blocks until the state has been
			// agreed upon by everyone
			agreedState, err := c.CommitState(newState)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if agreedState == nil {
				fmt.Printf("agreedState is nil: committed on a non-leader?")
				continue
			}

			agreedRaftState := agreedState.(*raftState)
			nUpdates++

			if nUpdates%200 == 0 {
				fmt.Printf("Performed %d updates. Current State value: %d\n",
					nUpdates, agreedRaftState.Value)
			}
		}
	}

	// Provide some time for leader election
	time.Sleep(5 * time.Second)

	// Run the 1000 updates on the leader
	// Barrier() will wait until updates have been applied
	if actor1.IsLeader() {
		updateState(consensus1)
	} else if actor2.IsLeader() {
		updateState(consensus2)
	} else if actor3.IsLeader() {
		updateState(consensus3)
	}

	// Wait for updates to arrive
	time.Sleep(5 * time.Second)

	// Shutdown raft and wait for it to complete
	// (ignoring errors)
	raft1.Shutdown().Error()
	raft2.Shutdown().Error()
	raft3.Shutdown().Error()

	// Final states
	finalState1, err := consensus1.GetCurrentState()
	if err != nil {
		fmt.Println(err)
		return
	}

	finalState2, err := consensus2.GetCurrentState()
	if err != nil {
		fmt.Println(err)
		return
	}

	finalState3, err := consensus3.GetCurrentState()
	if err != nil {
		fmt.Println(err)
		return
	}

	finalRaftState1 := finalState1.(*raftState)
	finalRaftState2 := finalState2.(*raftState)
	finalRaftState3 := finalState3.(*raftState)

	fmt.Printf("Raft1 final state: %d\n", finalRaftState1.Value)
	fmt.Printf("Raft2 final state: %d\n", finalRaftState2.Value)
	fmt.Printf("Raft3 final state: %d\n", finalRaftState3.Value)
}

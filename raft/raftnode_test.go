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
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	host "github.com/libp2p/go-libp2p-host"
	libp2praft "github.com/libp2p/go-libp2p-raft"
)

func TestNewRaftNode(t *testing.T) {
	peer1, err := NewNode(9999)
	if err != nil {
		t.Fatal("Create node1 error", err)
	}
	peer2, err := NewNode(9998)
	if err != nil {
		t.Fatal("Create node2 error", err)
	}

	peer3, err := NewNode(9997)
	if err != nil {
		t.Fatal("Create node3 error", err)
	}

	defer peer1.Close()
	defer peer2.Close()
	defer peer3.Close()

	peer1.Peerstore().AddAddrs(peer2.ID(), peer2.Addrs(), peerstore.PermanentAddrTTL)
	peer1.Peerstore().AddAddrs(peer3.ID(), peer3.Addrs(), peerstore.PermanentAddrTTL)
	peer2.Peerstore().AddAddrs(peer1.ID(), peer1.Addrs(), peerstore.PermanentAddrTTL)
	peer2.Peerstore().AddAddrs(peer3.ID(), peer3.Addrs(), peerstore.PermanentAddrTTL)
	peer3.Peerstore().AddAddrs(peer1.ID(), peer1.Addrs(), peerstore.PermanentAddrTTL)
	peer3.Peerstore().AddAddrs(peer2.ID(), peer2.Addrs(), peerstore.PermanentAddrTTL)

	pids := make([]peer.ID, 0)
	for _, p := range []host.Host{peer1, peer2, peer3} {
		pids = append(pids, p.ID())
	}

	// Create raft node
	raft1, consensus1, transport1, err := NewRaftNode(peer1, pids, nil, true)
	if err != nil {
		t.Fatal("Create raftnode node1 error", err)
	}
	raft2, consensus2, transport2, err := NewRaftNode(peer2, pids, nil, true)
	if err != nil {
		t.Fatal("Create raftnode node2 error", err)
	}
	raft3, consensus3, transport3, err := NewRaftNode(peer3, pids, nil, true)
	if err != nil {
		t.Fatal("Create raftnode node3 error", err)
	}
	defer transport1.Close()
	defer transport2.Close()
	defer transport3.Close()

	// Create the actors using the raft nodes
	actor1 := libp2praft.NewActor(raft1)
	actor2 := libp2praft.NewActor(raft2)
	actor3 := libp2praft.NewActor(raft3)

	// Set the actors so that we can CommitState() and GetCurrentState()
	consensus1.SetActor(actor1)
	consensus2.SetActor(actor2)
	consensus3.SetActor(actor3)

	// This function updates the cluster state commiting 1000 updates
	updateState := func(c *libp2praft.Consensus) {
		nUpdates := 0
		for {
			if nUpdates >= 2000 {
				break
			}
			newState := &raftState{nUpdates * 2}

			// CommitState() blocks until the state has been
			// agreed upon by everyone
			agreedState, err := c.CommitState(newState)
			if err != nil {
				t.Fatal(err)
				continue
			}
			if agreedState == nil {
				fmt.Println("agreedState is nil: commited on a non-leader?")
				continue
			}
			agreedRaftState := agreedState.(*raftState)
			nUpdates++

			if nUpdates%100 == 0 {
				fmt.Printf("Performed %d updates. Current state value: %d\n", nUpdates, agreedRaftState.Value)
			}
		}
	}

	// Provide some time for leader election
	time.Sleep(5 * time.Second)

	// Run the 1000 updates on the leader
	// Barrier() will wait until updates have been applied
	if actor1.IsLeader() {
		fmt.Println("I am leader: ", actor1)
		updateState(consensus1)
	} else if actor2.IsLeader() {
		fmt.Println("I am leader: ", actor2)
		updateState(consensus2)
	} else if actor3.IsLeader() {
		fmt.Println("I am leader: ", actor3)
		updateState(consensus3)
	}

	// Wait for updates to arrive.
	time.Sleep(5 * time.Second)

	// Shutdown raft and wait for it to complete
	// (ignoring errors)
	raft1.Shutdown().Error()
	raft2.Shutdown().Error()
	raft3.Shutdown().Error()

	// Final states
	finnalState1, err := consensus1.GetCurrentState()
	if err != nil {
		t.Fatal("Get Current state error ", err)
	}
	finnalState2, err := consensus2.GetCurrentState()
	if err != nil {
		t.Fatal("Get Current state error ", err)
	}
	finnalState3, err := consensus3.GetCurrentState()
	if err != nil {
		t.Fatal("Get current state error ", err)
	}

	finalRaftState1 := finnalState1.(*raftState)
	finalRaftState2 := finnalState2.(*raftState)
	finalRaftState3 := finnalState3.(*raftState)

	fmt.Printf("Raft1 final state: %d\n", finalRaftState1.Value)
	fmt.Printf("Raft2 final state: %d\n", finalRaftState2.Value)
	fmt.Printf("Raft3 final state: %d\n", finalRaftState3.Value)

}

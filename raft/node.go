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
	"context"
	"fmt"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	ma "github.com/multiformats/go-multiaddr"
)

var protocolVersion = "ipfs"

type raftState struct {
	Value int
}

// NewNode create a new node
func NewNode(port int) (host.Host, error) {

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)),
	}

	h, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	addr := h.Addrs()[0]
	maAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/%s/%s", protocolVersion, h.ID().Pretty()))
	fullAddr := addr.Encapsulate(maAddr)
	fmt.Println("I am running at addr: ", maAddr)
	fmt.Println("The full addr is: ", fullAddr)
	return h, err
}

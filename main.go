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
	"bufio"
	"flag"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p-core/network"
	raft "github.com/magicdb/raft"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	fmt.Println("handsome, raft run...")

	port := flag.Int("p", 0, "wait for incoming connections")
	dest := flag.String("d", "", "dest node to dail")

	flag.Parse()

	n, err := raft.NewNode(*port)
	if err != nil {
		log.Fatal(err)
	}

	// Set streamhandler
	n.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("Got a stream")
		if err := doEcho(s); err != nil {
			log.Println(err)
		} else {
			s.Close()
		}
	})

	// if not target. It's the first node
	if *dest == "" {
		log.Println("listening for connections")
		select {}
	}

	nodeAddr, err := ma.NewMultiaddr(*dest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(nodeAddr)
}

// do Echo reads a line of data a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}
	log.Printf("read:%s\n", str)
	_, err = s.Write([]byte(str))
	return err
}

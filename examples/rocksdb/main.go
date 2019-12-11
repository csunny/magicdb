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

/*
#cgo CFLAGS: -I/usr/local/lib/rocksdb/include
#cgo LDFLAGS: -L/usr/local/lib/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd
*/
import "C"
import (
	"fmt"

	"github.com/tecbot/gorocksdb"
)

var tmpPath = "/tmp/magicdb"

func main() {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))

	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)

	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, tmpPath)

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	wo := gorocksdb.NewDefaultWriteOptions()
	ro := gorocksdb.NewDefaultReadOptions()

	err = db.Put(wo, []byte("magic"), []byte("1"))
	value, err := db.Get(ro, []byte("magic"))

	v := fmt.Sprintf("%s", value.Data())

	fmt.Println(v == "1")

	var key interface{}
	key = "3"
	byteKey := []byte(fmt.Sprintf("%v", key.(interface{})))
	fmt.Println(byteKey)

	fmt.Println(string(byteKey) == "3")
}

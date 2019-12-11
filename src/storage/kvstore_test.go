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
package storage

import (
	"testing"

	"github.com/tecbot/gorocksdb"
)

var tmpPath = "/tmp/magicdb"

func buildOpts() *gorocksdb.Options {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))

	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	return opts
}

func TestNewKvStore(t *testing.T) {

	opts := buildOpts()
	store, err := NewKvStore(opts, tmpPath)
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
}

func TestOperate(t *testing.T) {
	opts := buildOpts()

	store, err := NewKvStore(opts, tmpPath)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	// Put test
	err = store.Put("foo", "bar")
	if err != nil {
		t.Fatal("Put k-v foo-bar error: ", err)
	}

	// Get test
	v, err := store.Get("foo")
	if err != nil {
		t.Fatal("Get k-v foo from store error: ", err)
	}

	if string(v) != "bar" {
		t.Fatal("Get value not equal set value, the get value is " + string(v) + " but excepted is 'bar' ")
	}
}

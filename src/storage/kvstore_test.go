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

	// Delete test
	err = store.Delete("foo")
	if err != nil {
		t.Fatal("Delete from store error ", err)
	}

	v2, err := store.Get("foo")
	if err != nil {
		t.Fatal("Get value error ", err)
	}

	if v2 != nil {
		t.Fatal("Delete key error from store, excepted nil but value get ", v2)
	}
}

func TestBatchOperate(t *testing.T) {
	opts := buildOpts()

	store, err := NewKvStore(opts, tmpPath)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	var _item []kv
	// Batch test
	for i := 0; i < 10; i++ {
		_item = append(_item, kv{i, i})
	}
	err = store.BatchPut(_item)
	if err != nil {
		t.Fatal("Batch put error ", err)
	}

	// Test if put
	for i := 0; i < 10; i++ {
		_v, err := store.Get(i)
		if err != nil {
			t.Fatal("Batch operate, Get error ", err)
		}
		if _v == nil {
			t.Fatal("Batch delete error, excepted nil, but got ", string(_v))
		}
	}

	var keys []item
	// Batch Delete
	for i := 0; i < 10; i++ {
		keys = append(keys, i)
	}
	err = store.BatchDelete(keys)
	if err != nil {
		t.Fatal("Batch Delete error ", err)
	}

	// Test if delete
	for i := 0; i < 10; i++ {
		_v, err := store.Get(i)
		if err != nil {
			t.Fatal("Batch operate, Get error ", err)
		}
		if _v != nil {
			t.Fatal("Batch delete error, excepted nil, but got ", _v)
		}
	}
}

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

// storage is the storage layer for magicdb
package storage

import (
	"sync"

	"github.com/tecbot/gorocksdb"
)

type kvstore struct {
	db *gorocksdb.DB
	mu sync.RWMutex
}

type item interface{}

type kvstorer interface {

	// Put a key-value to store
	Put(key, val item) error

	// Get get a key-value from store
	Get(key item) (item, error)

	// Delete a key from store
	Delete(key item) error

	//Batch operate
	//Batch Put
	BatchPut(kvpair ...kv) error

	// BatchDelete
	BatchDelete(k ...item) error

	// Close the kvstore
	Close()

	// Clear the kvstore and reuse
	Clear()
}

type kv struct {
	Key   string
	Value string
}

func newKvStore(opts *gorocksdb.Options, name string) (*kvstore, error) {
	db, err := gorocksdb.OpenDb(opts, name)
	if err != nil {
		return nil, err
	}
	store := &kvstore{db: db}
	return store, nil
}

// Put a key-value to store
func (s *kvstore) Put(k, v item) error {
	wo := gorocksdb.NewDefaultWriteOptions()
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.db.Put(wo, k.([]byte), v.([]byte))
	return err
}

// Get a key from store
func (s *kvstore) Get(k item) (item, error) {
	ro := gorocksdb.NewDefaultReadOptions()
	s.mu.Lock()
	defer s.mu.Unlock()
	value, err := s.db.Get(ro, k.([]byte))
	if err != nil {
		return nil, err
	}
	return value, nil
}

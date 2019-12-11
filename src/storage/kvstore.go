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
	Get(key item) ([]byte, error)

	// Delete a key from store
	Delete(key item) error

	//Batch operate
	//Batch Put
	BatchPut(kvpair ...kv) error

	// BatchDelete
	BatchDelete(k ...item) error

	// Close the kvstore
	Close()
}

type kv struct {
	Key   item
	Value item
}

// NewKvStore create a kvstore object
func NewKvStore(opts *gorocksdb.Options, name string) (*kvstore, error) {
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

	byteK := toBytes(k)
	byteV := toBytes(v)
	err := s.db.Put(wo, byteK, byteV)
	return err
}

// Get a key from store
func (s *kvstore) Get(k item) ([]byte, error) {
	ro := gorocksdb.NewDefaultReadOptions()
	s.mu.Lock()
	defer s.mu.Unlock()

	byteK := toBytes(k)
	value, err := s.db.Get(ro, byteK)
	if err != nil {
		return nil, err
	}
	return value.Data(), nil
}

// Delete the key-value pair from store
func (s *kvstore) Delete(k item) error {
	wo := gorocksdb.NewDefaultWriteOptions()
	s.mu.Lock()
	defer s.mu.Unlock()

	byteK := toBytes(k)
	err := s.db.Delete(wo, byteK)
	return err
}

// BatchPut batch put a batch of k-v pairs to store
func (s *kvstore) BatchPut(kvpair ...kv) error {
	wo := gorocksdb.NewDefaultWriteOptions()
	wb := gorocksdb.NewWriteBatch()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, _kv := range kvpair {
		byteK := toBytes(_kv.Key)
		byteV := toBytes(_kv.Value)
		wb.Put(byteK, byteV)
	}
	err := s.db.Write(wo, wb)
	return err
}

// BatchDelete delete a batch of kv pairs from store
func (s *kvstore) BatchDelete(k ...item) error {
	wo := gorocksdb.NewDefaultWriteOptions()
	wb := gorocksdb.NewWriteBatch()

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, _k := range k {
		byteK := toBytes(_k)
		wb.Delete(byteK)
	}
	err := s.db.Write(wo, wb)
	return err
}

// Close close the store db
func (s *kvstore) Close() {
	s.db.Close()
}

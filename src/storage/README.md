# magicdb storage use rocksdb to implement

At the first, use rocksdb impl a single node k-v storage server, and then use libp2p raft to rewrite.

go go go !!!


# rocksdb install

## Mac Install  
> brew install rocksdb

### Install gorocksdb

CGO_CFLAGS="-I/usr/local/lib/rocksdb/include" \
CGO_LDFLAGS="-L/usr/local/lib/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" \
  go get github.com/tecbot/gorocksdb
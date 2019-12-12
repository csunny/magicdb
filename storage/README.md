# magicdb storage use rocksdb to implement

At the first, use rocksdb impl a single node k-v storage server, and then use libp2p raft to rewrite.

go go go !!!


# rocksdb install

## Mac Install  
> brew install rocksdb

### Install gorocksdb

```
CGO_CFLAGS="-I/usr/local/lib/rocksdb/include" \
CGO_LDFLAGS="-L/usr/local/lib/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" \
  go get github.com/tecbot/gorocksdb

```

# Code Analysis

```
 git log --since=2018-01-01 --until=2020-12-31 --author="csunny" --pretty=tformat: --numstat | gawk '{ add += $1 ; subs += $2 ; loc += $1 - $2 } END { printf "added lines: %s removed lines : %s total lines: %s\n",add,subs,loc }' -
```

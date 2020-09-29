package jetton

import (
	"crypto/rand"
	"math/big"
	"sync"
)

// vNodeID is the node ID of the current replica instance.
type vNodeID int64

// uniqueID generates a unique ID.
func uniqueID() int64 {
	randI, err := rand.Int(rand.Reader, big.NewInt(255))
	if err != nil {
		panic(err)
	}
	saltI := randI.Int64()
	randA, err := rand.Int(rand.Reader, big.NewInt(255))
	if err != nil {
		panic(err)
	}

	saltA := randA.Int64()
	randB, err := rand.Int(rand.Reader, big.NewInt(255))
	if err != nil {
		panic(err)
	}
	saltB := randB.Int64()
	return saltI<<uint(16) | saltA<<uint(8) | saltB
}

func getVNodeID() vNodeID {
	return vNodeID(uniqueID())
}

type nodeCache struct {
	nodeID vNodeID
	// number of shards should be power of 2
	shards         int
	ShardedCounter []*shardedCounter
	shardMask      uint64
}

// shardedCounter represents each counter that has been sharded inside the given node.
type shardedCounter struct {
	nodeID    vNodeID
	lock      sync.RWMutex
	dotVector map[uint64]uint64
	storage   map[uint64]int64
}

func newShardCounter(nID vNodeID) *shardedCounter {
	return &shardedCounter{
		nodeID:    nID,
		dotVector: make(map[uint64]uint64),
		storage:   make(map[uint64]int64),
	}
}

// incrementBy allows passing an arbitrary delta to increment
// the current shard storage keyed value.
func (s *shardedCounter) incrementBy(key uint64, incr int64, tag uint64) {
	s.lock.Lock()
	s.storage[key] = s.storage[key] + incr
	s.dotVector[key] = tag
	s.lock.Unlock()
}

// decrementBy allows passing an arbitrary delta to decrement
// the current shard storage keyed value.
func (s *shardedCounter) decrementBy(key uint64, decr int64, tag uint64) {
	s.lock.Lock()
	s.storage[key] = s.storage[key] - decr
	s.dotVector[key] = tag
	s.lock.Unlock()
}

// set allows passing an arbitrary value to set
// the current shard storage keyed value.
func (s *shardedCounter) set(key uint64, val int64, tag uint64) {
	s.lock.Lock()
	s.storage[key] = val
	s.dotVector[key] = tag
	s.lock.Unlock()
}

// get returns value present or not. Second value denotes if value was
// present or not.
func (s *shardedCounter) get(key uint64) (int64, bool) {
	s.lock.RLock()
	v, ok := s.storage[key]
	s.lock.RUnlock()
	return v, ok
}

// delete the associated key and tag
func (s *shardedCounter) delete(key uint64) uint64 {
	s.lock.Lock()
	delete(s.storage, key)
	tag := s.dotVector[key]
	delete(s.dotVector, key)
	s.lock.Unlock()
	return tag
}

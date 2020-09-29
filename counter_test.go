// +build unit_tests all_tests

package jetton

import "testing"

func TestShardedCounterBasicOp(t *testing.T) {
	sc := newShardCounter(getVNodeID())
	key := "test1"
	sc.set(fnva1Hash(key), 12, 1)
	val, ok := sc.get(fnva1Hash(key))
	if val != 12 && ok != true {
		t.Errorf("expected return value 12 got %d", val)
	}
	sc.incrementBy(fnva1Hash(key), 2, 2)
	sc.decrementBy(fnva1Hash(key), 3, 3)
	val, ok = sc.get(fnva1Hash(key))
	if val != 11 && ok != true {
		t.Errorf("expected return value 11 got %d", val)
	}
}

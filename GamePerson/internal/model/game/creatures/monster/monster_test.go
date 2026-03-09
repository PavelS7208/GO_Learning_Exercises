package monster

import (
	"testing"
	"unsafe"
)

func TestGameMonsterSize(t *testing.T) {
	const expectedSize = 64
	actual := unsafe.Sizeof(monster{})
	if actual != expectedSize {
		t.Fatalf("GameMonster size MUST be %d bytes, got %d. "+
			"Check field order and alignment!", expectedSize, actual)
	}
}

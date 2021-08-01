package lungo

import (
	"math"
	"testing"
)

func TestByteToBinaryIEC(t *testing.T) {
	assertEqual(t, "0 B", byteToBinaryIEC(0))
	assertEqual(t, "27 B", byteToBinaryIEC(27))
	assertEqual(t, "999 B", byteToBinaryIEC(999))
	assertEqual(t, "1000 B", byteToBinaryIEC(1000))
	assertEqual(t, "1001 B", byteToBinaryIEC(1001))
	assertEqual(t, "1023 B", byteToBinaryIEC(1023))
	assertEqual(t, "1.0 KiB", byteToBinaryIEC(1024))
	assertEqual(t, "1.0 KiB", byteToBinaryIEC(1025))
	assertEqual(t, "8.0 EiB", byteToBinaryIEC(math.MaxInt64))
}

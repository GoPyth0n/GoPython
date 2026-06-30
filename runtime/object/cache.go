package object

import "math/big"

var (
	// Min and Max defined for your cache range
	cacheMin = -5
	cacheMax = 256
	// The cache array
	IntCache [262]*PyLongObject
)

func init() {
	for i := cacheMin; i <= cacheMax; i++ {
		IntCache[i-cacheMin] = &PyLongObject{Value: big.NewInt(int64(i))}
	}
}

func GetCachedInt(v int) *PyLongObject {
	if v >= cacheMin && v <= cacheMax {
		return IntCache[v-cacheMin]
	}
	return nil
}

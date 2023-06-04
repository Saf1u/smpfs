package disk

func CompareHashes(h1 []byte, h2 []byte) bool {
	for i := range h1 {
		if h1[i] != h2[i] {
			return false
		}
	}
	return true
}

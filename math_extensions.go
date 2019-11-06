package h3

func _ipow(base int, exp int) int {
	result := 1
	for exp != 0 {
		if exp&1 != 0 {
			result *= base
		}
		exp >>= uint64(1)
		base *= base
	}
	return result
}

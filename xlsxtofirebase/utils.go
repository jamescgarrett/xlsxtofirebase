package xlxstofirebase

// IntInArray - Check if int is in array of int
func IntInArray(integer int, integers []int) bool {
	for _, int := range integers {
		if int == integer {
			return true
		}
	}

	return false
}

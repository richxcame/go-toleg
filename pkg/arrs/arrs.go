package arrs

// Contains checks the value including in []values
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// HasKey checks existence of value with the given key
func HasKey(obj map[string]string, key string) (string, bool) {
	val, ok := obj[key]
	if ok {
		return val, true
	}
	return "", false
}

// HasMapWithKey is like HasKey but with array of maps
// If has map with the given key returns only value
func HasMapWithKey(arr []map[string]string, key string) (string, bool) {
	for _, v := range arr {
		val, has := HasKey(v, key)
		if has {
			return val, true
		}
	}
	return "", false
}

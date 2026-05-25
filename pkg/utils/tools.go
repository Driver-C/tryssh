package utils

// ToInterfaceSlice converts a typed slice to []interface{} using generics.
func ToInterfaceSlice[T any](s []T) []interface{} {
	if s == nil {
		return []interface{}{}
	}
	ret := make([]interface{}, len(s))
	for i, v := range s {
		ret[i] = v
	}
	return ret
}

// RemoveDuplicate removes duplicate strings from the slice, preserving order.
func RemoveDuplicate(s []string) []string {
	result := make([]string, 0, len(s))
	seen := make(map[string]struct{}, len(s))

	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// MaskSecret masks a secret string, returning a fixed-length indicator.
func MaskSecret(s string) string {
	if len(s) == 0 {
		return "<empty>"
	}
	return "****"
}

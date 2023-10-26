package maps

// Contains 判断 m1 是否包含 m2
func Contains[M ~map[K]V, K comparable, V comparable](m1, m2 M) bool {
	if len(m2) > len(m1) {
		return false
	}
	for k, v := range m2 {
		if vv, ok := m1[k]; !ok || vv != v {
			return false
		}
	}
	return true
}

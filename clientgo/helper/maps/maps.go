package maps

// Contains 判断 m1 是否包含 m2，二者有任意一个为 nil，或长度为空，则返回 false
func Contains[M ~map[K]V, K comparable, V comparable](m1, m2 M) bool {
	if m1 == nil || len(m1) == 0 ||
		m2 == nil || len(m2) == 0 {
		return false
	}
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

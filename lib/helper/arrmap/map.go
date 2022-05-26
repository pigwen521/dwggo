package arrmap

//key是否存在-string
func InMapStrKey(key string, arrs *map[string]string) bool {
	for k := range *arrs {
		if key == k {
			return true
		}

	}
	return false

}

//val是否存在-string
func InMapStrVal(val string, arrs *map[string]string) bool {
	for _, v := range *arrs {
		if val == v {
			return true
		}
	}
	return false
}

//两个map合并
func MergeMap(m1, m2 *map[string]interface{}) *map[string]interface{} {
	for k, v := range *m2 {
		(*m1)[k] = v
	}
	return *&m1
}

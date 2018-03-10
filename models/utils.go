package models

func stringSliceToMap(slice []string) map[string]interface{} {
	m := map[string]interface{}{}
	for _, s := range slice {
		m[s] = true
	}
	return m
}

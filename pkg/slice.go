package pkg

import "slices"

func Contains[T comparable](slice []T, item T) bool {
	//for _, i := range slice {
	//	if i == item {
	//		return true
	//	}
	//}
	//return false
	return slices.Contains(slice, item)
}

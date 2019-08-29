package helper

func ArrayRemove(array []int64, value int64) []int64 {
	newArray := make([]int64, 0)
	for _, item := range array {
		if item != value {
			newArray = append(newArray, item)
		}
	}
	return newArray
}

package utils

func CopyMap[S map[E]E, E int](originalMap S) S {
	newMap := make(S, len(originalMap))
	for i, v := range originalMap {
		newMap[i] = v
	}

	return newMap
}

func GetAllMapKeys[S map[E]E, E int](originalMap S) []E {
	keys := make([]E, 0, len(originalMap))
	for key := range originalMap {
		keys = append(keys, key)
	}

	return keys
}

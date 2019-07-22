// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

func sliceToMap(inputSlice []string) map[string]string {
	outputMap := make(map[string]string)
	for _, item := range inputSlice {
		outputMap[item] = item
	}
	return outputMap
}

func mapToSlice(inputMap map[string]string) []string {
	outputSlice := make([]string, 0, len(inputMap))

	for _, value := range inputMap {
		outputSlice = append(outputSlice, value)
	}

	return outputSlice
}

func addSliceToMap(inputSlice []string, inputMap map[string]string) map[string]string {
	for _, item := range inputSlice {
		inputMap[item] = item
	}
	return inputMap
}

func removeSliceFromMap(inputslice []string, inputMap map[string]string) map[string]string {
	for _, item := range inputslice {
		delete(inputMap, item)
	}
	return inputMap
}

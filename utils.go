package main

import "strings"

func stringInSlice(a string, list []interface{}) bool {
	for _, b := range list {
		if strings.Contains(a, b.(string)) {
			return true
		}
	}
	return false
}

package g

import (
	"regexp"
	"strings"
)

func KeyInMapWithFloat(s map[string]float64, j string) ( exist bool) {

	for keyInMap, _ := range s {
		match, _ := regexp.MatchString(j, keyInMap)
		if match {
			return true
		}
	}

	return false
}


func GetMetricAndTag(s map[string]float64, j string) (tag []string) {

	for key, _ := range s {

		keySplit := strings.Split(key, "@")

		if len(keySplit) == 2 {
			if keySplit[0] == j {
				tag = append (tag, keySplit[1])
			}
		}

	}

	return
}
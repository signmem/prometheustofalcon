package g

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var floatType = reflect.TypeOf(float64(0))

func KeyInMapWithFloat(s map[string]float64, j string) ( exist bool) {

	for keyInMap, _ := range s {
		match, _ := regexp.MatchString(j, keyInMap)
		if match {
			return true
		}
	}

	return false
}

func KeyinSliceWithChar(s []string, j string) (exit bool) {

	for _, i := range s {
		if i == j {
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

func LoadMetricJsonFile (filePath string) (validMetric []string) {
	if _, err := os.Stat(filePath); err == nil {

		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			Logger.Errorf("loadjsonfile %s read error", filePath)
			return
		}

		err = json.Unmarshal(content, &validMetric)

		if err != nil {
			Logger.Errorf("loadjsonfile %s json format error", filePath)
			return
		}

	} else {
		Logger.Errorf("loadjsonfile %s not exists", filePath)
	}
	return
}

func LoadCalMetricJsonFile(filePath string) (matchMetric []string, calMetric []MetricCalType) {
	if _, err := os.Stat(filePath); err == nil {

		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			Logger.Errorf("cal metrics jsonfile %s read error", filePath)
			return
		}

		err = json.Unmarshal(content, &calMetric)

		if err != nil {
			Logger.Errorf("cal metrics jsonfile %s json format error", filePath)
			return
		}

		for _, metrics := range calMetric {
			matchMetric = append(matchMetric, metrics.MetricSum)
			matchMetric = append(matchMetric, metrics.MetricCount)
		}

	} else {
		Logger.Errorf("metrics jsonfile %s not exists", filePath)
	}
	return
}


func GetValueToFloat(unk interface{}) (float64) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0
	}
	fv := v.Convert(floatType)
	return fv.Float()
}

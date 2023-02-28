package prom

import (
	"bytes"
	"fmt"
	"github.com/signmem/prometheustofalcon/g"
	"github.com/signmem/prometheustofalcon/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetMetricFromPrometheus() (allMetric []MetricValue) {
	server := g.Config().MetricServer.Server
	port := g.Config().MetricServer.Port
	metricAPI := g.Config().MetricServer.MetricAPI

	url := server + ":" + port + metricAPI
	metricFromHTTP, err := http.HttpApiGet(url, "","")
	if err != nil {
		g.Logger.Errorf("GetMetricFromPrometheus() error:%s", err)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(metricFromHTTP)
	newStr := buf.String()
	responseString := strings.Split( newStr, "\n")

	var metric MetricValue
	var totalline TotalLine

	timeNow := time.Now().Unix()
	metric.Timestamp = timeNow
	metric.Step = g.Config().Falcon.Step
	metric.Endpoint = g.Config().Falcon.Endpoint
	metric.Type = "GAUGE"

	for _, line := range responseString {
		if line == "" ||  strings.Contains(line, "#") {
			continue
		}

		lineSp := strings.Split(line, "}")

		//fmt.Println(len(lineSp))
		// continue
		if len(lineSp) < 2 {
			specialLine := strings.Split(line, " ")
			if len(specialLine) == 2  {
				metricName := specialLine[0]
				metric.Metric = MKmetric(metricName)
				value := strings.Replace(specialLine[1]," ", "", -1)
				metric.Value, _ = strconv.ParseFloat(value,  64)
				if metric.Value != metric.Value {
					continue
				}
				allMetric = append(allMetric, metric)
				continue
			}
		}
		// fmt.Println(line)
		totalline.Info, totalline.Value = lineSp[0], strings.Replace(lineSp[1],
			" ", "", -1)

		lineSp2 := strings.Split(totalline.Info, "{")
		var tags string
		metricName, tags := lineSp2[0], lineSp2[1]
		replacer := strings.NewReplacer( "\"", "", "}","", "pool=", "ceph_pool=")
		metric.Tags = replacer.Replace(tags)

		metric.Metric = MKmetric(metricName)

		metric.Value, _ = strconv.ParseFloat(totalline.Value,  64)
		// fmt.Println(metric)
		allMetric = append(allMetric, metric)
	}
	return allMetric
}


func GetProms() {
	for {
		getAllMetric := GetMetricFromPrometheus()
		error := SendMetric(getAllMetric)
		if error != nil {
			fmt.Println(error)
		}
		time.Sleep(time.Duration(g.Config().Falcon.Step) * time.Second)
	}
}

func MKmetric(metric string) (newmetric string) {
	status,_ :=  regexp.MatchString("^ceph*", metric)
	if status == false {
		if g.Config().Grafana {
			newmetric = "ceph." + strings.Replace(metric, "_", ".", -1)
		} else {
			newmetric = "ceph_" + metric
		}
	} else {
		if g.Config().Grafana {
			newmetric = strings.Replace(metric, "_", ".", -1)
		} else {
			newmetric = metric
		}
	}
	return newmetric
}
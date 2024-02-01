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

				// use to get mds_sum and dms_count metric info without tag

				if g.Config().MdsEnable {
					loadMdsName := g.Config().MdsMetric
					for _, metricDetail := range loadMdsName {

						if metricName == metricDetail.MetricSum ||
							metricName == metricDetail.MetricCount {
							g.MdsMetricNew[metricName], _ =
								strconv.ParseFloat(value,  64)
						}
					}
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
		var mdsvalue float64
		metricName := lineSp2[0]
		tags = lineSp2[1]

		if metricName == "ceph_mds_metadata" {

			metadataSplit := strings.Split(tags, ",")

			if len(metadataSplit) > 0 {
				for _, keyName := range metadataSplit {
					splitKey := strings.Split(keyName, "=")
					if splitKey[0] == "rank" {
						replacer := strings.NewReplacer( "\"", "", "}","", "pool=", "ceph_pool=")
						tagsvalue := replacer.Replace(splitKey[1])
						mdsvalue, _ = strconv.ParseFloat(tagsvalue, 64)
						metric.Value = mdsvalue
					}

					if splitKey[0] == "ceph_daemon" {
						replacer := strings.NewReplacer( "\"", "", "}","", "pool=", "ceph_pool=")
						tagsvalue := replacer.Replace(splitKey[1])
						metric.Tags = "ceph_daemon=" + tagsvalue
						metric.Type = "COUNTER"
					}
				}
				metric.Metric = "ceph_mds_status_change"
				allMetric = append(allMetric, metric)


				if mdsvalue >= 0 {
					metric.Metric = "ceph_mds_status_active"
					metric.Value = 1
				} else {
					metric.Metric = "ceph_mds_status_backup"
					metric.Value = 1
				}
				metric.Type = "GAUGE"
				allMetric = append(allMetric, metric)
				continue
			}
		} else {

			replacer := strings.NewReplacer( "\"", "", "}","", "pool=", "ceph_pool=")
			metric.Tags = replacer.Replace(tags)

			if metricName == "ceph_mgr_status" {
				metric.Metric = "ceph_mgr_status_change"
			} else {
				metric.Metric = MKmetric(metricName)
			}

			metric.Value, _ = strconv.ParseFloat(totalline.Value,  64)
		}

		// fmt.Println(metric)
		allMetric = append(allMetric, metric)


		// use to get mds_sum and dms_count metric info with tag

		if g.Config().MdsEnable {
			loadMdsName := g.Config().MdsMetric
			for _, metricDetail := range loadMdsName {

				if metricName == metricDetail.MetricSum ||
					metricName == metricDetail.MetricCount {

						metricWithTag := metricName + "@" + tags
						var metricValue float64
						metricValue, _ = strconv.ParseFloat(totalline.Value,64)

						if metricValue != metricValue {
							continue
						}

						g.MdsMetricNew[metricWithTag] = metricValue

				}
			}
		}


	}

	GetMdsCalAvg()

	if g.Config().Debug {
		g.Logger.Infof("new metrics %d, old metrids %d", len(g.MdsMetricNew), len(g.MdsMetricNew))
		for _, detail := range allMetric {
			g.Logger.Debugf("%s", detail.String())
		}
	}

	if  len(g.MdsMetricNew) > 0 {
		g.MdsMetricOld = make(map[string]float64)
		g.MdsMetricOld  = g.MdsMetricNew
		g.MdsMetricNew = make(map[string]float64)
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
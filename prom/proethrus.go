package prom

import (
	"bytes"
	"fmt"
	"github.com/signmem/prometheustofalcon/g"
	"github.com/signmem/prometheustofalcon/http"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetMetricFromPrometheus() (allMetric []MetricValue) {
        // 用于获取 /metrics 信息
        // 通过 /metrics 中 TYPE 中信息定义 type, 包含 (counter gauge histogram summary) || untyped

	server := g.Config().MetricServer.Server
	port := g.Config().MetricServer.Port
	metricAPI := g.Config().MetricServer.MetricAPI

	var metricFromHTTP io.ReadCloser
	var err error

	if g.Config().SslEnable == false {

		url := "http://" + server + ":" + port + metricAPI
		metricFromHTTP, err = http.HttpApiGet(url, "","")

		if err != nil {
			g.Logger.Errorf("GetMetricFromPrometheus() error:%s", err)
			return
		}
	}

	if g.Config().SslEnable == true {

		url := "https://" + server + ":" + port + metricAPI
		metricFromHTTP, err = http.HttpsApiGet(url, "")

		if err != nil {
			g.Logger.Errorf("GetMetricFromPrometheus() error:%s", err)
			return
		}
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(metricFromHTTP)
	newStr := buf.String()
	responseString := strings.Split(newStr, "\n")

	if g.Config().Debug {
		g.Logger.Debugf("/metrics respone lines :%d", len(responseString))
	}

	var metric MetricValue
	var totalline TotalLine

	timeNow := time.Now().Unix()
	metric.Timestamp = timeNow
	metric.Step = g.Config().Falcon.Step
	metric.Endpoint = g.Config().Falcon.Endpoint
	metric.Type = "GAUGE"

	for _, line := range responseString {

		if line == ""  {
			continue
		}

		// get metrics from line
		if strings.Contains(line, "#") {
			r, _ := regexp.Compile("^#\\ +TYPE")
			if r.MatchString(line) == true {
				metric_type := strings.ToUpper(line[strings.LastIndex(line, " ") +1:])

				if metric.Type == "UNTYPED" {
					metric_type = "GAUGE"
				}
				metric.Type = metric_type
			}
			// ceph /metrics 中有 untype 类型，暂配置称 GAUGE
			continue
		}

		lineSp := strings.Split(line, "}")

		if len(lineSp) == 1 {

			// 只针对不带 tags 的 metrics 进行处理

			specialLine := strings.Split(line, " ")
			if len(specialLine) == 2  {
				metricName := specialLine[0]
				value := strings.Replace(specialLine[1]," ", "", -1)
				metric.Value, _ = strconv.ParseFloat(value,  64)

				if metric.Value != metric.Value {
					continue
				}

				// use to get mds_sum and mds_count metric info without tag

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

				metric.Metric = MKmetric(metricName)
				allMetric = append(allMetric, metric)
				continue
			}
		}


		// g.Logger.Debugf("debug test line %s", lineSp)
		// continue

		totalline.Info, totalline.Value = lineSp[0], strings.Replace(lineSp[1],
			" ", "", -1)

		lineSp2 := strings.Split(totalline.Info, "{")

		var tags string
		var mdsvalue float64

		if len(lineSp2) == 0 {
			continue
		}

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

			replacer := strings.NewReplacer( "\"", "", "}","")
			metric.Tags = replacer.Replace(tags)

			if metricName == "ceph_mgr_status" {
				metric.Metric = "ceph_mgr_status_change"
			}

			metric.Value, _ = strconv.ParseFloat(totalline.Value,  64)
		}

		metric.Metric = MKmetric(metricName)

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

				GetMdsCalAvg()
			}
		}
	}



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

	// 用法: 如果配置了 grafana = true	
	// 把所有 _ 转换称为 .
	// prometheus:/metrics 默认 _ 作为连接符，无需处理

	if g.Config().Grafana {
		newmetric = strings.Replace(metric, "_", ".", -1)
	} else {
		newmetric = metric
	}
	return newmetric
}

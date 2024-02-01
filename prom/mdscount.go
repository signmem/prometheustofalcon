package prom

import (
	"fmt"
	"github.com/signmem/prometheustofalcon/g"
	"strings"
	"time"
)


func GetMdsCalAvg() {

	metrics := mdsCalAvg()

	if g.Config().Debug {
		g.Logger.Debugf("len of avg:% d", len(metrics))
	}

	error := SendMetric(metrics)

	if error != nil {
		fmt.Println(error)
	}

}


func mdsCalAvg() (metricAvgs []MetricValue){

	// use to calcutlate mds avg load

	var metricAvg MetricValue

	timeNow := time.Now().Unix()
	metricAvg.Timestamp = timeNow
	metricAvg.Step = g.Config().Falcon.Step
	metricAvg.Endpoint = g.Config().Falcon.Endpoint
	metricAvg.Type = "GAUGE"



	if len(g.MdsMetricNew) > 0 && len(g.MdsMetricOld) > 0 {

		metricsInfo := g.Config().MdsMetric

		for  _, metricsDetail := range metricsInfo {

			sum := metricsDetail.MetricSum
			count := metricsDetail.MetricCount
			name := metricsDetail.MetricName

			if g.KeyInMapWithFloat(g.MdsMetricNew, sum) && g.KeyInMapWithFloat(g.MdsMetricNew, count) &&
				g.KeyInMapWithFloat(g.MdsMetricOld, sum) && g.KeyInMapWithFloat(g.MdsMetricOld, count) {

				metricAvg.Metric = name

				tags := g.GetMetricAndTag(g.MdsMetricNew, sum)

				if len(tags) > 0 {
					for _, tag := range tags {

						metricNameSum := sum + "@" + tag
						metricNameCount := count + "@" + tag

						sumValue := g.MdsMetricNew[metricNameSum] - g.MdsMetricOld[metricNameSum]
						countValue := g.MdsMetricNew[metricNameCount] - g.MdsMetricOld[metricNameCount]

						/*
						if g.Config().Debug {
							if name == "ceph_mds_reply_latency_avg" {
								g.Logger.Debugf("metric:%s, newsum: %v", name, g.MdsMetricNew[metricNameSum] )
								g.Logger.Debugf("metric:%s, oldsum: %v", name, g.MdsMetricOld[metricNameSum])
								g.Logger.Debugf("metric:%s, newcount: %v", name, g.MdsMetricNew[metricNameCount])
								g.Logger.Debugf("metric:%s, oldcount: %v", name, g.MdsMetricOld[metricNameCount] )
							}
						}
						*/


						if countValue == 0 || sumValue == 0  {
							metricAvg.Value = 0
						} else {
							metricAvg.Value = sumValue / countValue
						}

						replacer := strings.NewReplacer( "\"", "", "}","", "pool=", "ceph_pool=")
						metricAvg.Tags = replacer.Replace(tag)

						metricAvgs = append(metricAvgs, metricAvg)
					}
				} else {
					sumValue := g.MdsMetricNew[sum] - g.MdsMetricOld[sum]
					countValue := g.MdsMetricNew[count] - g.MdsMetricOld[count]
					if countValue == 0 || sumValue == 0  {
						metricAvg.Value = 0
					} else {

						if g.Config().Debug {
							g.Logger.Debugf( "metric: %s, sum:%v, count%v",sum, sumValue, countValue )
						}

						metricAvg.Value = sumValue / countValue
					}

					metricAvgs = append(metricAvgs, metricAvg)
				}

			}
		}



		if g.Config().Debug {

			for _, metrics := range metricAvgs {
				g.Logger.Debugf("metrics info: %s", metrics.String())
			}

		}


	}


	return metricAvgs

}
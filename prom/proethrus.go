package prom

import (
	"fmt"
	"github.com/coreos/go-log/log"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/signmem/prometheustofalcon/g"
	"github.com/signmem/prometheustofalcon/http"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

var (
	STEP  int64
	HOSTNAME string
	ValidMetric []string    // 用于过滤并只收集当前 list 中的 metrics 其他 metrics 不需要上报
	MatchCalMetricList []string	// 用于匹配自身计算用的 metrics, 匹配后，更新 CalMetric 中的值
	CalMetricDict  []g.MetricCalType	// 用于自身 metrics 计算用的常量，被动更新
	SumMetrics		[]string  // use to sum metrics values
	AllMetrics []string   // use to match metrics
)

func getMetricFromServer() (httpresponse string, err error) {

	server := g.Config().MetricServer.Server
	port := g.Config().MetricServer.Port
	metricAPI := g.Config().MetricServer.MetricAPI

	var metricFromHTTP io.ReadCloser

	if g.Config().SslEnable == false {

		url := "http://" + server + ":" + port + metricAPI
		metricFromHTTP, err = http.HttpApiGet(url, "","")

		if err != nil {
			g.Logger.Errorf("GetMetricFromPrometheus() error:%s", err)
			return "", err
		}
	}

	if g.Config().SslEnable == true {

		url := "https://" + server + ":" + port + metricAPI
		metricFromHTTP, err = http.HttpsApiGet(url, "")

		if err != nil {
			g.Logger.Errorf("GetMetricFromPrometheus() error:%s", err)
			return "", err
		}
	}

	responseBody, err := ioutil.ReadAll(metricFromHTTP)

	return string(responseBody), nil
}


func GetMetricFromPrometheus() (getAllMetric []MetricValue) {
	// 用于获取 /metrics 信息
	// 通过 /metrics 中 TYPE 中信息定义 type, 包含 (counter gauge histogram summary) || untyped

	metricsString, err := getMetricFromServer()  // get all info from html
	if err != nil {
		return
	}

	if len(metricsString) == 0 {
		return
	}

	getAllMetric, matchMetricDict, err := genMetricFormat(metricsString)

	if err != nil {
		return
	}

	if  len(MatchCalMetricList) > 0 && len(CalMetricDict) > 0 {
		calMatchMetrics := calMatchMetricValues(CalMetricDict, matchMetricDict)
		getAllMetric = append(getAllMetric, calMatchMetrics...)
	}

	g.Logger.Debugf("total metrics: %d\n", len(getAllMetric))

	if g.Config().Falcon.Enable {
		_ = SendMetric(getAllMetric)
	}

	return
}

func genMetricFormat(info string) (getAllMetric []MetricValue, matchMetricDict []MetricValue, err error) {

	// getAllMetric  all of the metrics (包含了计算前的 metrics)
	// matchMetricDict 只包含用于计算用的 metrics 信息

	timestamp := time.Now().Unix()

	var metric MetricValue

	// totalCountMetric  use to auto sum values
	var totalCountMetric []MetricValue

	parser := &expfmt.TextParser{}
	families, err := parser.TextToMetricFamilies(strings.NewReader(info))

	if err != nil {
		log.Errorf("failed to parse input: %w", err)
		return getAllMetric, matchMetricDict,  err
	}

	for _, val := range families {

		for _, m := range val.GetMetric() {

			metric.Metric = val.GetName()

			// 只有自定义需要获取的 metrics 时候才需要执行下面过滤操作
			// 如果不过滤，则获取所有 metrics
			if (len(AllMetrics)) > 0 {
				if g.KeyinSliceWithChar(AllMetrics, metric.Metric) == false {
					continue
				}
			}

			// 或者 mertric 对应 values, 根据不同类型进行判定
			switch val.GetType() {
			case dto.MetricType_COUNTER:
				metric.Value = m.GetCounter().GetValue()
				metric.Type = "COUNTER"
			case dto.MetricType_GAUGE:
				metric.Value = m.GetGauge().GetValue()
				metric.Type = "GAUGE"
			case dto.MetricType_UNTYPED:
				metric.Value = m.GetUntyped().GetValue()
				metric.Type = "GAUGE"
			case dto.MetricType_SUMMARY:
				metric.Value = m.GetSummary().GetSampleSum()
				metric.Type = "SUMMARY"
			default:
				// 部分 metrics type unkonw , 暂时作为 gauge type
				metric.Value = m.GetGauge().GetValue()
				metric.Type = "GAUGE"
			}

			metric.Tags = ""
			var metricTags string

			metric.Metric = val.GetName()
			metric.Step = STEP
			metric.Endpoint = HOSTNAME
			metric.Timestamp = timestamp

			// 获取指标的 metrics

			for n, label := range m.GetLabel() {
				if n == len(m.Label) - 1 {
					tag := fmt.Sprintf("%s=%s" , label.GetName(), label.GetValue())
					metric.Tags = metricTags + tag
				} else {
					// 多个 tags 处理方法
					tag := fmt.Sprintf("%s=%s," , label.GetName(), label.GetValue())
					metric.Tags = metricTags + tag
				}
			}

			// 对需要计算的 metric 进行收集， 放外部进行处理 (matchMetricDict)
			if (len(MatchCalMetricList) > 0 ){
				for _, name := range MatchCalMetricList {
					if name == metric.Metric {
						matchMetricDict = append(matchMetricDict, metric)
					}
				}
			}

			// 需要对计算 count 总数的 metric 进行处理， 避免过多 metric 进行上报
			var appendGrap bool
			if (len(SumMetrics) > 0) {
				for _, name := range SumMetrics {
					if metric.Metric == name {
						totalCountMetric = metricAddSlice(metric, totalCountMetric)

						// not going to send
						appendGrap = true
						break
					}
				}
			}

			if appendGrap == true {
				continue
			}

			getAllMetric = append(getAllMetric, metric)
		}
	}

	g.Logger.Debugf("totalCountMetric len:%d\n", len(totalCountMetric))
	
	// count metrics 只需要在循环结束后放入 getAllMetric 变量中
	getAllMetric = append(getAllMetric, totalCountMetric...)

	if g.Config().Debug {
		for _, metric := range getAllMetric {
			g.Logger.Debugf("%s", metric.String())
		}
	}

	return getAllMetric, matchMetricDict,nil
}


func metricAddSlice(newStruct MetricValue, w []MetricValue)(q []MetricValue) {

	// 用于自动化把 metricValue 加入 []MetricValue
	// 自动化对 value 进行 sum 计算

	if metricInSlice(newStruct, w) == true {
		q = metricAddValue(newStruct, w)
	} else {
		q = append(w, newStruct)
	}
	return q
}

func metricInSlice(newStruct MetricValue, w []MetricValue) bool {
	// use to judge metric in []metrics

	for _, info := range w {
		if info.Metric == newStruct.Metric {
			return true
		}
	}
	return false
}


func metricAddValue(newStruct MetricValue, w []MetricValue)(q []MetricValue) {

	for _, metric := range w {
		if metric.Metric == newStruct.Metric {
			metric.Add(newStruct.Value)
		}
		q = append(q, metric)
	}
	return
}



func GetProms() {
	for {

		if time.Now().Unix() % g.Config().Falcon.Step == 0 {
			_ = GetMetricFromPrometheus()
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func calMatchMetricValues( calMetricDict []g.MetricCalType,
	calMetric []MetricValue) (calMatchMetric []MetricValue) {

		// calMetricList = list == prom.MatchCalMetricList
		// calMetricDict = dict == prom.CalMetricDict
		// calMetric = 只包含 prom.MatchCalMetricList 中的所有 metrics 信息  
		// calMatchMetric = 通过计算 metricsum / metriccount 对应的 metricname

		for _, calInfo := range calMetricDict {

			sumName := calInfo.MetricSum
			var tags []string

			for _, metricInfo := range  calMetric {

				if metricInfo.Metric == sumName {

					if g.KeyinSliceWithChar(tags, metricInfo.Tags) == false {
						tags = append(tags, metricInfo.Tags)
					}
				}
			}

			for _, tag := range tags {

				var newMetric MetricValue
				var countName, metricName string
				var countVale, sumValue interface{}

				for _, info := range calMetricDict {

					if info.MetricSum == sumName {
						countName = info.MetricCount
						metricName = info.MetricName
						break
					}
				}

				newMetric.Metric = metricName

				for _, metricInfo := range calMetric {
					if metricInfo.Metric == sumName && metricInfo.Tags == tag {
						sumValue = metricInfo.Value
					}

					if metricInfo.Metric == countName && metricInfo.Tags == tag {
						countVale = metricInfo.Value
						newMetric.Endpoint = metricInfo.Endpoint
						newMetric.Timestamp = metricInfo.Timestamp
						newMetric.Step = metricInfo.Step
						newMetric.Type = metricInfo.Type
					}
				}

				if g.GetValueToFloat(countVale) == 0 {
					newMetric.Value = float64(0)

				} else {
					newMetric.Value = g.GetValueToFloat(sumValue) / g.GetValueToFloat(countVale)
				}

				newMetric.Tags = tag

				calMatchMetric = append(calMatchMetric, newMetric)
			}

		}

		return calMatchMetric
}

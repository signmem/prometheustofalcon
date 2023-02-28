package prom

import (
	"github.com/signmem/prometheustofalcon/g"
	"github.com/signmem/prometheustofalcon/http"
	"encoding/json"
)

func SendMetric(metrics []MetricValue) (err error) {

	httpUrl := g.Config().Falcon.Url + g.Config().Falcon.Api

	if g.Config().Debug {
		for _, metric := range metrics {
			g.Logger.Debugf("metric: %s", metric.String())
		}
	}

	metricsBytes, err := json.Marshal(metrics)
	if err != nil {
		g.Logger.Errorf("SendMetric() json error %s", err)
		return err
	}

	resp, err := http.HttpApiPost(httpUrl, metricsBytes, "")
	if err != nil {
		g.Logger.Errorf("SendMetric() http post error %s", err)
		return err
	}
	err = resp.Close()
	if err != nil {
		g.Logger.Errorf("SendMetric() http close error %s", err)
		return err
	}
	return nil


}
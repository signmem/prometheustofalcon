package prom

import (
	"fmt"
	"github.com/signmem/prometheustofalcon/g"
)

type Metric struct {
	Name string
	Tags string
	Value float64
}

type TotalLine struct {
	Info string
	Value string
}



type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

func (this *MetricValue) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Type:%s, Tags:%s, Step:%d, " +
			"Time:%d, Value:%v>",
		this.Endpoint,
		this.Metric,
		this.Type,
		this.Tags,
		this.Step,
		this.Timestamp,
		this.Value,
	)
}

func (this *MetricValue) Add(value interface{}) ( metric MetricValue) {
	this.Value = g.GetValueToFloat(this.Value) + g.GetValueToFloat(value)
	return
}

type MetaData struct {
	Metric      string            `json:"metric"`
	Endpoint    string            `json:"endpoint"`
	Timestamp   int64             `json:"timestamp"`
	Step        int64             `json:"step"`
	Value       float64           `json:"value"`
	CounterType string            `json:"counterType"`
	Tags        map[string]string `json:"tags"`
}

func (t *MetaData) String() string {
	return fmt.Sprintf("<MetaData Endpoint:%s, Metric:%s, Timestamp:%d, Step:%d, Value:%f, Tags:%v>",
		t.Endpoint, t.Metric, t.Timestamp, t.Step, t.Value, t.Tags)
}

type MetricCalType struct {
	MetricSum	string		`json:"metricsum"`
	MetricCount	string		`json:"metriccount"`
	MetricName 	string 		`json:"metricname"`
}


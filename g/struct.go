package g

type GlobalConfig struct {
	Debug			bool		`json:"debug"`
	LogFile			string		`json:"logfile"`
	LogMaxAge		int			`json:"logmaxage"`
	LogRotateAge	int			`json:"logrotateage"`
	Grafana  		bool 		`json:"grafana"`
	MetricServer    *MetricDetail 		`json:"metricserver"`
	Falcon 			*Falcon 	`json:"falcon"`
	MdsEnable 		bool		`json:"mdsenable"`
	MdsMetric 		[]*MdsMetric 	`json:"mdsmetric"`
}

type MdsMetric struct {
	MetricSum 		string		`json:"metricsum"`
	MetricCount 	string 		`json:"metriccount"`
	MetricName		string		`json:"metricname"`
}

type MetricDetail struct {
	Enable 			bool 		`json:"enable"`
	Server 			string 		`json:"server"`
	Port 			string 		`json:"port"`
	MetricAPI       string  	`json:"metricapi"`
}

type Falcon struct {
	Step		int64 		`json:"step"`
	Endpoint	string		`json:"endpoint"`
	Url 		string 		`json:"url"`
	Api 		string 		`json:"api"`
}
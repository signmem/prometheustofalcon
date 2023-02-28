package g

type GlobalConfig struct {
	Debug			bool		`json:"debug"`
	LogFile			string		`json:"logfile"`
	LogMaxAge		int			`json:"logmaxage"`
	LogRotateAge	int			`json:"logrotateage"`
	Grafana  		bool 		`json:"grafana"`
	MetricServer    *MetricDetail 		`json:"metricserver"`
	Falcon 			*Falcon 	`json:"falcon"`
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
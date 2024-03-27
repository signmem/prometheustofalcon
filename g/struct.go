package g

type GlobalConfig struct {
	Debug			bool			`json:"debug"`
	LogFile			string			`json:"logfile"`
	ValidMetricFile		string			`json:"validmetricfile"`
	CalMetricFile		string			`json:"calmetricfile"`
	SumMetricFile 	string			`json:"summetricfile"`
	LogMaxAge		int				`json:"logmaxage"`
	LogRotateAge	int				`json:"logrotateage"`
	MetricServer    *MetricDetail 	`json:"metricserver"`
	Falcon 			*Falcon 		`json:"falcon"`
	SslEnable		bool 			`json:"sslenable"`
	TLS				*TLS			`json:"tls"`
}

type TLS struct {
	CaFile 			string 			`json:"cafile"`
	CertFile 		string 			`json:"certfile"`
	KeyFile 		string			`json:"keyfile"`
}

type CountMetric struct {
	Value 			interface{}		`json:"value"`
	MetricName		string			`json:"metric"`
}


type MetricDetail struct {
	Enable 			bool 		`json:"enable"`
	Server 			string 		`json:"server"`
	Port 			string 		`json:"port"`
	MetricAPI       string  	`json:"metricapi"`
}

type Falcon struct {
	Enable 		bool 		`json:"enable"`
	Step		int64 		`json:"step"`
	Url 		string 		`json:"url"`
	Api 		string 		`json:"api"`
}


type MetricCalType struct {
	MetricSum	string		`json:"metricsum"`
	MetricCount	string		`json:"metriccount"`
	MetricName 	string 		`json:"metricname"`
}
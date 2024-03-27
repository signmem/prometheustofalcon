# prometheustofalcon

> use to set golang read cfg.json template.
> log info into logfile.


# log vendor

>  确保 vendor/github.com/coreos/go-log/log/fields.go 文件被修改
> 1 修改 full_time 满足格式要求

```
"full_time":  time.Now().Format("2006-01-02 15:04:05.999"),  // time of log entry

```

> 2 修改 logger.verbose = true  由于属于内部变量无法外部修改

```
logger.verbose = true
```
> 3 logger example 

```
[2022-04-26 18:12:46.65] [DEBUG] [28062] [commands.go:33] >>> [main] msg=debug: yes
[2022-04-26 18:13:30.858] [DEBUG] [28224] [commands.go:33] >>> [main] msg=debug: yes
[2022-04-26 18:13:30.858] [INFO] [28224] [commands.go:33] >>> [main] msg=info: yes
[2022-04-26 18:13:30.858] [WARNING] [28224] [commands.go:33] >>> [main] msg=warning: yes
```

# 架构说明  

> 因为 ceph 每个版本暴露的 /metrics 格式都不一样    
> 因此通过 ceph_exporter 对 ceph 进行监控，统一 metric name 格式    
> ceph_exporter 通过 /metrics api 提供监控信息    
> 当 prometheustofalcon 访问  ceph_exporter:port/metrics 即达到获取当前 ceph 监控数据状态    
> 获取信息后转换为 falcon 上报格式，并向 falcon agent 上报    
> 每个 ceph 集群只需要跑一套  ceph_exporter, prometheustofalcon 即可  

# 新功能   
> 1. /metrics 包含了大量的数据，但很多监控指标并不需要，因此增加自定义监控指标功能，只录入指定的指标，不录入则丢弃  
> 2. /metrics 例如 ceph_osd_up 指标，只需要获取 sum(ceph_osd_up) 一条数据，但集群有 5000 个 OSD, 则上报了很多多余指标，因此程序增加对某个指标进行 sum 统计功能    
> 3. /etrics 例如 需要对 metric1 / metric2 后的结果进行报警，falcon 当前没有类似功能，因此程序增加自行计算功能    

## 用法说明   
> 指定执行监控指标方法, 配置中通过 validmetricfile 指定 json 文件，把需要上报的 metrics 自行定义到 json list 中      
example   
```
[
        "ceph_cluster_total_used_bytes",
        "ceph_cluster_total_bytes"
]
```

> 对特定指标执行 sum 功能, 配置中通过 summetricfile 指定 json 文件，把需要上报的 metrics 自行定义到 json list 中   
example   
```
[
        "ceph_osd_op_r",
        "ceph_osd_op_w",
	"ceph_osd_up",
        "ceph_osd_in"
]
```   

> 指定自定义计算指标方法，配置通过 calmetricfile 定义 json 文件，把 metric 定义到 json 文件中  
> 即可通过 "metricsum"/"metriccount"  并把值赋至 "metricname" 这个新的 metric 中   
example   
```
[
        {
            "metricsum": "ceph_mds_server_req_setfilelock_latency_sum",
            "metriccount": "ceph_mds_server_req_setfilelock_latency_count",
            "metricname": "ceph_mds_server_req_setfilelock_avg"
        }
]
```

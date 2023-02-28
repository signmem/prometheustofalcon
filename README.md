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

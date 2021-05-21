package eskywalking

import (
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"log"
	"os"
)

// Load config from file
func Load(key string) *config {
	config := DefaultConfig()
	if err := econf.UnmarshalKey(key, config); err != nil {
		elog.Panic("unmarshal key", elog.FieldErr(err))
	}
	return config
}

func (config *config) Build() *go2sky.Tracer {
	rpcOptions := make([]reporter.GRPCReporterOption, 0, 8)
	if config.GRPCReporterOptions != nil {
		if config.GRPCReporterOptions.Log != nil {
			file, err := os.Open(config.GRPCReporterOptions.Log.FilePath)
			if err != nil {
				elog.Panic("open file", elog.FieldErr(err))
			}
			defer file.Close()
			logger := log.New(file, config.GRPCReporterOptions.Log.Prefix, config.GRPCReporterOptions.Log.Flag)
			rpcOptions = append(rpcOptions, reporter.WithLogger(logger))
		}
		if config.GRPCReporterOptions.Auth != "" {
			rpcOptions = append(rpcOptions, reporter.WithAuthentication(config.GRPCReporterOptions.Auth))
		}
		if config.GRPCReporterOptions.CheckInterval != 0 {
			rpcOptions = append(rpcOptions, reporter.WithCheckInterval(config.GRPCReporterOptions.CheckInterval))
		}
		//if config.GRPCReporterOptions.CDSInterval != 0 {
		//	rpcOptions = append(rpcOptions, reporter.WithCDS(config.GRPCReporterOptions.CDSInterval))
		//}
		if config.GRPCReporterOptions.InstanceProps != nil {
			rpcOptions = append(rpcOptions, reporter.WithInstanceProps(config.GRPCReporterOptions.InstanceProps))
		}
		if config.GRPCReporterOptions.MaxSendQueueSize >= 0 {
			rpcOptions = append(rpcOptions, reporter.WithMaxSendQueueSize(config.GRPCReporterOptions.MaxSendQueueSize))
		}
	}
	r, err := reporter.NewGRPCReporter(config.ServerAddr, rpcOptions...)
	if err != nil {
		elog.Panic("build reporter", elog.FieldErr(err))
	}

	tracerOptions := make([]go2sky.TracerOption, 0, 6)
	tracerOptions = append(tracerOptions, go2sky.WithReporter(r))
	if config.TracerOptions != nil {
		if config.TracerOptions.Instance != "" {
			tracerOptions = append(tracerOptions, go2sky.WithInstance(config.TracerOptions.Instance))
		}
		if config.TracerOptions.SamplingRate >= 0 {
			tracerOptions = append(tracerOptions, go2sky.WithSampler(config.TracerOptions.SamplingRate))
		}
		if config.TracerOptions.KeyCount >= 0 && config.TracerOptions.ValueSize >= 0 {
			tracerOptions = append(tracerOptions, go2sky.WithCorrelation(config.TracerOptions.KeyCount, config.TracerOptions.ValueSize))
		}
	}
	tracer, err := go2sky.NewTracer(config.ServiceName, tracerOptions...)
	if err != nil {
		elog.Panic("build tracer", elog.FieldErr(err))
	}
	return tracer
}


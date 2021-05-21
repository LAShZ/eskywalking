package eskywalking

import "time"

type LogOptions struct {
	FilePath string
	Prefix string
	Flag int
}

type GRPCReporterOptions struct {
	Log *LogOptions
	CheckInterval time.Duration
	MaxSendQueueSize int
	InstanceProps map[string]string
	//Creds credentials.TransportCredentials
	Auth string
	//CDSInterval time.Duration
}

type TracerOptions struct {
	Instance     string
	SamplingRate float64
	KeyCount     int
	ValueSize    int
}

type config struct {
	ServiceName string
	ServerAddr string
	GRPCReporterOptions *GRPCReporterOptions
	TracerOptions *TracerOptions
	PanicOnError bool
}

func DefaultConfig() *config {
	return &config{
		ServiceName: "default",
		ServerAddr: "localhost:11800",
		GRPCReporterOptions: nil,
		TracerOptions: nil,
		PanicOnError: true,
	}
}

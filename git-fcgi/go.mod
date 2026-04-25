module go-queue

go 1.25.8

require (
	github.com/pingcap/tidb/v8 v8.5.6
	github.com/prometheus/procfs v0.19.2
	github.com/spf13/pflag v1.0.10
	k8s.io/klog/v2 v2.140.0
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/pingcap/errors v0.11.5-0.20250523034308-74f78ae071ee // indirect
	github.com/pingcap/failpoint v0.0.0-20240528011301-b51a646c7c86 // indirect
	github.com/pingcap/log v1.1.1-0.20250917021125-19901e015dc9 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/pingcap/tidb/v8 => github.com/pingcap/tidb v0.0.0-20260413061245-ae18096e0237

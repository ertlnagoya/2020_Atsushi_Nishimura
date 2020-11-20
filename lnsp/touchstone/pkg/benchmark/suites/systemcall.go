package suites

import (
	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/runtime"
	"github.com/lnsp/touchstone/pkg/util"
	"github.com/sirupsen/logrus"
)

var PerformanceSystemcall = []benchmark.Benchmark{
	&SystemcallScore{},
}

const defaultUnixbenchImage = "paipoi/unixbench"

// RunInSysbench executes a specific sysbench benchmark and returns the application logs.
func RunInUnixbench(bm benchmark.Benchmark, client *runtime.Client, handler string, args []string) ([]byte, error) {
	var (
		sandboxID   = benchmark.ID(bm)
		containerID = benchmark.ID(bm)
	)
	// Pull image
	if err := client.PullImage(defaultUnixbenchImage, nil); err != nil {
		return nil, err
	}
	// Perform benchmark
	sandbox := client.InitLinuxSandbox(sandboxID)
	pod, err := client.StartSandbox(sandbox, handler)
	if err != nil {
		return nil, err
	}
	container, err := client.CreateContainer(sandbox, pod, containerID, defaultUnixbenchImage, args)
	if err != nil {
		return nil, err
	}
	if err := client.StartContainer(container); err != nil {
		return nil, err
	}
	logs, err := client.WaitForLogs(container)
	if err != nil {
		return nil, err
	}
	// Cleanup container and sandbox
	if err := client.StopAndRemoveContainer(container); err != nil {
		return nil, err
	}
	if err := client.StopAndRemoveSandbox(pod); err != nil {
		return nil, err
	}
	logrus.WithField("name", bm.Name()).Debugf("unixbench logs: %v", string(logs))
	return logs, nil
}


// CPUTime measures the total time taken by a CPU heavy task.
type SystemcallScore struct{}

func (SystemcallScore) Name() string {
	return "systemcall"
}

func (bm *SystemcallScore) Run(client *runtime.Client, handler string) (benchmark.Report, error) {
	logs, err := RunInUnixbench(bm, client, handler, []string{
	})
	if err != nil {
		return nil, err
	}
	return benchmark.ValueReport{
		"System Benchmarks Index Score (Partial Only)": util.ParsePrefixedLine(logs, "System Benchmarks Index Score (Partial Only)"),
	}, nil
}

func (SystemcallScore) Labels() []string {
	return []string{"System Benchmarks Index Score (Partial Only)"}
}



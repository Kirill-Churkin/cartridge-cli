package bench

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tarantool/cartridge-cli/cli/context"
)

// Main benchmark function.
func Run(ctx context.BenchCtx) error {
	rand.Seed(time.Now().UnixNano())

	if err := verifyOperationsPercentage(&ctx); err != nil {
		return err
	}

	cluster, err := isCluster(ctx)
	if err != nil {
		return err
	}

	if cluster {
		if err := verifyClusterTopology(ctx); err != nil {
			return err
		}
		ctx.URL = (*ctx.Leaders)[0]
	}

	// Connect to tarantool and preset space for benchmark.
	tarantoolConnection, err := createConnection(ctx)
	if err != nil {
		return err
	}
	defer tarantoolConnection.Close()

	printConfig(ctx, tarantoolConnection)

	if err := spacePreset(ctx, tarantoolConnection); err != nil {
		return err
	}

	if err := preFillBenchmarkSpaceIfRequired(ctx); err != nil {
		return err
	}

	fmt.Println("Benchmark start")
	fmt.Println("...")

	// The "context" will be used to stop all "connectionLoop" when the time is out.

	benchStart := benchOneInstance
	if cluster {
		benchStart = benchCluster
	}

	benchData := getBenchData(ctx)

	if err := benchStart(ctx, &benchData); err != nil {
		return err
	}

	benchData.results.duration = time.Since(benchData.startTime).Seconds()
	benchData.results.requestsPerSecond = int(float64(benchData.results.handledRequestsCount) / benchData.results.duration)

	if err := dropBenchmarkSpace(tarantoolConnection); err != nil {
		return err
	}
	fmt.Println("Benchmark stop.")

	printResults(benchData.results)
	return nil
}

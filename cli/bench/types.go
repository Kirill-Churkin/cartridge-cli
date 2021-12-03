package bench

import (
	"github.com/FZambia/tarantool"
	"github.com/tarantool/cartridge-cli/cli/context"
)

// Results describes set of benchmark results.
type Results struct {
	handledRequestsCount int     // Count of all executed requests.
	successResultCount   int     // Count of successful request in all connections.
	failedResultCount    int     // Count of failed request in all connections.
	duration             float64 // Benchmark duration.
	requestsPerSecond    int     // Cumber of requests per second - the main measured value.
}

type Request interface {
	execute()
}

type InsertRequest struct {
	ctx                 context.BenchCtx
	tarantoolConnection *tarantool.Connection
	results             *Results
}

type SelectRequest struct {
	ctx                   context.BenchCtx
	tarantoolConnection   *tarantool.Connection
	results               *Results
	getRandomTupleCommand string
}

type UpdateRequest struct {
	ctx                   context.BenchCtx
	tarantoolConnection   *tarantool.Connection
	results               *Results
	getRandomTupleCommand string
}

type RequestsPool struct {
	request Request
	count   int
}

type RequestsSequence struct {
	requests            [requestTypesCount]RequestsPool
	currentRequestIndex int
	currentCounter      int
}

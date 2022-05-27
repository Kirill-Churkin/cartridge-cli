package bench

import (
	"math/rand"
	"reflect"

	"github.com/FZambia/tarantool"
	"github.com/tarantool/cartridge-cli/cli/common"
)

// insertOperationOnConnection execute insert operation with specified connection.
func insertOperationOnConnection(tarantoolConnection *tarantool.Connection, request *Request) {
	_, err := tarantoolConnection.Exec(
		tarantool.Insert(
			benchSpaceName,
			[]interface{}{
				common.RandomString(request.ctx.KeySize),
				common.RandomString(request.ctx.DataSize),
			}))
	request.results.incrementRequestsCounters(err)
}

func shardingInsertOperationOnConnection(tarantoolConnection *tarantool.Connection, request *Request) {
	key := common.RandomString(request.ctx.KeySize)
	value := common.RandomString(request.ctx.DataSize)
	bucket := calculateBucket(key, request.ctx)
	_, err := tarantoolConnection.Exec(tarantool.Call(vshardRWCommand, []interface{}{
		bucket,
		benchSpaceInsertCommand,
		[]interface{}{[]interface{}{
			key,
			bucket,
			value,
		}},
	}))
	request.results.incrementRequestsCounters(err)
}

// insertOperation execute insert operation.
func insertOperation(request *Request) {
	request.onConnectionOperation(request.tarantoolConnection, request)
}

// clusterInsertOperation execute insert operation on cluster topology.
func clusterInsertOperation(request *Request) {
	tarantoolConnection := request.clusterNodesConnections.getNextConnectionsPool().getNextConnection()
	request.onConnectionOperation(tarantoolConnection, request)
}

// selectOperationOnConnection execute select operation with specified connection.
func selectOperationOnConnection(tarantoolConnection *tarantool.Connection, request *Request) {
	_, err := tarantoolConnection.Exec(tarantool.Call(
		getRandomTupleCommand,
		[]interface{}{rand.Int()}))
	request.results.incrementRequestsCounters(err)
}

func shardingSelectOperationOnConnection(tarantoolConnection *tarantool.Connection, request *Request) {
	_, err := tarantoolConnection.Exec(tarantool.Call(
		getRandomTupleCommand,
		[]interface{}{rand.Int()}))
	request.results.incrementRequestsCounters(err)
}

// selectOperation execute select operation.
func selectOperation(request *Request) {
	request.onConnectionOperation(request.tarantoolConnection, request)
}

// clusterSelectOperation execute select operation on cluster topology.
func clusterSelectOperation(request *Request) {
	tarantoolConnection := request.clusterNodesConnections.getNextConnectionsPool().getNextConnection()
	request.onConnectionOperation(tarantoolConnection, request)
}

// updateOperationOnConnection execute update operation with specified connection.
func updateOperationOnConnection(tarantoolConnection *tarantool.Connection, request *Request) {
	getRandomTupleResponse, err := tarantoolConnection.Exec(
		tarantool.Call(getRandomTupleCommand,
			[]interface{}{rand.Int()}))
	if err == nil {
		data := getRandomTupleResponse.Data
		if len(data) > 0 {
			key := reflect.ValueOf(data[0]).Index(0).Elem().String()
			_, err := tarantoolConnection.Exec(
				tarantool.Update(
					benchSpaceName,
					benchSpacePrimaryIndexName,
					[]interface{}{key},
					[]tarantool.Op{tarantool.Op(
						tarantool.OpAssign(
							2,
							common.RandomString(request.ctx.DataSize)))}))
			request.results.incrementRequestsCounters(err)
		}
	}
}

func shardingUpdateOperationOnConnection(tarantoolConnection *tarantool.Connection, request *Request) {
	getRandomTupleResponse, err := tarantoolConnection.Exec(
		tarantool.Call(getRandomTupleCommand,
			[]interface{}{rand.Int()}))
	if err == nil {
		data := getRandomTupleResponse.Data
		if len(data) > 0 {
			key := reflect.ValueOf(data[0]).Index(0).Elem().String()
			bucket := calculateBucket(key, request.ctx)
			value := common.RandomString(request.ctx.DataSize)
			_, err := tarantoolConnection.Exec(tarantool.Call(vshardRWCommand, []interface{}{
				bucket,
				benchSpaceUpdateCommand,
				[]interface{}{
					benchSpaceName,
					benchSpacePrimaryIndexName,
					[]interface{}{key},
					[]tarantool.Op{tarantool.Op(
						tarantool.OpAssign(
							3,
							value))},
				}}))
			request.results.incrementRequestsCounters(err)
		}
	}
}

// updateOperation execute update operation.
func updateOperation(request *Request) {
	request.onConnectionOperation(request.tarantoolConnection, request)
}

// clusterUpdateOperation execute update operation on cluster topology.
func clusterUpdateOperation(request *Request) {
	tarantoolConnection := request.clusterNodesConnections.getNextConnectionsPool().getNextConnection()
	request.onConnectionOperation(tarantoolConnection, request)
}

// deleteOperationOnConnection execute delete operation with specified connection.
func deleteOperationOnConnection(tarantoolConnection *tarantool.Connection, request *Request) {
	getRandomTupleResponse, err := tarantoolConnection.Exec(
		tarantool.Call(getRandomTupleCommand,
			[]interface{}{rand.Int()}))
	if err == nil {
		data := getRandomTupleResponse.Data
		if len(data) > 0 {
			key := reflect.ValueOf(data[0]).Index(0).Elem().String()
			_, err := tarantoolConnection.Exec(
				tarantool.Delete(
					benchSpaceName,
					benchSpacePrimaryIndexName,
					[]interface{}{key},
				))
			request.results.incrementRequestsCounters(err)
		}
	}
}

func shardingDeleteOperationOnConnection(tarantoolConnection *tarantool.Connection, request *Request) {
	getRandomTupleResponse, err := tarantoolConnection.Exec(
		tarantool.Call(getRandomTupleCommand,
			[]interface{}{rand.Int()}))
	if err == nil {
		data := getRandomTupleResponse.Data
		if len(data) > 0 {
			key := reflect.ValueOf(data[0]).Index(0).Elem().String()
			bucket := calculateBucket(key, request.ctx)
			_, err := tarantoolConnection.Exec(tarantool.Call(vshardRWCommand, []interface{}{
				bucket,
				benchSpaceUpdateCommand,
				[]interface{}{
					benchSpaceName,
					benchSpacePrimaryIndexName,
					[]interface{}{key},
				}}))
			request.results.incrementRequestsCounters(err)
		}
	}
}

// deleteOperation execute delete operation.
func deleteOperation(request *Request) {
	request.onConnectionOperation(request.tarantoolConnection, request)
}

// clusterDeleteOperation execute delete operation on cluster topology.
func clusterDeleteOperation(request *Request) {
	tarantoolConnection := request.clusterNodesConnections.getNextConnectionsPool().getNextConnection()
	request.onConnectionOperation(tarantoolConnection, request)
}

// getNext return next operation in operations sequence.
func (requestsSequence *RequestsSequence) getNext() *Request {
	// If at the moment the number of remaining requests = 0,
	// then find a new generator, which requests count > 0.
	// If new generator has requests count = 0, then repeat.
	requestsSequence.findNewRequestsGeneratorMutex.Lock()
	defer requestsSequence.findNewRequestsGeneratorMutex.Unlock()
	for requestsSequence.currentCounter == 0 {
		// Increase the index, which means logical switching to a new generator.
		requestsSequence.currentRequestIndex++
		requestsSequence.currentRequestIndex %= len(requestsSequence.requests)
		// Get new generator by index.
		nextRequestsGenerator := &requestsSequence.requests[requestsSequence.currentRequestIndex]
		// Get requests count for new operation.
		requestsSequence.currentCounter = nextRequestsGenerator.count
	}
	// Logical taking of a single request.
	requestsSequence.currentCounter--
	return &requestsSequence.requests[requestsSequence.currentRequestIndex].request
}

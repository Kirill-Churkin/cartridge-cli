package bench

import (
	"math/rand"
	"reflect"

	"github.com/FZambia/tarantool"
	"github.com/tarantool/cartridge-cli/cli/common"
)

func (request InsertRequest) execute() {
	_, err := request.tarantoolConnection.Exec(
		tarantool.Insert(
			benchSpaceName,
			[]interface{}{common.RandomString(request.ctx.KeySize), common.RandomString(request.ctx.DataSize)}))
	incrementRequest(err, request.results)
}

func (request SelectRequest) execute() {
	_, err := request.tarantoolConnection.Exec(tarantool.Call(request.getRandomTupleCommand, []interface{}{rand.Int()}))
	incrementRequest(err, request.results)
}

func (request UpdateRequest) execute() {
	getRandomTupleResponse, err := request.tarantoolConnection.Exec(tarantool.Call(request.getRandomTupleCommand, []interface{}{rand.Int()}))
	if err == nil {
		data := getRandomTupleResponse.Data
		if len(data) > 0 {
			key := reflect.ValueOf(data[0]).Index(0).Elem().String()
			_, err := request.tarantoolConnection.Exec(
				tarantool.Update(
					benchSpaceName,
					benchSpacePrimaryIndexName,
					[]interface{}{key},
					[]tarantool.Op{tarantool.Op(tarantool.OpAssign(2, common.RandomString(request.ctx.DataSize)))}))
			incrementRequest(err, request.results)
		}
	}
}

func (requestsSequence RequestsSequence) getNext() Request {
	for requestsSequence.currentCounter == 0 {
		requestsSequence.currentRequestIndex++
		requestsSequence.currentRequestIndex %= requestTypesCount
		requestsSequence.currentCounter = requestsSequence.requests[requestsSequence.currentRequestIndex].count
	}
	requestsSequence.currentCounter--
	return requestsSequence.requests[requestsSequence.currentRequestIndex].request
}

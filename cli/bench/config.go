package bench

import "fmt"

var (
	benchSpaceName             = "__benchmark_space__"
	benchSpacePrimaryIndexName = "__bench_primary_key__"
	PreFillingCount            = 1000000
	vshardRWCommand            = "vshard.router.callrw"
	vshardRCommand             = "vshard.router.callro"
	getRandomTupleCommand      = fmt.Sprintf(
		"box.space.%s.index.%s:random",
		benchSpaceName,
		benchSpacePrimaryIndexName,
	)
	benchSpaceInsertCommand = fmt.Sprintf(
		"box.space.%s:insert",
		benchSpaceName,
	)
	benchSpaceUpdateCommand = fmt.Sprintf(
		"box.space.%s:update",
		benchSpaceName,
	)
	benchSpaceDeleteCommand = fmt.Sprintf(
		"box.space.%s:delete",
		benchSpaceName,
	)
)

package micro

type RunMode uint8

const (
	ByHost   RunMode = 1
	ByDocker RunMode = 2
	ByK8s    RunMode = 3
)

package micro

type HttpVerb uint8

const (
	GET    HttpVerb = 1
	POST   HttpVerb = 2
	PUT    HttpVerb = 3
	DELETE HttpVerb = 4
)

package c2s

type Rpc struct {
	Node          string
	ServiceMethod string
	Data          []byte
}

type Login struct {
	Name string
	Pass string
}

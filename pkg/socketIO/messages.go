package socketIO

const (
	CHELLO   = "2probe"
	SHELLO   = "3probe"
	ACKHELLO = "5"
	PING     = "2"
	PONG     = "3"
	MESSAGE  = "42"
)

type IncomingMessage struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type OutgoingMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

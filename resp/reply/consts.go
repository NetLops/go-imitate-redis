package reply

type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func (p *PongReply) ToBytes() []byte {
	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

type okReply struct {
}

var okBytes = []byte("+OK\r\n")

// ToBytes marshal redis.Reply
func (r *okReply) ToBytes() []byte {
	return okBytes
}

// 节约内存
var theOkReply = new(okReply)

// MakeOkReply returns a ok reply
func MakeOkReply() *okReply {
	return theOkReply
}

type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

var emptyMultiBulkBytes = []byte("*0\r\n")

type EmptyMultiBulkReply struct {
}

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

type Noreply struct {
}

var noBytes = []byte("")

func (n *Noreply) ToBytes() []byte {
	return noBytes
}

package reply

type UnknownErrReply struct {
}

var unknownErrBytes = []byte("-Err unknown\r\n")

func (u *UnknownErrReply) Error() string {
	return "-Err unknown"
}

func (u *UnknownErrReply) ToBytes() []byte {
	return unknownErrBytes
}

func MakeUnknownReply() *UnknownErrReply {
	return &UnknownErrReply{}
}

// 写给redis 参数有错误的时候
type ArgNumErrReply struct {
	Cmd string //  提示客户那个参数错误
}

func (a *ArgNumErrReply) Error() string {
	return "-ERR Wrong number of arguments for '" + a.Cmd + "' command"
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR Wrong number of arguments for '" + a.Cmd + "' command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

type SyntaxErrReply struct {
}

var syntaxErrBytes = []byte("-Err syntax error\r\n")
var theSyntaxErReply = &SyntaxErrReply{}

// MakeSyntaxErrReply creates syntax error
func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErReply
}

func (s SyntaxErrReply) Error() string {
	return "Err syntax error"
}

// ToBytes marshals redis.Reply
func (s SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

// 数据类型错误
type WrongTypeErrReply struct {
}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

func (r *WrongTypeErrReply) Error() string {
	return "-WRONGTYPE Operation against a key holding the wrong kind of value"
}

func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

// ProtocolErrReply 用户的协议不符合 RESP规范
// Example：
// $、-、+、*、:
type ProtocolErrReply struct {
	Msg string
}

func (p *ProtocolErrReply) Error() string {
	return "ERR Protocol error: '" + p.Msg
}

func (p *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + p.Msg + "'\r\n")
}

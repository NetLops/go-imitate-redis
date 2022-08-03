package resp

type Connection interface {
	Write([]byte) error
	GetDBIndex() int // 获取DB
	SelectDB(int)    // 切换DB
}

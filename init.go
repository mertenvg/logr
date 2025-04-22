package logr

func init() {
	pool.New = func() any {
		return &Message{}
	}
	SetBufferSize(10000)
}

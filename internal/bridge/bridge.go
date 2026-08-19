package bridge

type Bridge struct {
	ID       string
	UUID     string
	HTTPAddr string
}

func New(id, uuid, httpAddr string) *Bridge {
	return &Bridge{
		ID:       id,
		UUID:     uuid,
		HTTPAddr: httpAddr,
	}
}

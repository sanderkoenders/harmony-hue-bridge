package bridge

type Bridge struct {
	ID     string
	UUID   string
	IpAddr string
	Port   int
}

func New(id, uuid, ipAddr string, port int) *Bridge {
	return &Bridge{
		ID:     id,
		UUID:   uuid,
		IpAddr: ipAddr,
		Port:   port,
	}
}

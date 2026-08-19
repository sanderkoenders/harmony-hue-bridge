package bridge

type Bridge struct {
	Username string
	UUID     string
	IpAddr   string
	Port     int
}

func New(id, uuid, ipAddr string, port int) *Bridge {
	return &Bridge{
		Username: id,
		UUID:     uuid,
		IpAddr:   ipAddr,
		Port:     port,
	}
}

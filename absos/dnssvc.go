package absos

import "net"

type DnsSvc interface {
	LookupIP(host string) ([]net.IP, error)
}

type dnsSvcImpl struct{}

var dnsSvcImplInstance = dnsSvcImpl{}

func NewDnsSvc() DnsSvc {
	return dnsSvcImplInstance
}

func (dnsSvcImpl) LookupIP(host string) ([]net.IP, error) {
	return net.LookupIP(host)
}

package ngns

import (
	"net"
)

type ipPool struct {
	networks  []*net.IPNet
	netIndex  int
	currentIP net.IP
}

func newPool(networks []*net.IPNet) ipPool {
	currentIP := net.IP{}
	if len(networks) > 0 {
		network := networks[0]
		currentIP = network.IP.Mask(network.Mask)
	}

	return ipPool{
		networks:  networks,
		currentIP: currentIP,
	}
}

func (p *ipPool) next() (net.IP, bool) {
	if len(p.networks) == 0 {
		return net.IP{}, false
	}

	currentNetwork := p.networks[p.netIndex]
	if currentNetwork.Contains(p.currentIP) {
		defer incIP(p.currentIP)
		res := make(net.IP, net.IPv4len)
		copy(res, p.currentIP)
		return res, true
	}

	p.netIndex++
	if len(p.networks) <= p.netIndex {
		return net.IP{}, false
	}

	currentNetwork = p.networks[p.netIndex]
	p.currentIP = currentNetwork.IP.Mask(currentNetwork.Mask)
	return p.next()
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

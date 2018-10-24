package ngns

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Port represents network service port
type Port uint16

// Scanner is an instance of the dnrtlib scanner
type Scanner struct {
	pool    ipPool
	ports   []Port
	jobs    chan *scanJob
	results chan<- *ScanResult
	workers int
	wg      sync.WaitGroup
}

// ScannerConfig is a config used to create new instance of Scanner
type ScannerConfig struct {
	IPNetworks []*net.IPNet
	Ports      []Port
	Results    chan<- *ScanResult
	Workers    int
}

// ScanResult represents network scan result
type ScanResult struct {
	IP   net.IP
	Port Port
}

type scanJob struct {
	ip   net.IP
	port Port
}

// NewScanner creates new instance of the dnrtlib scanner
func NewScanner(config ScannerConfig) Scanner {
	return Scanner{
		pool:    newPool(config.IPNetworks),
		ports:   config.Ports,
		jobs:    make(chan *scanJob),
		results: config.Results,
		workers: config.Workers,
	}
}

// Scan starts the scan
func (s *Scanner) Scan() {
	for i := 0; i < s.workers; i++ {
		go s.scanWorker()
	}

	for {
		ip, more := s.pool.next()
		if !more {
			break
		}

		for _, port := range s.ports {
			s.wg.Add(1)
			s.jobs <- &scanJob{
				ip:   ip,
				port: port,
			}
		}
	}

	s.wg.Wait()
	close(s.results)
}

func (s *Scanner) scanWorker() {
	for j := range s.jobs {
		if s.checkPort(j.ip, j.port) {
			s.results <- &ScanResult{
				IP:   j.ip,
				Port: j.port,
			}
		}

		s.wg.Done()
	}
}

func (s *Scanner) checkPort(ip net.IP, port Port) (ok bool) {
	addressString := fmt.Sprintf("%s:%d", ip.String(), port)
	conn, err := net.DialTimeout("tcp", addressString, time.Second)
	if err == nil {
		defer conn.Close()
		return true
	}

	return false
}

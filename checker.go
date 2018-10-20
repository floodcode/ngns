package ngns

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// CheckerJob stores job configuration for Bruteforcer
type CheckerJob struct {
	IP   net.IP
	Port Port
	User string
	Pass string
}

// ServiceChecker checks single service credentials pair
type ServiceChecker interface {
	Check(CheckerJob) (ok bool)
}

// SSHChecker checks single SSH credentials pair
type SSHChecker struct {
}

// Check on success returns brute result
func (SSHChecker) Check(job CheckerJob) (ok bool) {
	sshConfig := &ssh.ClientConfig{
		User:            job.User,
		Auth:            []ssh.AuthMethod{ssh.Password(job.Pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second,
	}

	host := fmt.Sprintf("%s:%d", job.IP, 22)
	client, err := dial("tcp", host, sshConfig)
	if err == nil {
		client.Conn.Close()
		client.Close()
		return true
	}

	return false
}

// custom dial function with connection deadline
// used to prevent zombie connections
func dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	conn, err := net.DialTimeout(network, addr, config.Timeout)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Now().Add(time.Minute))

	c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(c, chans, reqs), nil
}

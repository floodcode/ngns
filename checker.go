package ngns

import (
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const defaultTimeout = time.Second * 4

// CheckerJob stores job configuration for Bruteforcer
type CheckerJob struct {
	IP   net.IP
	Port Port
	User string
	Pass string
}

// ServiceChecker checks single service credentials pair
type ServiceChecker interface {
	Check(CheckerJob) (summary string, ok bool)
}

// SSHChecker checks single SSH credentials pair
type SSHChecker struct {
}

// Check on success returns brute result
func (SSHChecker) Check(job CheckerJob) (summary string, ok bool) {
	sshConfig := &ssh.ClientConfig{
		User:            job.User,
		Auth:            []ssh.AuthMethod{ssh.Password(job.Pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         defaultTimeout,
	}

	host := fmt.Sprintf("%s:%d", job.IP, job.Port)
	client, conn, err := dial("tcp", host, sshConfig)
	if err != nil {
		return "", false
	}

	defer client.Close()
	defer conn.Close()

	distroMap := map[string]string{
		"fedora": "Fedora",
		"ubuntu": "Ubuntu",
		"debian": "Debian",
	}

	cmdOutput, _ := runCommand(client, conn, "cat /etc/os-release")
	if len(cmdOutput) > 0 {
		for id, distro := range distroMap {
			if strings.Contains(cmdOutput, "ID="+id) {
				return distro, true
			}
		}
	}

	cmdOutput, _ = runCommand(client, conn, "system resource print")
	if strings.Contains(cmdOutput, "MikroTik") {
		return "MikroTik", true
	}

	return "", true
}

// custom dial function with connection deadline
// used to prevent zombie connections
func dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, net.Conn, error) {
	conn, err := net.DialTimeout(network, addr, defaultTimeout)
	if err != nil {
		return nil, nil, err
	}

	conn.SetDeadline(time.Now().Add(defaultTimeout))

	c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, nil, err
	}

	return ssh.NewClient(c, chans, reqs), conn, nil
}

func runCommand(client *ssh.Client, conn net.Conn, command string) (string, error) {
	conn.SetDeadline(time.Now().Add(defaultTimeout))
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	out, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

package ngns

import (
	"net"
)

// BruteResult represents valid credentials of the bruteforced service
type BruteResult struct {
	User string
	Pass string
}

// BruteConfig stores configuration for Bruteforcer
type BruteConfig struct {
	ServiceChecker ServiceChecker
	IP             net.IP
	Port           Port
	User           StringSource
	Pass           StringSource
}

// Bruteforcer used to bruteforce TCP services
type Bruteforcer struct {
}

// Brute on success returns brute result
func (b Bruteforcer) Brute(conf BruteConfig) (res BruteResult, ok bool) {
	for {
		user, ok := conf.User.Next()
		if !ok {
			break
		}

		for {
			pass, ok := conf.Pass.Next()
			if !ok {
				break
			}

			checked := conf.ServiceChecker.Check(CheckerJob{
				IP:   conf.IP,
				Port: conf.Port,
				User: user,
				Pass: pass,
			})

			if checked {
				return BruteResult{
					User: user,
					Pass: pass,
				}, true
			}
		}

		conf.Pass.Reset()
	}

	return BruteResult{}, false
}

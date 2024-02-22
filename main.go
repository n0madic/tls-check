package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"time"
)

type TLSVersion struct {
	Name       string
	Version    uint16
	Deprecated bool
}

type CheckResult struct {
	Supported    bool
	Expired      bool
	ExpiringSoon bool
	ExpiryDate   time.Time
	Error        error
}

func checkTLSVersion(server string, tlsVersion TLSVersion) CheckResult {
	config := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tlsVersion.Version,
		MaxVersion:         tlsVersion.Version,
	}

	conn, err := tls.Dial("tcp", server, config)
	if err != nil {
		return CheckResult{Supported: false, Error: err}
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	now := time.Now()
	expired := cert.NotAfter.Before(now)
	expiringSoon := now.Add(30 * 24 * time.Hour).After(cert.NotAfter)

	return CheckResult{
		Supported:    true,
		Expired:      expired,
		ExpiringSoon: expiringSoon,
		ExpiryDate:   cert.NotAfter,
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s server1.com[:port] [server2.com[:port] ...]\n", os.Args[0])
		os.Exit(1)
	}

	tlsVersions := []TLSVersion{
		{"TLS v1.0", tls.VersionTLS10, true},
		{"TLS v1.1", tls.VersionTLS11, true},
		{"TLS v1.2", tls.VersionTLS12, false},
		{"TLS v1.3", tls.VersionTLS13, false},
	}

	for _, arg := range os.Args[1:] {
		server := arg
		if _, _, err := net.SplitHostPort(arg); err != nil {
			server = arg + ":443"
		}

		fmt.Printf("Checking supported TLS versions for %s\n", server)

		for _, tlsVersion := range tlsVersions {
			fmt.Printf("Testing %s...\t", tlsVersion.Name)
			result := checkTLSVersion(server, tlsVersion)

			resultColor := "\033[32m" // Green for supported
			if tlsVersion.Deprecated || result.Error != nil {
				resultColor = "\033[31m" // Red for deprecated
			}
			if result.Supported {
				fmt.Printf("%sSupported\033[0m", resultColor)
			} else {
				fmt.Print("Not supported")
				if result.Error != nil {
					fmt.Printf(" (%s)", result.Error)
				}
			}
			if result.Expired {
				fmt.Printf(" [Certificate expired on %s]", result.ExpiryDate.Format("2006-01-02"))
			} else if result.ExpiringSoon {
				fmt.Printf(" [Certificate expiring soon on %s]", result.ExpiryDate.Format("2006-01-02"))
			}
			fmt.Println() // New line for readability
		}
	}
}

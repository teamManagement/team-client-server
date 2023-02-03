//go:build none_Safe

package config

var (
	InsecureSkipVerify = true
	HttpsTlsConfig     = &tls.Config{
		InsecureSkipVerify: InsecureSkipVerify,
	}
)

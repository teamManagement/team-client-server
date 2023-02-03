//go:build safe

package config

import "crypto/tls"

var (
	InsecureSkipVerify = false
	HttpsTlsConfig     = &tls.Config{
		InsecureSkipVerify: InsecureSkipVerify,
	}
)

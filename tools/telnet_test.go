package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTelnetHostRangeNetInterfaces(t *testing.T) {
	a := assert.New(t)
	addr, ok := TelnetHostRangeNetInterfaces("apps.byzk.cn:80")
	//addr, ok := TelnetHostRangeNetInterfaces("127.0.0.1:80")
	//addr, ok := TelnetHostRangeNetInterfaces("baidu.com:80")
	a.True(ok)
	a.NotNil(addr)
}

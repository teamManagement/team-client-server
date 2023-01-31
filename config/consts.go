package config

var TeamworkSm4Key = []byte("teamwork queue message transfer fixed key!!!")[:16]

var (
	localWebServerAddress = "https://" + serverDomain
	localWSServerAddress  = "wss://" + serverDomain
)

func LocalWebServerAddress() string {
	return localWebServerAddress
}

func LocalWSServerAddress() string {
	return localWSServerAddress
}

func IsDebug() bool {
	return debug == "1"
}

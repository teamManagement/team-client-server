package config

var serverDomain string

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

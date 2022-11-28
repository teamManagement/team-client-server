package remoteserver

import "team-client-server/vos"

type InsideEventHandlerFn = func(userInfo *vos.UserInfo) error

type InsideEventName string

const (
	InsideEventNameFlushAppList InsideEventName = "app-desktop-list-flush"
)

var (
	eventMaps = make(map[InsideEventName]InsideEventHandlerFn)
)

func RegistryInsideEvent(eventName InsideEventName, fn InsideEventHandlerFn) {
	eventMaps[eventName] = fn
}

func InvokeInsideEvent(eventName InsideEventName, userInfo *vos.UserInfo) error {
	fn, ok := eventMaps[eventName]
	if !ok {
		return nil
	}
	return fn(userInfo)
}

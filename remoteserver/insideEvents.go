package remoteserver

type InsideEventHandlerFn = func(userInfo *UserInfo) error

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

func InvokeInsideEvent(eventName InsideEventName, userInfo *UserInfo) error {
	fn, ok := eventMaps[eventName]
	if !ok {
		return nil
	}
	return fn(userInfo)
}

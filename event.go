package mggo

type eventHadler func(interface{})

var events map[string][]eventHadler

// EventType is event type publish
type eventType int

// EventTypeClient is publish event only client
// EventTypeServer is publish event only server
// EventTypeGlobal is publish event server and client
const (
    EventTypeClient eventType = iota + 1
    EventTypeServer
    EventTypeGlobal
)

func init() {
    events = map[string][]eventHadler{}
}

// EventSubscribe is subscribe event by event name
func EventSubscribe(eventName string, handler eventHadler) {
    if v, ok := events[eventName]; ok {
        v = append(v, handler)
    } else {
        events[eventName] = []eventHadler{handler}
    }
}

// EventPublish is publick
func EventPublish(eventName string, et eventType, users []int, params ...interface{}) {
    // send server
    if v, ok := events[eventName]; ok {
        for _, handler := range v {
            if et > EventTypeClient {
                handler(params)
            }
        }
    }
    // send server
    if et != EventTypeServer {
        sendSockets(eventName, users, params)
    }
}

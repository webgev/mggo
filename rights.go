package mggo

// Rights
const (
	RRightGuest   = 0
	RRightUser    = 4
	RRightEditor  = 8
	RRightManager = 16
	RRightAdmin   = 32
)

var Rights = map[string]int{}
var rights = map[string]int{}
var rightsView = map[string]int{}

func init() {
	Rights["Guest"] = RRightGuest
	Rights["User"] = RRightUser
	Rights["Editor"] = RRightEditor
	Rights["Manager"] = RRightManager
	Rights["Admin"] = RRightAdmin
}

// AppendRight registration rigth
func AppendRight(method string, right int) {
	rights[method] = right
}

// AppendViewRight registration rigth by view
func AppendViewRight(method string, right int) {
	rightsView[method] = right
}

// GetRightMethod get right method by method name
func GetRightMethod(method string) (int, bool) {
	val, ok := rights[method]
	return val, ok
}

// GetRightView get right view by controller name and action
func GetRightView(controller string, action string) (int, bool) {
	val, ok := rightsView[controller+"."+action]
	if !ok {
		val, ok = rightsView[controller]
	}
	return val, ok
}

// CheckRight - check right in method
func CheckRight(method string, right int, hard bool) bool {
	val, ok := GetRightMethod(method)
	if hard && !ok {
		return false
	}
	return right >= val
}

// CheckViewRight - check right in view
func CheckViewRight(controller string, action string, right int, hard bool) bool {
	val, ok := GetRightView(controller, action)
	if hard && !ok {
		return false
	}
	return right >= val
}

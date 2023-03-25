package stats

type AlertFilterConstructor func(attrs map[string]string) AlertFilter

var alertQueryRegistry = make(map[string]AlertFilterConstructor)

func RegisterAlertFilter(name string, q AlertFilterConstructor) {
	alertQueryRegistry[name] = q
}

func LookupAlertFilter(name string, attrs map[string]string) AlertFilter {
	return alertQueryRegistry[name](attrs)
}

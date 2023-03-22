package stats

type AlertQueryConstructor func(attrs map[string]string) AlertQuery

var alertQueryRegistry = make(map[string]AlertQueryConstructor)

func RegisterAlertQuery(name string, q AlertQueryConstructor) {
	alertQueryRegistry[name] = q
}

func LookupAlertQuery(name string, attrs map[string]string) AlertQuery {
	return alertQueryRegistry[name](attrs)
}

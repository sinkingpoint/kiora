package alerts

type AlertsCmd struct {
	Get  AlertsGetCmd  `cmd:"" help:"Get alerts."`
	Post AlertsPostCmd `cmd:"" help:"Post alerts."`
	Test AlertsTestCmd `cmd:"" help:"Test alerts."`
}

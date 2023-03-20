package alerts

type AlertsCmd struct {
	Get  AlertsGetCmd  `cmd:"" help:"Get alerts."`
	Post AlertsPostCmd `cmd:"" help:"Post alerts."`
}

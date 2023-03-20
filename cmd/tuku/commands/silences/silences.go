package silences

type SilencesCmd struct {
	Post SilencePostCmd `cmd:"" help:"Add a silence."`
}

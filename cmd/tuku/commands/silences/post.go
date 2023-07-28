package silences

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/cmd/tuku/commands"
	"github.com/sinkingpoint/kiora/internal/stubs"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
)

type SilencePostCmd struct {
	Comment   string   `help:"The comment describing the silence."`
	Author    string   `help:"The author of the silence."`
	StartTime string   `help:"When the silence should start."`
	Duration  string   `help:"How long after the StartTime the silence should last for." short:"d" required:""`
	Matchers  []string `arg:"" help:"The matchers for the silence." required:""`
}

func (a *SilencePostCmd) Run(ctx *commands.Context) error {
	var startTime time.Time
	var err error

	if a.StartTime != "" {
		startTime, err = time.Parse(time.RFC3339, a.StartTime)
		if err != nil {
			return errors.Wrap(err, "failed to parse start time")
		}
	} else {
		startTime = stubs.Time.Now()
	}

	duration, err := time.ParseDuration(a.Duration)
	if err != nil {
		return errors.Wrap(err, "failed to parse duration")
	}

	endTime := startTime.Add(duration)

	matchers, err := parseMatchers(a.Matchers)
	if err != nil {
		return err
	}

	silence, err := ctx.Kiora.PostSilence(model.Silence{
		Comment:   a.Comment,
		Creator:   a.Author,
		StartTime: startTime,
		EndTime:   endTime,
		Matchers:  matchers,
	})
	if err != nil {
		return err
	}

	out, err := ctx.Formatter.Marshal(silence)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}

func parseMatchers(rawMatchers []string) ([]model.Matcher, error) {
	matchers := []model.Matcher{}
	for _, m := range rawMatchers {
		matcher := model.Matcher{}
		if err := matcher.UnmarshalText(m); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse matcher %q", m))
		}

		matchers = append(matchers, matcher)
	}

	return matchers, nil
}

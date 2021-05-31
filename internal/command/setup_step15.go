package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
)

type Step15 struct {
	DefaultMessageTexts []domain.CustomMessageText
}

func (s *Step15) Step() domain.Step {
	return domain.Step15
}

func (s *Step15) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep15(ctx, s)
}

func (c *Commands) SetupStep15(ctx context.Context, step *Step15) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		events := make([]eventstore.EventPusher, 0)

		for _, text := range step.DefaultMessageTexts {
			mailEvents, _, err := c.setDefaultMessageText(ctx, iamAgg, &text)
			if err != nil {
				return nil, err
			}
			events = append(events, mailEvents...)
		}

		logging.Log("SETUP-4insR").Info("default mail template/text set up")
		return events, nil
	}
	return c.setup(ctx, step, fn)
}

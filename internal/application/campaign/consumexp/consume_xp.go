package consumexp

import (
	"context"
	applicationcampaign "meye-core/internal/application/campaign"
	domaincampaign "meye-core/internal/domain/campaign"
	"meye-core/internal/domain/event"
)

type UseCase struct {
	pjRepository    domaincampaign.PjRepository
	eventsPublisher event.Publisher
}

func New(pjRepo domaincampaign.PjRepository, eventPub event.Publisher) *UseCase {
	return &UseCase{
		pjRepository:    pjRepo,
		eventsPublisher: eventPub,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input applicationcampaign.ConsumeXpInput) error {
	pj, err := uc.pjRepository.FindByID(ctx, input.PjID)
	if err != nil {
		return err
	}

	if pj == nil {
		return domaincampaign.ErrPjNotFound
	}

	pj.ConsumeXp(input.Xp.Basic, input.Xp.Special, input.Xp.Supernatural)

	if err = uc.pjRepository.Save(ctx, pj); err != nil {
		return err
	}

	return uc.eventsPublisher.Publish(ctx, pj.UncommittedEvents())
}

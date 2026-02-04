package campaign_test

import (
	"meye-core/internal/domain/campaign"
	"meye-core/tests/mocks"
	"meye-core/tests/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCampaign_AddPJ(t *testing.T) {
	type want struct {
		pj  *campaign.PJ
		err error
	}

	var idServiceMock *mocks.MockIdentificationService

	defaultParams := campaign.PJCreateParameters{
		Name:                     "test",
		Weight:                   70,
		Height:                   175,
		Age:                      25,
		Look:                     5,
		Charisma:                 5,
		Villainy:                 5,
		Heroism:                  5,
		PjType:                   campaign.PJTypeSupernatural,
		IsPhysicalTalented:       true,
		IsMentalTalented:         false,
		IsCoordinationTalented:   false,
		IsPhysicalSkillsTalented: true,
		IsMentalSkillsTalented:   false,
		IsEnergySkillsTalented:   false,
	}

	tests := []struct {
		name       string
		campaign   *campaign.Campaign
		userID     string
		params     campaign.PJCreateParameters
		want       want
		setupMocks func()
	}{
		{
			name:     "PJ is added correctly when user is invited",
			campaign: testdata.Campaign(t),
			userID:   testdata.UserID,
			params:   defaultParams,
			want: want{
				pj: campaign.CreatePJWithoutValidation(
					testdata.PJID,
					testdata.UserID,
					defaultParams.Name,
					defaultParams.Weight,
					defaultParams.Height,
					defaultParams.Age,
					defaultParams.Look,
					defaultParams.Charisma,
					defaultParams.Villainy,
					defaultParams.Heroism,
					defaultParams.PjType,
					campaign.CreateBasicStatsWithoutValidation(
						campaign.CreatePhysicalWithoutValidation(0, 0, 0, 0, true),
						campaign.CreateMentalWithoutValidation(0, 0, 0, 0, false),
						campaign.CreateCoordinationWithoutValidation(0, 0, 0, 0, false),
						0,
					),
					campaign.CreateSpecialStatsWithoutValidation(
						campaign.CreatePhysicalSkillsWithoutValidation(0, 0, true),
						campaign.CreateMentalSkillsWithoutValidation(0, 0, false),
						campaign.CreateEnergySkillsWithoutValidation(0, 0, false),
						0,
						false,
					),
					campaign.CreateSupernaturalStatsWithoutValidation([]campaign.Skill{
						campaign.CreateSkillWithoutValidation([]uint{0}),
					}),
				),
				err: nil,
			},
			setupMocks: func() {
				idServiceMock.EXPECT().
					GenerateID().
					Return(testdata.PJID).
					Times(1)
			},
		},
		{
			name:     "PJ is not added when user is not invited",
			campaign: testdata.Campaign(t),
			userID:   "not_invited",
			params:   defaultParams,
			want: want{
				pj:  nil,
				err: campaign.ErrUserNotInvited,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			idServiceMock = mocks.NewMockIdentificationService(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			got, err := tt.campaign.AddPJ(tt.userID, tt.params, idServiceMock)

			assert.Equal(t, tt.want.pj, got)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

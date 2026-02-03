package campaign_test

import (
	"meye-core/internal/domain/campaign"
	"meye-core/tests/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicStats_GetRequiredXP(t *testing.T) {
	tests := []struct {
		name       string
		talent     campaign.BasicTalentType
		basicStats campaign.BasicStats
		want       int
	}{
		{
			name:       "Works correclty with physical talent",
			talent:     campaign.BasicTalentPhysical,
			basicStats: testdata.BasicStats(),
			want:       1176,
		},
		{
			name:       "Works correclty with mental talent",
			talent:     campaign.BasicTalentMental,
			basicStats: testdata.BasicStats(),
			want:       1096,
		},
		{
			name:       "Works correclty with coordination talent",
			talent:     campaign.BasicTalentCoordination,
			basicStats: testdata.BasicStats(),
			want:       1016,
		},
		{
			name:       "Works correclty with energy talent",
			talent:     campaign.BasicTalentEnergy,
			basicStats: testdata.BasicStats(),
			want:       1280,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.basicStats.GetRequiredXP(test.talent)

			assert.Equal(t, test.want, got)
		})
	}
}

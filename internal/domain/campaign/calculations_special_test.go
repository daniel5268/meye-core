package campaign_test

import (
	"meye-core/internal/domain/campaign"
	"meye-core/tests/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecialStats_GetRequiredXP(t *testing.T) {
	tests := []struct {
		name         string
		specialStats campaign.SpecialStats
		want         int
	}{
		{
			name:         "Works correctly with physical skills talented",
			specialStats: testdata.SpecialStatsWithPhysicalTalent(),
			want:         4380,
		},
		{
			name:         "Works correctly with energy skills talented",
			specialStats: testdata.SpecialStatsWithEnergyTalent(),
			want:         4210,
		},
		{
			name:         "Works correctly with mental skills talented",
			specialStats: testdata.SpecialStatsWithMentalTalent(),
			want:         4350,
		},
		{
			name:         "Works correctly with energy tank talented (cheaper energy tank)",
			specialStats: testdata.SpecialStatsWithEnergyTankTalent(),
			want:         4200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.specialStats.GetRequiredXP()

			assert.Equal(t, test.want, got)
		})
	}
}

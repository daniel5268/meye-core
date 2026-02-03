package campaign_test

import (
	"meye-core/internal/domain/campaign"
	"meye-core/tests/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecialStats_GetRequiredXP(t *testing.T) {
	tests := []struct {
		name          string
		basicTalent   campaign.BasicTalentType
		specialTalent campaign.SpecialTalentType
		specialStats  campaign.SpecialStats
		want          int
	}{
		{
			name:          "Works correctly for special talent physical",
			basicTalent:   campaign.BasicTalentPhysical,
			specialTalent: campaign.SpecialTalentPhysical,
			specialStats:  testdata.SpecialStats(),
			want:          4380,
		},
		{
			name:          "Works correctly for special talent energy",
			basicTalent:   campaign.BasicTalentPhysical,
			specialTalent: campaign.SpecialTalentEnergy,
			specialStats:  testdata.SpecialStats(),
			want:          4210,
		},
		{
			name:          "Works correctly for special talent mental",
			basicTalent:   campaign.BasicTalentPhysical,
			specialTalent: campaign.SpecialTalentMental,
			specialStats:  testdata.SpecialStats(),
			want:          4350,
		},
		{
			name:          "Works correctly for basic talent energy",
			basicTalent:   campaign.BasicTalentEnergy,
			specialTalent: campaign.SpecialTalentMental,
			specialStats:  testdata.SpecialStats(),
			want:          4200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.specialStats.GetRequiredXP(test.basicTalent, test.specialTalent)

			assert.Equal(t, test.want, got)
		})
	}
}

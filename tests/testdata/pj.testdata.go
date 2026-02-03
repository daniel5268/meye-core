package testdata

import "meye-core/internal/domain/campaign"

func BasicStats() campaign.BasicStats {
	return campaign.CreateBasicStatsWithoutValidation(
		campaign.CreatePhysicalWithoutValidation(13, 13, 13, 13),
		campaign.CreateMentalWithoutValidation(23, 23, 23, 23),
		campaign.CreateCoordinationWithoutValidation(33, 33, 33, 33),
		44,
	)
}

func SpecialStats() campaign.SpecialStats {
	return campaign.CreateSpecialStatsWithoutValidation(
		campaign.CreatePhysicalSkillsWithoutValidation(115, 225),
		campaign.CreateMentalSkillsWithoutValidation(215, 155),
		campaign.CreateEnergySkillsWithoutValidation(300, 210),
		30,
	)
}

const PJID = "test-pj-id"

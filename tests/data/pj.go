package data

import "meye-core/internal/domain/campaign"

func BasicStatsWithPhysicalTalent() campaign.BasicStats {
	return campaign.CreateBasicStatsWithoutValidation(
		campaign.CreatePhysicalWithoutValidation(13, 13, 13, 13, true),
		campaign.CreateMentalWithoutValidation(23, 23, 23, 23, false),
		campaign.CreateCoordinationWithoutValidation(33, 33, 33, 33, false),
		44,
	)
}

func BasicStatsWithMentalTalent() campaign.BasicStats {
	return campaign.CreateBasicStatsWithoutValidation(
		campaign.CreatePhysicalWithoutValidation(13, 13, 13, 13, false),
		campaign.CreateMentalWithoutValidation(23, 23, 23, 23, true),
		campaign.CreateCoordinationWithoutValidation(33, 33, 33, 33, false),
		44,
	)
}

func BasicStatsWithCoordinationTalent() campaign.BasicStats {
	return campaign.CreateBasicStatsWithoutValidation(
		campaign.CreatePhysicalWithoutValidation(13, 13, 13, 13, false),
		campaign.CreateMentalWithoutValidation(23, 23, 23, 23, false),
		campaign.CreateCoordinationWithoutValidation(33, 33, 33, 33, true),
		44,
	)
}

func BasicStatsWithNoTalents() campaign.BasicStats {
	return campaign.CreateBasicStatsWithoutValidation(
		campaign.CreatePhysicalWithoutValidation(13, 13, 13, 13, false),
		campaign.CreateMentalWithoutValidation(23, 23, 23, 23, false),
		campaign.CreateCoordinationWithoutValidation(33, 33, 33, 33, false),
		44,
	)
}

func SpecialStatsWithPhysicalTalent() campaign.SpecialStats {
	return campaign.CreateSpecialStatsWithoutValidation(
		campaign.CreatePhysicalSkillsWithoutValidation(115, 225, true),
		campaign.CreateMentalSkillsWithoutValidation(215, 155, false),
		campaign.CreateEnergySkillsWithoutValidation(300, 210, false),
		30,
		false, // isEnergyTalented
	)
}

func SpecialStatsWithMentalTalent() campaign.SpecialStats {
	return campaign.CreateSpecialStatsWithoutValidation(
		campaign.CreatePhysicalSkillsWithoutValidation(115, 225, false),
		campaign.CreateMentalSkillsWithoutValidation(215, 155, true),
		campaign.CreateEnergySkillsWithoutValidation(300, 210, false),
		30,
		false, // isEnergyTalented
	)
}

func SpecialStatsWithEnergyTalent() campaign.SpecialStats {
	return campaign.CreateSpecialStatsWithoutValidation(
		campaign.CreatePhysicalSkillsWithoutValidation(115, 225, false),
		campaign.CreateMentalSkillsWithoutValidation(215, 155, false),
		campaign.CreateEnergySkillsWithoutValidation(300, 210, true),
		30,
		false, // isEnergyTalented
	)
}

func SpecialStatsWithEnergyTankTalent() campaign.SpecialStats {
	return campaign.CreateSpecialStatsWithoutValidation(
		campaign.CreatePhysicalSkillsWithoutValidation(115, 225, false),
		campaign.CreateMentalSkillsWithoutValidation(215, 155, true),
		campaign.CreateEnergySkillsWithoutValidation(300, 210, false),
		30,
		true, // isEnergyTalented - for cheaper energy tank
	)
}

const PjID = "test-pj-id"

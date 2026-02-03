package campaign

const (
	levelStepSpecial   = 100
	energyDefaultCost  = 10
	energyTalentedCost = energyDefaultCost / 2
)

func (ss *SpecialStats) GetRequiredXP(basicTalent BasicTalentType, specialTalent SpecialTalentType) int {
	physicalGroup := ss.physical.getGroup()
	energyGroup := ss.energy.getGroup()
	mentalGroup := ss.mental.getGroup()

	physicalXP := getGroupRequiredXP(physicalGroup, levelStepSpecial, getSpecialFirstLevelCost(specialTalent, SpecialTalentPhysical))
	energyXP := getGroupRequiredXP(energyGroup, levelStepSpecial, getSpecialFirstLevelCost(specialTalent, SpecialTalentEnergy))
	mentalXP := getGroupRequiredXP(mentalGroup, levelStepSpecial, getSpecialFirstLevelCost(specialTalent, SpecialTalentMental))

	energyCost := getEnergyTankCost(basicTalent)

	energyTankXP := energyCost * ss.energyTank

	return int(physicalXP + energyXP + mentalXP + energyTankXP)
}

func (ps PhysicalSkills) getGroup() []uint {
	return []uint{ps.empowerment + ps.vitalControl}
}

func (ms MentalSkills) getGroup() []uint {
	return []uint{ms.mentalControl + ms.ilusion}
}

func (es EnergySkills) getGroup() []uint {
	return []uint{es.energyHandling + es.objectHandling}
}

func getSpecialFirstLevelCost(talent, target SpecialTalentType) uint {
	var firtsLevelCost uint = 1

	if talent != target {
		firtsLevelCost = 2
	}

	return firtsLevelCost
}

func getEnergyTankCost(basicTalent BasicTalentType) uint {
	if basicTalent == BasicTalentEnergy {
		return energyTalentedCost
	}

	return energyDefaultCost
}

func (ss SpecialStats) getOrderedStats() []uint {
	return []uint{
		ss.physical.empowerment,
		ss.physical.vitalControl,
		ss.energy.energyHandling,
		ss.energy.objectHandling,
		ss.mental.ilusion,
		ss.mental.mentalControl,
		ss.energyTank,
	}
}

func (ss SpecialStats) isHigherThan(ssB SpecialStats) bool {
	stats := ss.getOrderedStats()
	statsB := ssB.getOrderedStats()

	return isHigherThan(stats, statsB)
}

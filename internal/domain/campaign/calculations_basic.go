package campaign

const (
	levelStepBasic = 10
	costLife       = 5
)

func (bs *BasicStats) GetRequiredXP(basicTalent BasicTalentType) int {
	physicalGroup := bs.physical.getGroup()
	mentalGroup := bs.mental.getGroup()
	coordGroup := bs.coordination.getGroup()

	physicalXP := getGroupRequiredXP(physicalGroup, levelStepBasic, getBasicFirstLevelCost(basicTalent, BasicTalentPhysical))
	mentalXP := getGroupRequiredXP(mentalGroup, levelStepBasic, getBasicFirstLevelCost(basicTalent, BasicTalentMental))
	coordXP := getGroupRequiredXP(coordGroup, levelStepBasic, getBasicFirstLevelCost(basicTalent, BasicTalentCoordination))

	lifeXP := costLife * bs.life

	return int(physicalXP + mentalXP + coordXP + lifeXP)
}

func (p Physical) getGroup() []uint {
	return []uint{p.strength, p.agility, p.speed, p.resistance}
}

func (m Mental) getGroup() []uint {
	return []uint{m.inteligence, m.wisdom, m.concentration, m.will}
}

func (c Coordination) getGroup() []uint {
	return []uint{c.precision, c.calculation, c.coordRange, c.reflexes}
}

func getBasicFirstLevelCost(talent, target BasicTalentType) uint {
	var firtsLevelCost uint = 1

	if talent != target {
		firtsLevelCost = 3
	}

	return firtsLevelCost
}

func (bs BasicStats) isHigherThan(bsB BasicStats) bool {
	stats := bs.getOrderedStats()
	statsB := bsB.getOrderedStats()

	return isHigherThan(stats, statsB)
}

func (bs BasicStats) getOrderedStats() []uint {
	groups := [][]uint{
		bs.physical.getGroup(),
		bs.mental.getGroup(),
		bs.coordination.getGroup(),
		{bs.life},
	}

	return concatenateGroups(groups)
}

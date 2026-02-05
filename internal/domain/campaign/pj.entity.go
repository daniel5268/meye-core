package campaign

type PJType string

const (
	PJTypeHuman        PJType = "human"
	PJTypeSupernatural PJType = "supernatural"
)

type Physical struct {
	strength   uint
	agility    uint
	speed      uint
	resistance uint
	isTalented bool
}

func (p Physical) Strength() uint   { return p.strength }
func (p Physical) Agility() uint    { return p.agility }
func (p Physical) Speed() uint      { return p.speed }
func (p Physical) Resistance() uint { return p.resistance }
func (p Physical) IsTalented() bool { return p.isTalented }

func CreatePhysicalWithoutValidation(strength, agility, speed, resistance uint, isTalented bool) Physical {
	return Physical{
		strength:   strength,
		agility:    agility,
		speed:      speed,
		resistance: resistance,
		isTalented: isTalented,
	}
}

type Mental struct {
	inteligence   uint
	wisdom        uint
	concentration uint
	will          uint
	isTalented    bool
}

func (m Mental) Inteligence() uint   { return m.inteligence }
func (m Mental) Wisdom() uint        { return m.wisdom }
func (m Mental) Concentration() uint { return m.concentration }
func (m Mental) Will() uint          { return m.will }
func (m Mental) IsTalented() bool    { return m.isTalented }

func CreateMentalWithoutValidation(inteligence, wisdom, concentration, will uint, isTalented bool) Mental {
	return Mental{
		inteligence:   inteligence,
		wisdom:        wisdom,
		concentration: concentration,
		will:          will,
		isTalented:    isTalented,
	}
}

type Coordination struct {
	precision   uint
	calculation uint
	coordRange  uint
	reflexes    uint
	isTalented  bool
}

func (c Coordination) Precision() uint   { return c.precision }
func (c Coordination) Calculation() uint { return c.calculation }
func (c Coordination) Range() uint       { return c.coordRange }
func (c Coordination) Reflexes() uint    { return c.reflexes }
func (c Coordination) IsTalented() bool  { return c.isTalented }

func CreateCoordinationWithoutValidation(precision, calculation, coordRange, reflexes uint, isTalented bool) Coordination {
	return Coordination{
		precision:   precision,
		calculation: calculation,
		coordRange:  coordRange,
		reflexes:    reflexes,
		isTalented:  isTalented,
	}
}

type BasicStats struct {
	physical     Physical
	mental       Mental
	coordination Coordination
	life         uint
}

func (bs BasicStats) Physical() Physical         { return bs.physical }
func (bs BasicStats) Mental() Mental             { return bs.mental }
func (bs BasicStats) Coordination() Coordination { return bs.coordination }
func (bs BasicStats) Life() uint                 { return bs.life }

func CreateBasicStatsWithoutValidation(physical Physical, mental Mental, coordination Coordination, life uint) BasicStats {
	return BasicStats{
		physical:     physical,
		mental:       mental,
		coordination: coordination,
		life:         life,
	}
}

type PhysicalSkills struct {
	empowerment  uint
	vitalControl uint
	isTalented   bool
}

func (ps PhysicalSkills) Empowerment() uint  { return ps.empowerment }
func (ps PhysicalSkills) VitalControl() uint { return ps.vitalControl }
func (ps PhysicalSkills) IsTalented() bool   { return ps.isTalented }

func CreatePhysicalSkillsWithoutValidation(empowerment, vitalControl uint, isTalented bool) PhysicalSkills {
	return PhysicalSkills{
		empowerment:  empowerment,
		vitalControl: vitalControl,
		isTalented:   isTalented,
	}
}

type MentalSkills struct {
	ilusion       uint
	mentalControl uint
	isTalented    bool
}

func (ms MentalSkills) Ilusion() uint       { return ms.ilusion }
func (ms MentalSkills) MentalControl() uint { return ms.mentalControl }
func (ms MentalSkills) IsTalented() bool    { return ms.isTalented }

func CreateMentalSkillsWithoutValidation(ilusion, mentalControl uint, isTalented bool) MentalSkills {
	return MentalSkills{
		ilusion:       ilusion,
		mentalControl: mentalControl,
		isTalented:    isTalented,
	}
}

type EnergySkills struct {
	objectHandling uint
	energyHandling uint
	isTalented     bool
}

func (es EnergySkills) ObjectHandling() uint { return es.objectHandling }
func (es EnergySkills) EnergyHandling() uint { return es.energyHandling }
func (es EnergySkills) IsTalented() bool     { return es.isTalented }

func CreateEnergySkillsWithoutValidation(objectHandling, energyHandling uint, isTalented bool) EnergySkills {
	return EnergySkills{
		objectHandling: objectHandling,
		energyHandling: energyHandling,
		isTalented:     isTalented,
	}
}

type SpecialStats struct {
	physical         PhysicalSkills
	mental           MentalSkills
	energy           EnergySkills
	energyTank       uint
	isEnergyTalented bool
}

func (ss SpecialStats) Physical() PhysicalSkills { return ss.physical }
func (ss SpecialStats) Mental() MentalSkills     { return ss.mental }
func (ss SpecialStats) Energy() EnergySkills     { return ss.energy }
func (ss SpecialStats) EnergyTank() uint         { return ss.energyTank }
func (ss SpecialStats) IsEnergyTalented() bool   { return ss.isEnergyTalented }

func CreateSpecialStatsWithoutValidation(physical PhysicalSkills, mental MentalSkills, energy EnergySkills, energyTank uint, isEnergyTalented bool) SpecialStats {
	return SpecialStats{
		physical:         physical,
		mental:           mental,
		energy:           energy,
		energyTank:       energyTank,
		isEnergyTalented: isEnergyTalented,
	}
}

type Skill struct {
	transformations []uint
}

func (s Skill) Transformations() []uint { return s.transformations }

func CreateSkillWithoutValidation(transformations []uint) Skill {
	return Skill{
		transformations: transformations,
	}
}

type SupernaturalStats struct {
	skills []Skill
}

func (ss SupernaturalStats) Skills() []Skill { return ss.skills }

func CreateSupernaturalStatsWithoutValidation(skills []Skill) *SupernaturalStats {
	return &SupernaturalStats{
		skills: skills,
	}
}

type XP struct {
	basic        uint
	special      uint
	supernatural uint
}

func (xp XP) Basic() uint        { return xp.basic }
func (xp XP) Special() uint      { return xp.special }
func (xp XP) Supernatural() uint { return xp.supernatural }

func CreateXPWithoutValidation(basic, special, supernatural uint) XP {
	return XP{
		basic:        basic,
		special:      special,
		supernatural: supernatural,
	}
}

type PJ struct {
	id                string
	userID            string
	name              string
	weight            uint
	height            uint
	age               uint
	look              uint
	charisma          int
	villainy          uint
	heroism           uint
	pjType            PJType
	basicStats        BasicStats
	specialStats      SpecialStats
	supernaturalStats *SupernaturalStats
	xp                XP
}

// Getter methods
func (p *PJ) ID() string                            { return p.id }
func (p *PJ) UserID() string                        { return p.userID }
func (p *PJ) Name() string                          { return p.name }
func (p *PJ) Weight() uint                          { return p.weight }
func (p *PJ) Height() uint                          { return p.height }
func (p *PJ) Age() uint                             { return p.age }
func (p *PJ) Look() uint                            { return p.look }
func (p *PJ) Charisma() int                         { return p.charisma }
func (p *PJ) Villainy() uint                        { return p.villainy }
func (p *PJ) Heroism() uint                         { return p.heroism }
func (p *PJ) Type() PJType                          { return p.pjType }
func (p *PJ) BasicStats() BasicStats                { return p.basicStats }
func (p *PJ) SpecialStats() SpecialStats            { return p.specialStats }
func (p *PJ) SupernaturalStats() *SupernaturalStats { return p.supernaturalStats }
func (p *PJ) XP() XP                                { return p.xp }

// CreatePJWithoutValidation creates a PJ instance without validation.
// This function is intended to be used by adapters (like database repositories)
// when reconstructing entities from external sources.
func CreatePJWithoutValidation(
	id string,
	userID string,
	name string,
	weight uint,
	height uint,
	age uint,
	look uint,
	charisma int,
	villainy uint,
	heroism uint,
	pjType PJType,
	basicStats BasicStats,
	specialStats SpecialStats,
	supernaturalStats *SupernaturalStats,
	xp XP,
) *PJ {
	return &PJ{
		id:                id,
		userID:            userID,
		name:              name,
		weight:            weight,
		height:            height,
		age:               age,
		look:              look,
		charisma:          charisma,
		villainy:          villainy,
		heroism:           heroism,
		pjType:            pjType,
		basicStats:        basicStats,
		specialStats:      specialStats,
		supernaturalStats: supernaturalStats,
		xp:                xp,
	}
}

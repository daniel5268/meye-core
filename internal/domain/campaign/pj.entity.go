package campaign

type BasicTalentType string

const (
	BasicTalentPhysical     BasicTalentType = "physical"
	BasicTalentMental       BasicTalentType = "mental"
	BasicTalentCoordination BasicTalentType = "coordination"
	BasicTalentEnergy       BasicTalentType = "energy"
)

type PJType string

const (
	PJTypeHuman        PJType = "human"
	PJTypeSupernatural PJType = "supernatural"
)

type SpecialTalentType string

const (
	SpecialTalentPhysical SpecialTalentType = "physical"
	SpecialTalentMental   SpecialTalentType = "mental"
	SpecialTalentEnergy   SpecialTalentType = "energy"
)

type Physical struct {
	strength   uint
	agility    uint
	speed      uint
	resistance uint
}

func (p Physical) Strength() uint   { return p.strength }
func (p Physical) Agility() uint    { return p.agility }
func (p Physical) Speed() uint      { return p.speed }
func (p Physical) Resistance() uint { return p.resistance }

func CreatePhysicalWithoutValidation(strength, agility, speed, resistance uint) Physical {
	return Physical{
		strength:   strength,
		agility:    agility,
		speed:      speed,
		resistance: resistance,
	}
}

type Mental struct {
	inteligence   uint
	wisdom        uint
	concentration uint
	will          uint
}

func (m Mental) Inteligence() uint   { return m.inteligence }
func (m Mental) Wisdom() uint        { return m.wisdom }
func (m Mental) Concentration() uint { return m.concentration }
func (m Mental) Will() uint          { return m.will }

func CreateMentalWithoutValidation(inteligence, wisdom, concentration, will uint) Mental {
	return Mental{
		inteligence:   inteligence,
		wisdom:        wisdom,
		concentration: concentration,
		will:          will,
	}
}

type Coordination struct {
	precision   uint
	calculation uint
	coordRange  uint
	reflexes    uint
}

func (c Coordination) Precision() uint   { return c.precision }
func (c Coordination) Calculation() uint { return c.calculation }
func (c Coordination) Range() uint       { return c.coordRange }
func (c Coordination) Reflexes() uint    { return c.reflexes }

func CreateCoordinationWithoutValidation(precision, calculation, coordRange, reflexes uint) Coordination {
	return Coordination{
		precision:   precision,
		calculation: calculation,
		coordRange:  coordRange,
		reflexes:    reflexes,
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
}

func (ps PhysicalSkills) Empowerment() uint  { return ps.empowerment }
func (ps PhysicalSkills) VitalControl() uint { return ps.vitalControl }

func CreatePhysicalSkillsWithoutValidation(empowerment, vitalControl uint) PhysicalSkills {
	return PhysicalSkills{
		empowerment:  empowerment,
		vitalControl: vitalControl,
	}
}

type MentalSkills struct {
	ilusion       uint
	mentalControl uint
}

func (ms MentalSkills) Ilusion() uint       { return ms.ilusion }
func (ms MentalSkills) MentalControl() uint { return ms.mentalControl }

func CreateMentalSkillsWithoutValidation(ilusion, mentalControl uint) MentalSkills {
	return MentalSkills{
		ilusion:       ilusion,
		mentalControl: mentalControl,
	}
}

type EnergySkills struct {
	objectHandling uint
	energyHandling uint
}

func (es EnergySkills) ObjectHandling() uint { return es.objectHandling }
func (es EnergySkills) EnergyHandling() uint { return es.energyHandling }

func CreateEnergySkillsWithoutValidation(objectHandling, energyHandling uint) EnergySkills {
	return EnergySkills{
		objectHandling: objectHandling,
		energyHandling: energyHandling,
	}
}

type SpecialStats struct {
	physical   PhysicalSkills
	mental     MentalSkills
	energy     EnergySkills
	energyTank uint
}

func (ss SpecialStats) Physical() PhysicalSkills { return ss.physical }
func (ss SpecialStats) Mental() MentalSkills     { return ss.mental }
func (ss SpecialStats) Energy() EnergySkills     { return ss.energy }
func (ss SpecialStats) EnergyTank() uint         { return ss.energyTank }

func CreateSpecialStatsWithoutValidation(physical PhysicalSkills, mental MentalSkills, energy EnergySkills, energyTank uint) SpecialStats {
	return SpecialStats{
		physical:   physical,
		mental:     mental,
		energy:     energy,
		energyTank: energyTank,
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
	basicTalent       BasicTalentType
	specialTalent     SpecialTalentType
	basicStats        BasicStats
	specialStats      SpecialStats
	supernaturalStats *SupernaturalStats
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
func (p *PJ) BasicTalent() BasicTalentType          { return p.basicTalent }
func (p *PJ) SpecialTalent() SpecialTalentType      { return p.specialTalent }
func (p *PJ) BasicStats() BasicStats                { return p.basicStats }
func (p *PJ) SpecialStats() SpecialStats            { return p.specialStats }
func (p *PJ) SupernaturalStats() *SupernaturalStats { return p.supernaturalStats }

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
	basicTalent BasicTalentType,
	specialTalent SpecialTalentType,
	basicStats BasicStats,
	specialStats SpecialStats,
	supernaturalStats *SupernaturalStats,
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
		basicTalent:       basicTalent,
		specialTalent:     specialTalent,
		basicStats:        basicStats,
		specialStats:      specialStats,
		supernaturalStats: supernaturalStats,
	}
}

type PJCreateParameters struct {
	Name          string
	Weight        uint
	Height        uint
	Age           uint
	Look          uint
	Charisma      int
	Villainy      uint
	Heroism       uint
	PjType        PJType
	BasicTalent   BasicTalentType
	SpecialTalent SpecialTalentType
}

func (c *Campaign) GetPendingUserInvitation(userID string) *Invitation {
	for i := range c.invitations {
		if c.invitations[i].UserID() == userID && c.invitations[i].State() == InvitationStatePending {
			return &c.invitations[i]
		}
	}

	return nil
}

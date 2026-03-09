package person

// PersonDTO - Data Transfer Object для сериализации/десериализации Person
type PersonDTO struct {
	Name       string     `json:"name" xml:"Name" yaml:"name"`
	Type       PersonType `json:"type" xml:"Type" yaml:"type"`
	Health     uint32     `json:"health" xml:"Health" yaml:"health"`
	Mana       uint32     `json:"mana" xml:"Mana" yaml:"mana"`
	Level      uint32     `json:"level" xml:"Level" yaml:"level"`
	Gold       uint32     `json:"gold" xml:"Gold" yaml:"gold"`
	Respect    uint32     `json:"respect" xml:"Respect" yaml:"respect"`
	Strength   uint32     `json:"strength" xml:"Strength" yaml:"strength"`
	Experience uint32     `json:"experience" xml:"Experience" yaml:"experience"`
	HasHouse   bool       `json:"has_house" xml:"HasHouse" yaml:"has_house"`
	HasWeapon  bool       `json:"has_weapon" xml:"HasWeapon" yaml:"has_weapon"`
	HasFamily  bool       `json:"has_family" xml:"HasFamily" yaml:"has_family"`
	X          int32      `json:"x" xml:"X" yaml:"x"`
	Y          int32      `json:"y" xml:"Y" yaml:"y"`
	Z          int32      `json:"z" xml:"Z" yaml:"z"`
}

// ToDTO преобразует Person в PersonDTO
func ToDTO(p Person) PersonDTO {
	return PersonDTO{
		Name:       p.Name(),
		Type:       p.Type(),
		Health:     p.Health(),
		Mana:       p.Mana(),
		Level:      p.Level(),
		Gold:       p.Gold(),
		Respect:    p.Respect(),
		Strength:   p.Strength(),
		Experience: p.Experience(),
		HasHouse:   p.HasHouse(),
		HasWeapon:  p.HasWeapon(),
		HasFamily:  p.HasFamily(),
		X:          p.X(),
		Y:          p.Y(),
		Z:          p.Z(),
	}
}

// FromDTO создает Person из PersonDTO
func FromDTO(dto PersonDTO) (Person, error) {
	return NewPerson(
		WithName(dto.Name),
		WithType(dto.Type),
		WithHealth(dto.Health),
		WithMana(dto.Mana),
		WithLevel(dto.Level),
		WithGold(dto.Gold),
		WithRespect(dto.Respect),
		WithStrength(dto.Strength),
		WithExperience(dto.Experience),
		WithHouse(dto.HasHouse),
		WithWeapon(dto.HasWeapon),
		WithFamily(dto.HasFamily),
		WithCoordinates(dto.X, dto.Y, dto.Z),
	)
}

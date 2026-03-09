package monster

// MonsterDTO - Data Transfer Object для сериализации/десериализации Monster
type MonsterDTO struct {
	Name     string `json:"name" xml:"Name" yaml:"name"`
	Health   uint32 `json:"health" xml:"Health" yaml:"health"`
	Mana     uint32 `json:"mana" xml:"Mana" yaml:"mana"`
	Gold     uint32 `json:"gold" xml:"Gold" yaml:"gold"`
	HasHouse bool   `json:"has_house" xml:"HasHouse" yaml:"has_house"`
	X        int32  `json:"x" xml:"X" yaml:"x"`
	Y        int32  `json:"y" xml:"Y" yaml:"y"`
	Z        int32  `json:"z" xml:"Z" yaml:"z"`
}

// ToDTO преобразует Monster в MonsterDTO
func ToDTO(m Monster) MonsterDTO {
	return MonsterDTO{
		Name:     m.Name(),
		Health:   m.Health(),
		Mana:     m.Mana(),
		Gold:     m.Gold(),
		HasHouse: m.HasHouse(),
		X:        m.X(),
		Y:        m.Y(),
		Z:        m.Z(),
	}
}

// FromDTO создает Monster из MonsterDTO
func FromDTO(dto MonsterDTO) (Monster, error) {
	return NewMonster(
		WithName(dto.Name),
		WithHealth(dto.Health),
		WithMana(dto.Mana),
		WithGold(dto.Gold),
		WithHouse(dto.HasHouse),
		WithCoordinates(dto.X, dto.Y, dto.Z),
	)
}

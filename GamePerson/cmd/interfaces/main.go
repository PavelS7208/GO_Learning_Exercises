package main

// Некие функции уровня бизнес логики работающие с игровыми объектами при помощи интерфейсов

/*
// Функция работает ТОЛЬКО с позицией
func CalculateDistance(a, b Positioned) float64 {
	dx := float64(a.X() - b.X())
	dy := float64(a.Y() - b.Y())
	dz := float64(a.Z() - b.Z())
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// Функция работает только с боем
func ApplyDamage(target Living, damage uint32) error {
	current := target.Health()
	if damage >= current {
		return target.SetHealth(0)
	}
	return target.SetHealth(current - damage)
}

// Функция работает с богатством
func TransferGold(from, to Wealthy, amount uint32) error {
	if from.Gold() < amount {
		return fmt.Errorf("insufficient gold")
	}
	if err := from.SetGold(from.Gold() - amount); err != nil {
		return err
	}
	return to.SetGold(to.Gold() + amount)
}

// Система опыта
type ExperienceSystem struct{}

func (es *ExperienceSystem) AddExperience(char Experienced, exp uint32) error {
	newExp := char.Experience() + exp
	if err := char.SetExperience(newExp); err != nil {
		return err
	}

	// Повышение уровня
	if newExp >= 10 && char.Level() < 10 {
		return char.SetLevel(char.Level() + 1)
	}
	return nil
}

// Использует только нужные интерфейсы
func (es *ExperienceSystem) Process(char Experienced) {
	// Работает и с Person, и с другими сущностями
}
*/

/*
// Работает и с Person, и с Monster, и с любым Positioned
func IsNearby(a, b Positioned, maxDistance float64) bool {
	return CalculateDistance(a, b) <= maxDistance
}

// Работает с любым Living
func Heal(target Living, amount uint32) error {
	current := target.Health()
	// ...
}


func RenderOnMap(entities []Actor) {
    // Работает с Person, Monster, NPC  .....
}


func MoveToTarget(p Movable, target Positioned) { ... }


// Mock для тестов

type mockLiving struct {
    health uint32
}

func (m *mockLiving) Health() uint32 { return m.health }
func (m *mockLiving) SetHealth(h uint32) error {
    m.health = h
    return nil
}

func TestHeal(t *testing.T) {
    target := &mockLiving{health: 50}
    Heal(target, 30)
    assert.Equal(t, uint32(80), target.Health())
}


*/

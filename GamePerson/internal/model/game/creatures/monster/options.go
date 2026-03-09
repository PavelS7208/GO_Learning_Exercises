package monster

import "fmt"

// ------------------ Для конструктора---------------------------------
// ----------------- Опции свойств -------------------------------------

func WithName(name string) Option {
	return func(p *monster) error {
		if err := p.SetName(name); err != nil {
			return fmt.Errorf("failed to apply option WithName: %w", err)
		}
		return nil
	}
}

func WithCoordinates(x, y, z int32) Option {
	return func(m *monster) error {
		if err := m.SetX(x); err != nil {
			return fmt.Errorf("failed to apply option WithCoordinates (X): %w", err)
		}
		if err := m.SetY(y); err != nil {
			return fmt.Errorf("failed to apply option WithCoordinates (Y): %w", err)
		}
		if err := m.SetZ(z); err != nil {
			return fmt.Errorf("failed to apply option WithCoordinates (Z): %w", err)
		}
		return nil
	}
}

func WithGold(gold uint32) Option {
	return func(m *monster) error {
		if err := m.SetGold(gold); err != nil {
			return fmt.Errorf("WithGold: %w", err)
		}
		return nil
	}
}

func WithHealth(health uint32) Option {
	return func(m *monster) error {
		if err := m.SetHealth(health); err != nil {
			return fmt.Errorf("WithHealth: %w", err)
		}
		return nil
	}
}

func WithMana(mana uint32) Option {
	return func(m *monster) error {
		if err := m.SetMana(mana); err != nil {
			return fmt.Errorf("WithMana: %w", err)
		}
		return nil
	}
}

func WithHouse(has bool) Option {
	return func(m *monster) error {
		if err := m.SetHouse(has); err != nil {
			return fmt.Errorf("WithHouse: %w", err)
		}
		return nil
	}
}

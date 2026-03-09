package person

import "fmt"

func WithName(name string) Option {
	return func(p *person) error {
		if err := p.SetName(name); err != nil {
			return fmt.Errorf("WithName: %w", err)
		}
		return nil
	}
}

func WithLevel(level uint32) Option {
	return func(p *person) error {
		if err := p.SetLevel(level); err != nil {
			return fmt.Errorf("WithLevel: %w", err)
		}
		return nil
	}
}

func WithHealth(health uint32) Option {
	return func(p *person) error {
		if err := p.SetHealth(health); err != nil {
			return fmt.Errorf("WithHealth: %w", err)
		}
		return nil
	}
}

func WithMana(mana uint32) Option {
	return func(p *person) error {
		if err := p.SetMana(mana); err != nil {
			return fmt.Errorf("WithMana: %w", err)
		}
		return nil
	}
}

func WithGold(gold uint32) Option {
	return func(p *person) error {
		if err := p.SetGold(gold); err != nil {
			return fmt.Errorf("WithGold: %w", err)
		}
		return nil
	}
}

func WithStrength(strength uint32) Option {
	return func(p *person) error {
		if err := p.SetStrength(strength); err != nil {
			return fmt.Errorf("WithStrength: %w", err)
		}
		return nil
	}
}

func WithRespect(respect uint32) Option {
	return func(p *person) error {
		if err := p.SetRespect(respect); err != nil {
			return fmt.Errorf("WithRespect: %w", err)
		}
		return nil
	}
}

func WithExperience(exp uint32) Option {
	return func(p *person) error {
		if err := p.SetExperience(exp); err != nil {
			return fmt.Errorf("WithExperience: %w", err)
		}
		return nil
	}
}

func WithType(pt PersonType) Option {
	return func(p *person) error {
		if err := p.SetType(pt); err != nil {
			return fmt.Errorf("WithType: %w", err)
		}
		return nil
	}
}

func WithCoordinates(x, y, z int32) Option {
	return func(p *person) error {
		if err := p.SetX(x); err != nil {
			return fmt.Errorf("WithCoordinates.X: %w", err)
		}
		if err := p.SetY(y); err != nil {
			return fmt.Errorf("WithCoordinates.Y: %w", err)
		}
		if err := p.SetZ(z); err != nil {
			return fmt.Errorf("WithCoordinates.Z: %w", err)
		}
		return nil
	}
}

func WithHouse(has bool) Option {
	return func(p *person) error {
		if err := p.SetHouse(has); err != nil {
			return fmt.Errorf("WithHouse: %w", err)
		}
		return nil
	}
}

func WithWeapon(has bool) Option {
	return func(p *person) error {
		if err := p.SetWeapon(has); err != nil {
			return fmt.Errorf("WithWeapon: %w", err)
		}
		return nil
	}
}

func WithFamily(has bool) Option {
	return func(p *person) error {
		if err := p.SetFamily(has); err != nil {
			return fmt.Errorf("WithFamily: %w", err)
		}
		return nil
	}
}

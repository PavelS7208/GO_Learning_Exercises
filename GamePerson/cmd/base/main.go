package main

import (
	"GamePerson/internal/model/game/creatures/monster"
	"GamePerson/internal/model/game/creatures/person"
	"fmt"
)

func main() {

	p, err := person.NewPerson(createGamePersonConfig()...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(p)

	m, err1 := monster.NewMonster(createGameMonsterConfig()...)
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}
	fmt.Println(m)

	if err = p.SetGold(2000); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)

	if err = m.SetHealth(30); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(m)
}

func createGamePersonConfig() []person.Option {

	const x, y, z = -2_000, 1_000, 0
	const name = "Wizard First 353535                                       3"
	const personType = person.PersonTypeBuilder
	const gold = 2_000
	const mana = 50
	const health = 100
	const respect = 1
	const strength = 1
	const experience = 1
	const level = 1

	return []person.Option{
		person.WithName(name),
		person.WithType(personType),
		person.WithCoordinates(x, y, z),
		person.WithGold(gold),
		person.WithMana(mana),
		person.WithHealth(health),
		person.WithRespect(respect),
		person.WithStrength(strength),
		person.WithExperience(experience),
		person.WithLevel(level),
		person.WithHouse(true),
		person.WithFamily(true),
	}
}

func createGameMonsterConfig() []monster.Option {

	const x, y, z = -3_000, 1_000, 0
	const name = "Monster Vasya"
	const gold = 5_000
	const mana = 500
	const health = 300

	return []monster.Option{
		monster.WithName(name),
		monster.WithCoordinates(x, y, z),
		monster.WithGold(gold),
		monster.WithMana(mana),
		monster.WithHealth(health),
		monster.WithHouse(true),
	}
}

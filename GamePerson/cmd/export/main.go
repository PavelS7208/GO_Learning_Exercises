package main

import (
	"GamePerson/internal/model/game/creatures/base/entity"
	"GamePerson/internal/model/game/creatures/monster"
	"GamePerson/internal/model/game/creatures/person"
	"fmt"
	"os"
)

func main() {
	// Единый инстанс проверки целостности
	ic := entity.NewIntegrityChecker()

	// Создаём сериализаторы
	personSerializer := person.NewSerializer(ic)
	monsterSerializer := monster.NewSerializer(ic)

	// === Создание персонажа ===
	p, err := person.NewPerson(
		person.WithName("BuilderBob"),
		person.WithType(person.PersonTypeBuilder),
		person.WithCoordinates(100, 200, 50),
		person.WithGold(5000),
	)
	if err != nil {
		panic(err)
	}

	// === Сериализация в JSON ===
	jsonData, err := personSerializer.ToJSON(p)
	if err != nil {
		panic(err)
	}
	fmt.Println("Person JSON:", string(jsonData))

	// === Десериализация с валидацией ===
	// Атака: подмена здоровья в сохранении
	maliciousJSON := []byte(`{"name":"Hacker","type":0,"health":999999,"mana":100,"level":1,"gold":0,"respect":0,"strength":5,"experience":0,"has_house":false,"has_weapon":true,"has_family":false,"x":0,"y":0,"z":0}`)

	_, err = personSerializer.FromJSON(maliciousJSON)
	if err != nil {
		fmt.Printf("✅ Атака кривыми данными заблокирована: %v\n", err)
		// Вывод: "JSON integrity check failed: health 999999 > max 1000"
	} else {
		fmt.Println("❌ УЯЗВИМОСТЬ: атака кривыми данными прошла!")
	}

	// Типичный случай создания — просто и понятно
	personJSON := []byte(`{"name":"Aragon","type":2, "health":500 }`)
	per, err := person.NewFromJSON(personJSON)
	if err != nil {
		panic(err)
	}
	fmt.Println(per)

	// Специальный случай — кастомная валидация (не реализовано, типа будущее развитие)
	//checker := &entity.LoggingIntegrityChecker{Logger: auditLog}
	//m, err := monster.NewFromJSONWithIntegrity(jsonData, checker)

	// === Работа с монстром ===
	m, err := monster.NewMonster(
		monster.WithName("Dragon"),
		//monster.WithAggression(95),
		monster.WithCoordinates(500, 600, 100),
	)
	if err != nil {
		panic(err)
	}

	// Сериализация монстра
	monsterJSON, _ := monsterSerializer.ToJSON(m)

	// Десериализация с той же логикой валидации
	m2, err := monsterSerializer.FromJSON(monsterJSON)
	if err != nil {
		panic(err) // Если данные валидны — ошибки не будет
	}
	_ = m2 // восстановленный монстр

	// === Сохранение в файл ===
	if err := os.WriteFile("save_person.json", jsonData, 0644); err != nil {
		panic(err)
	}
}

/*

func main() {
	// Примеры сериализации Person
	fmt.Println("=== Person Serialization Examples ===")
	demonstratePersonSerialization()

	fmt.Println("\n=== Monster Serialization Examples ===")
	demonstrateMonsterSerialization()
}

func demonstratePersonSerialization() {
	// Создаем Person
	p, err := person.NewPerson(
		person.WithName("Aragon"),
		person.WithType(person.PersonTypeWarrior),
		person.WithHealth(850),
		person.WithMana(300),
		person.WithLevel(8),
		person.WithGold(5000),
		person.WithStrength(9),
		person.WithExperience(7),
		person.WithRespect(10),
		person.WithWeapon(true),
		person.WithHouse(true),
		person.WithFamily(true),
		person.WithCoordinates(100, 200, 50),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Original Person:", p)
	fmt.Println()

	serializer := person.NewSerializer()

	// JSON
	fmt.Println("--- JSON ---")
	jsonData, err := serializer.ToJSON(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

	// Десериализация из JSON
	pFromJSON, err := serializer.FromJSON(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deserialized from JSON:", pFromJSON)
	fmt.Println()

	// XML
	fmt.Println("--- XML ---")
	xmlData, err := serializer.ToXML(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(xmlData))

	// Десериализация из XML
	pFromXML, err := serializer.FromXML(xmlData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deserialized from XML:", pFromXML)
	fmt.Println()

	// YAML
	fmt.Println("--- YAML ---")
	yamlData, err := serializer.ToYAML(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(yamlData))

	// Десериализация из YAML
	pFromYAML, err := serializer.FromYAML(yamlData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deserialized from YAML:", pFromYAML)
}

func demonstrateMonsterSerialization() {
	// Создаем Monster
	m, err := monster.NewMonster(
		monster.WithName("Dragon"),
		monster.WithHealth(500),
		monster.WithMana(800),
		monster.WithGold(10000),
		monster.WithHouse(false),
		monster.WithCoordinates(-50, 300, 150),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Original Monster:", m)
	fmt.Println()

	serializer := monster.NewSerializer()

	// JSON
	fmt.Println("--- JSON ---")
	jsonData, err := serializer.ToJSON(m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

	// Десериализация из JSON
	mFromJSON, err := serializer.FromJSON(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deserialized from JSON:", mFromJSON)
	fmt.Println()

	// XML
	fmt.Println("--- XML ---")
	xmlData, err := serializer.ToXML(m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(xmlData))

	// Десериализация из XML
	mFromXML, err := serializer.FromXML(xmlData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deserialized from XML:", mFromXML)
	fmt.Println()

	// YAML
	fmt.Println("--- YAML ---")
	yamlData, err := serializer.ToYAML(m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(yamlData))

	// Десериализация из YAML
	mFromYAML, err := serializer.FromYAML(yamlData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deserialized from YAML:", mFromYAML)
}
*/

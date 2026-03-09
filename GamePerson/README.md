# GamePersons — пример системы управления игровыми персонажами в Go

## Постановка задачи

За основу взято одно из домашних заданий курса "Глубокий GO" от Balun-а с существенными доработками требований к системе.

**Базовое задание:**
Разработать систему для хранения игровых персонажей с ограничением в **64 байта на персонажа** без использования heap-аллокаций и внешних ссылок. Сериализация в JSON/XML/YAML. Паттерн Functional Options для создания объектов.

**Обязательные атрибуты:**
```
Координаты X, Y, Z:  [-2_000_000_000 … 2_000_000_000]  → int32 (по 4 байта)
Золото:              [0 … 2_000_000_000]               → uint32 (4 байта)
Мана:                [0 … 1000]                        → 10 бит
Здоровье:            [0 … 1000]                        → 10 бит
Уровень:             [0 … 10]                          → 4 бита
Сила:                [0 … 10]                          → 4 бита
Опыт:                [0 … 10]                          → 4 бита
Уважение:            [0 … 10]                          → 4 бита
Тип персонажа:       [строитель/кузнец/воин]           → 2 бита
Флаги (дом/оружие/семья):                              → 3 бита
Имя:                 [до 42 символов ASCII]            → 42 байта
```

**Расширенные требования:**
- Поддержка нескольких типов существ (Person, Monster и т.д.)
- Расширенные требования к созданию через Паттерн Functional Options. Дефолтные значения при создании
- Сериализация в JSON/XML/YAML. Создание персонажей из JSON/XML/YAML 
- Валидация на всех уровнях
- Расширяемая архитектура
- Покрытие тестами

---

## Основные результаты работы

### 1. **Собственная библиотека BitPack**

Разработана с нуля библиотека для работы с битовыми полями:
- Поддержка булевых (один бит), знаковых и беззнаковых целых чисел (2-64 бита)
- Generic API с compile-time проверкой типов
- [Подробная документация BitPack](GamePerson/internal/bitpack/BitPack_readme.md)


### 2. Архитектура

```
┌─────────────────────────────────────────────────────────┐
│  cmd/                    (Application Layer)            │
│  ├── base/      - Демонстрация базовых операций         │
│  ├── export/    - Пример сериализации                   │
│  └── interfaces/- Работа с интерфейсами                 │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│  internal/model/game/creatures/  (Domain Layer)         │
│  ├── base/entity/    - Интерфейсы и валидация           │
│  ├── base/serializer/- Generic сериализация             │
│  ├── person/         - Бизнес-логика персонажей         │
│  └── monster/        - Бизнес-логика монстров           │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│  internal/model/bitpack/  (Infrastructure Layer)        │
│  ├── person/   - Схема упаковки для Person              │
│  └── monster/  - Схема упаковки для Monster             │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│  internal/bitpack/        (Core Library)                │
│  - Низкоуровневая работа с битами                       │
└─────────────────────────────────────────────────────────┘
```

**Принципы дизайна:**
- **Разделение ответственности:** битовая упаковка ≠ бизнес-логика
- **Инкапсуляция:** приватные структуры, публичные интерфейсы
- **Расширяемость:** добавление нового типа требует минимум изменений
- **Валидация на всех уровнях:** битовые поля → бизнес-логика → сериализация

### 3. Битовая структура персонажа

**Person (64 байта):**
```go
type person struct {
    name   [42]byte  // 42: имя (без heap-аллокации!)
    packed [6]byte   //  6: упакованные атрибуты (48 бит)
    gold   uint32    //  4: золото
    x, y, z int32    // 12: координаты
    // Total: 64 байта
}
```

**Битовая схема для packed (48 бит):**
```
Биты  0- 5: длина имени     (6 бит)
Биты  6- 9: уважение        (4 бита)
Биты 10-13: сила            (4 бита)
Биты 14-17: опыт            (4 бита)
Биты 18-21: уровень         (4 бита)
Биты 22-23: тип персонажа   (2 бита)
Бит  24:    есть дом        (1 бит)
Бит  25:    есть оружие     (1 бит)
Бит  26:    есть семья      (1 бит)
Биты 27-36: мана            (10 бит)
Биты 37-46: здоровье        (10 бит)
Бит  47:    резерв          (1 бит)
```

**Monster (64 байта с явным паддингом):**
```go
type monster struct {
    name   [42]byte  // 42: имя
    _      [2]byte   //  2: явный padding для выравнивания
    packed [4]byte   //  4: упакованные атрибуты (32 бита)
    gold   uint32    //  4: золото
    x, y, z int32    // 12: координаты
    // Total: 64 байта
}
```

### 4. **Интерфейсы персонажей**

```go
// Композиция интерфейсов для гибкости
type Person interface {
    entity.Combatant      // Named + Positioned + Movable + Living + Strong + Armed
    entity.Wealthy        // Gold management
    entity.Magical        // Mana management
    entity.Experienced    // Level & Experience
    entity.Reputable      // Respect
    entity.PropertyOwner  // House
    entity.FamilyMember   // Family
    
    Type() PersonType
    SetType(PersonType) error
}
```

**Compile-time проверка соответствия:**
```go
var _ entity.Combatant = (*person)(nil)
var _ Person = (*person)(nil)
// Компилятор гарантирует реализацию всех методов
```

### 5. **Functional Options Pattern**  для создания объекта

```go
// Создание с явными настройками
person, err := person.NewPerson(
    person.WithName("Gandalf"),
    person.WithType(person.PersonTypeBuilder),
    person.WithCoordinates(-2000, 1000, 0),
    person.WithHealth(100),
    person.WithMana(50),
    person.WithLevel(10),
    person.WithHouse(true),
)

// Гибкость: порядок опций не важен
// Безопасность: валидация каждой опции
// Расширяемость: добавление новых опций не ломает существующий код
```

### 6. Обобщенная система сериализации на основе дженериков

```go
// Универсальный сериализатор для любых сущностей
type Serializer[E any, D any] struct {
    converter Converter[E, D]
    integrity *entity.IntegrityChecker
}

// Единый инстанс проверки целостности
	ic := entity.NewIntegrityChecker()

	// Создаём сериализаторы
	personSerializer := person.NewSerializer(ic)
	monsterSerializer := monster.NewSerializer(ic)

// Использование:
	p, _ := person.NewPerson(
		person.WithName("BuilderBob"),
		person.WithType(person.PersonTypeBuilder),
		person.WithCoordinates(100, 200, 50),
		person.WithGold(5000),
	)

	// === Сериализация в JSON ===
	jsonData, err := personSerializer.ToJSON(p)

     m, _ := monster.NewMonster(
		monster.WithName("Dragon"),
		monster.WithCoordinates(500, 600, 100),
	)
	
	// Сериализация монстра
	monsterJSON, _ := monsterSerializer.ToJSON(m)

	// Десериализация с той же логикой валидации
	m2, err := monsterSerializer.FromJSON(monsterJSON)
	
//  Или так. Дефолтная проверка валидации
serializer := person.NewSerializer(nil)
jsonData, _ := serializer.ToJSON(p)
xmlData, _ := serializer.ToXML(p)
yamlData, _ := serializer.ToYAML(p)

// Десериализация с автоматической валидацией
restored, err := person.NewFromJSON(jsonData)
```

### 7. **Валидации**

```go
// Уровень 1: Валидация конфигурации битовых полей (compile-time)
var healthField = bitpack.MustNewUIntBitField(37, 46, 1000)
// Паникует при старте, если конфигурация некорректна

// Уровень 2: Бизнес-валидация (runtime)
func (p *person) SetHealth(health uint32) error {
    if health > config.PersonMaxHealth {
        return fmt.Errorf("health %d exceeds maximum %d", health, config.PersonMaxHealth)
    }
    personbitpack.SetHealthUnchecked(&p.packed, health)
    return nil
}

// Уровень 3: Валидация целостности (после десериализации)
func (p *person) Validate() error {
    // Проверка инвариантов: длина имени, мусор в буфере, диапазоны
}
```

---

## Технические решения

### Checked vs Unchecked операции

**Performance-критичные участки:**
```go
// Бизнес-слой делает валидацию → использует unchecked для скорости
func (p *person) SetMana(mana uint32) error {
    if mana > config.PersonMaxMana {
        return fmt.Errorf("mana exceeds maximum")
    }
    // Гарантия валидности → пропускаем повторную проверку
    personbitpack.SetManaUnchecked(&p.packed, mana)  // ~2x быстрее
    return nil
}
```

### Использование промежуточного слоя (DTO) для сериализации сложной структуры персонажа
```go
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
func ToDTO(p Person) PersonDTO 

// FromDTO создает Person из PersonDTO
func FromDTO(dto PersonDTO) (Person, error)
```


## Структура проекта

```
GamePerson/
├── cmd/
│   ├── base/           # Демонстрация базовых операций
│   ├── export/         # Примеры сериализации
│   └── interfaces/     # Работа с интерфейсами
├── internal/
│   ├── bitpack/        #  Библиотека битовых полей
│   │   ├── bit_field.go         # UInt/Int/BoolBitField
│   │   ├── bit_field_error.go   # Структурированные ошибки библиотеки
│   │   ├── bit_pack.go          # API для работы с битовыми полями
│   │   └── types.go             # Базовые типы
│   └── model/
│       ├── config/              # Конфигурация лимитов
│       ├── bitpack/
│       │   ├── person/          # Схема упаковки Person
│       │   └── monster/         # Схема упаковки Monster
│       └── game/creatures/
│           ├── base/
│           │   ├── entity/      # Интерфейсы (Combatant, Living, etc)
│           │   └── serializer/  # Generic сериализатор
│           ├── person/          # Реализация Person
│           │   ├── person.go
│           │   ├── attributes.go
│           │   ├── options.go
│           │   ├── dto.go
│           │   └── serialize.go
│           └── monster/         # Реализация Monster
│               └── ...
└── README.md
```

---

## Примеры использования

### Создание персонажа

```go
person, err := person.NewPerson(
    person.WithName("Gandalf the Grey"),
    person.WithType(person.PersonTypeBuilder),
    person.WithCoordinates(100, 200, 0),
    person.WithHealth(100),
    person.WithMana(50),
    person.WithLevel(5),
    person.WithStrength(8),
    person.WithExperience(7),
    person.WithRespect(9),
    person.WithGold(1000),
    person.WithHouse(true),
    person.WithWeapon(true),
    person.WithFamily(false),
)
```

### Сериализация

```go
// В JSON
serializer := person.NewSerializer(nil)
jsonData, _ := serializer.ToJSON(person)

// Из JSON с валидацией
restored, err := person.NewFromJSON(jsonData)
if err != nil {
    log.Fatal("Deserialization failed:", err)
}
```

### Работа с интерфейсами

```go
// Полиморфизм через интерфейсы
func HealCreature(c entity.Living) {
    current := c.Health()
    c.SetHealth(current + 10)
}

HealCreature(person)  //  работает
HealCreature(monster) //  работает
```

---

## Возможные улучшения

- Добавить поддержку concurrent доступов (sync.RWMutex wrapper и т.д.)
- Поддержка эффективного создания и работы с большим кол-вом объектов 
- Реализовать persistence layer (сохранение в БД/файлы)
- Добавить события и event sourcing
- Расширить систему типов (Elf, Dwarf, Dragon, и т.д.)
- Benchmark, Fuzz-тесты






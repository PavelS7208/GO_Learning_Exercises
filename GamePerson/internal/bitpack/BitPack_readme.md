# BitPack — библиотека работы с битовыми полями в GO

Библиотека **BitPack** предоставляет типобезопасные и производительные примитивы для работы с битовыми полями в Go. Позволяет компактно хранить несколько независимых значений (целые числа, флаги) внутри одного 64-битного слова или среза байтов (от 1 до 8) с автоматической валидацией диапазонов и корректной обработкой знаковых значений.

## Содержание

- [Архитектура](#архитектура)
- [API для работы](#api-для-работы)
- [Сценарии работы](#сценарии-работы)
- [Обработка ошибок](#обработка-ошибок)
- [Конкурентность](#конкурентность)
- [Примеры использования](#примеры-использования)


## Архитектура

Библиотека построена на основе битовых полей, которые инкапсулируют конфигурацию (позиции битов, диапазоны значений) и предоставляют безопасные методы для работы с данными.

### Базовые типы

| Тип | Описание | Применение |
|-----|----------|------------|
| `BitSet64` | `uint64` — основной контейнер для 64 бит | Хранение упакованных данных в памяти |
| `PackedN` | `[N/8]byte` (N=8,16,24,32,40,48,56,64) | Типобезопасное представление для сериализации |
| `BitPosition` | `uint8` (type alias) — позиция бита (0-63) | Указание границ битовых полей |

**Доступные типы PackedN:**
```go
type Packed8  [1]byte // 8 бит
type Packed16 [2]byte // 16 бит
type Packed24 [3]byte // 24 бита
type Packed32 [4]byte // 32 бита
type Packed40 [5]byte // 40 бит
type Packed48 [6]byte // 48 бит
type Packed56 [7]byte // 56 бит
type Packed64 [8]byte // 64 бита
```

**Пример использования `PackedN` типов:**
```go
var buffer bitpack.Packed24 // Типобезопасный 3-байтный контейнер (24 бита)

// Запись значения
bitset := bitpack.BitSet64(0x123456)
bitpack.PackBytes(buffer[:], bitset)

// Чтение значения
restored := bitpack.UnpackBytes(buffer[:]) // 0x123456
```

### Типы битовых полей

Библиотека предоставляет три специализированных типа для работы с разными видами данных:

#### `UIntBitField` — беззнаковые целые

Предназначен для хранения неотрицательных целых чисел с контролем максимального значения.

```go
// Конструктор с валидацией диапазона
field, err := bitpack.NewUIntBitField(0, 11, 4095) // 12 бит, макс. 4095
if err != nil {
    log.Fatal(err)
}

// Извлечение значения
value := field.Get(bitset) // uint64

// Запись значения с валидацией
bitset, err = field.Update(bitset, 1234)

// Запись без валидации (быстрее, но требует уверенности в корректности значения)
bitset = field.UpdateUnchecked(bitset, 1234)
```

**Дополнительные методы:**
- `Width() uint8` — возвращает ширину поля в битах
- `String() string` — строковое представление для отладки

#### `IntBitField` — знаковые целые

Реализует хранение целых чисел со знаком в дополнительном коде с автоматическим знаковым расширением при извлечении.

```go
// Конструктор с явными границами диапазона
field, err := bitpack.NewIntBitField(0, 7, -100, 100) // 8 бит, от -100 до 100
if err != nil {
    log.Fatal(err)
}

// Конструктор с полным диапазоном для заданной ширины
field, err := bitpack.NewIntBitFieldAuto(0, 7) // 8 бит: -128..127

// Извлечение с автоматическим знаковым расширением
value := field.Get(bitset) // int64

// Запись значения с валидацией
bitset, err = field.Update(bitset, -42)

// Запись без валидации
bitset = field.UpdateUnchecked(bitset, -42)
```

**Алгоритм знакового расширения:**
```
Для 4-битного поля со значением 0b1101 (-3):

1. Извлечение:       (bitset >> start) & mask → 0b1101 (13)
2. Сдвиг влево на 60: 13 << 60 → 0xD000_0000_0000_0000
3. Арифметический сдвиг вправо на 60: 
   знаковый бит (1) распространяется → 0xFFFF_FFFF_FFFF_FFFD = -3
```

#### `BoolBitField` — булевы флаги

Оптимизированное представление для хранения одиночных битов-флагов.

```go
flag, _ := bitpack.NewBoolBitField(15)

// Чтение значения
if flag.Get(bitset) {
    // флаг установлен
}

// Запись значения (всегда без ошибок)
bitset, _ = flag.Update(bitset, true)
bitset = flag.UpdateUnchecked(bitset, true)

// Оптимизированные операции
bitset = flag.Set(bitset)    // Установить в true
bitset = flag.Clear(bitset)  // Установить в false
bitset = flag.Toggle(bitset) // Инвертировать
```

## API для работы

Библиотека предоставляет методы для работы с данными, упакованными в байтовые срезы.

### Generic функции

generic-функции для работы с различными типами:

```go
// Чтение беззнаковых значений
func GetUIntFieldAs[T UnsignedInteger](packed []byte, field UIntBitField) T

// Чтение знаковых значений
func GetIntFieldAs[T SignedInteger](packed []byte, field IntBitField) T

// Чтение булевых значений
func GetBoolField(packed []byte, field BoolBitField) bool

// Запись с валидацией
func SetUIntFieldAs[T UnsignedInteger](packed []byte, field UIntBitField, value T) error
func SetIntFieldAs[T SignedInteger](packed []byte, field IntBitField, value T) error
func SetBoolField(packed []byte, field BoolBitField, value bool) error

// Запись без валидации (быстрее)
func SetUIntFieldUncheckedAs[T UnsignedInteger](packed []byte, field UIntBitField, value T)
func SetIntFieldUncheckedAs[T SignedInteger](packed []byte, field IntBitField, value T)
func SetBoolFieldUnchecked(packed []byte, field BoolBitField, value bool)
```

**Поддерживаемые типы:**
- `UnsignedInteger`: `uint8`, `uint16`, `uint32`, `uint64`
- `SignedInteger`: `int8`, `int16`, `int32`, `int64`

**Пример использования:**
```go
packet := make([]byte, 8) // 64 бита

field := bitpack.MustNewUIntBitField(0, 15, 65535)

// Работа с uint16
err := bitpack.SetUIntFieldAs[uint16](packet, field, 42)
value := bitpack.GetUIntFieldAs[uint16](packet, field) // 42

// Работа с int32
tempField := bitpack.MustNewIntBitField(16, 31, -1000, 1000)
err = bitpack.SetIntFieldAs[int32](packet, tempField, -250)
temp := bitpack.GetIntFieldAs[int32](packet, tempField) // -250
```

### Checked vs Unchecked операции

Библиотека предоставляет два варианта операций записи для оптимизации производительности:

#### Checked (с валидацией)

**Проверки:**
- Длина slice (`len(packed) > 0` и `len(packed) <= 8`)
- Диапазон значения (соответствие `Min`/`Max` или `Max` для беззнаковых)

**Пример:**
```go
// Внешние данные - всегда проверяем
err := bitpack.SetUIntFieldAs[uint16](packet, field, userInput)
if err != nil {
    return fmt.Errorf("invalid value: %w", err)
}
```

#### Unchecked (без валидации)

**Что пропускается:**
- Проверка длины slice
- Проверка диапазона значения

**Когда использовать:**
- Данные уже валидированы на уровне бизнес-логики
- Batch обработка больших объемов данных
- Критичные по производительности участки

**Пример:**
```go
// Данные уже проверены бизнес-логикой
func (p *Person) SetHealth(health uint32) error {
    // Бизнес-валидация
    if health > MaxHealth {
        return fmt.Errorf("health exceeds maximum")
    }
    
    // Гарантия корректности → используем unchecked для скорости
    bitpack.SetUIntFieldUncheckedAs[uint32](p.packed[:], healthField, health)
    return nil
}
```

## Сценарии работы

### Валидация и безопасность

#### 1. Валидация конфигурации (при создании поля)

```go
// Некорректный диапазон битов
_, err := bitpack.NewUIntBitField(10, 5, 255)
// err: "bit field range error: start position (10) must be <= end position (5)"

// Выход за пределы 64 бит
_, err := bitpack.NewUIntBitField(0, 70, 255)
// err: "bit field range error: end position (70) must be < 64"

// Максимум превышает ёмкость поля
_, err := bitpack.NewUIntBitField(0, 3, 20) // 4 бита = макс 15
// err: "value 20 exceeds capacity of 4-bit field (max allowed: 15)"
```

#### 2. Валидация значений (при записи Checked операциями)

```go
field := bitpack.MustNewUIntBitField(0, 3, 10) // 4 бита, но макс 10

// Попытка записать значение вне диапазона
_, err := field.Update(0, 15)
// err: "value 15 exceeds capacity of 4-bit field (max allowed: 10)"

// Checked операции также проверяют длину slice
var packet []byte // пустой slice
err := bitpack.SetUIntFieldAs[uint8](packet, field, 5)
// err: "packed slice is empty"

packet = make([]byte, 10) // слишком большой
err = bitpack.SetUIntFieldAs[uint8](packet, field, 5)
// err: "packed slice too large: 10 bytes (max 8)"
```


#### 3. Использование Must* конструкторов для статических конфигураций

Это переносит валидацию на этап компиляции/инициализации для глобальных констант:

```go
// В глобальных переменных или константах структуры:
var (
    VersionField   = bitpack.MustNewUIntBitField(0, 3, 15)
    PriorityField  = bitpack.MustNewUIntBitField(4, 6, 7)
    EncryptedFlag  = bitpack.MustNewBoolBitField(7)
)

// Если конфигурация невалидна - программа упадет при запуске
// Это ЖЕЛАЕМОЕ поведение для статических конфигураций!
```

#### 4. Использование PackedN типов

Это предотвращает ошибки работы с неправильной длиной среза:

```go
var buffer bitpack.Packed32 // Типобезопасный 4-байтный контейнер

// Запись
bitpack.PackBytes(buffer[:], bitset)

// Чтение
bitset := bitpack.UnpackBytes(buffer[:])
```




## Обработка ошибок

### Типы ошибок

```go
const (
    KindStartAfterEnd      // start > end в конфигурации
    KindEndOutOfRange      // end >= 64
    KindValueOverflow      // значение превышает максимум
    KindValueUnderflow     // значение меньше минимума (signed)
    KindValueRangeInverted // min > max в конфигурации
    KindPositionOutOfRange // позиция >= 64 для bool поля
    KindSliceEmpty         // len(packed) == 0
    KindSliceTooLarge      // len(packed) > 8
)
```

### Примеры обработки ошибок

```go
// Проверка типа ошибки
err := bitpack.SetUIntFieldAs[uint32](packet, field, value)
if errors.Is(err, bitpack.ErrValueOverflow) {
    log.Printf("Value too large for field")
}

// Извлечение деталей
var bpErr *bitpack.Error
if errors.As(err, &bpErr) {
    switch bpErr.Kind {
    case bitpack.KindValueOverflow:
        log.Printf("Value %d exceeds max %d for %d-bit field",
            bpErr.Details.Value, bpErr.Details.AllowedMax, bpErr.Details.BitWidth)
    case bitpack.KindSliceTooLarge:
        log.Printf("Slice too large: %d bytes (max 8)",
            bpErr.Details.SliceLength)
    }
}
```

### Проверка на непересечение полей

При проектировании сложных структур с множеством битовых полей важно проверять на **Непересечение** — битовые поля не должны перекрываться в хранилище, чтобы изменение одного поля не приводило к изменению другого. На уровне кода работы с полем работа идет как с атомарным объектом, проверки что диапазон битов может быть общим - не происходит. Рекомендуется добавлять unit-тест, проверяющий корректность раскладки полей:

```go
package mypackage

import (
	"testing"
	"bitpack"
)

// Тестовая структура: компактный заголовок сетевого пакета (32 бита)
var (
    VersionField  = bitpack.MustNewUIntBitField(0, 3, 15)     // 4 бита
    PriorityField = bitpack.MustNewUIntBitField(4, 6, 7)      // 3 бита
    EncryptedFlag = bitpack.MustNewBoolBitField(7)            // 1 бит
    SequenceField = bitpack.MustNewUIntBitField(8, 23, 65535) // 16 бит
    ReservedField = bitpack.MustNewUIntBitField(24, 30, 127)  // 7 бит
    LastFlag      = bitpack.MustNewBoolBitField(31)           // 1 бит
)

// TestPacketHeaderLayout проверяет непересечение и полное покрытие битов
func TestPacketHeaderLayout(t *testing.T) {
    fields := []struct {
        name  string
        start bitpack.BitPosition
        end   bitpack.BitPosition
    }{
        {"Version", VersionField.Start, VersionField.End},
        {"Priority", PriorityField.Start, PriorityField.End},
        {"Encrypted", EncryptedFlag.Position, EncryptedFlag.Position},
        {"Sequence", SequenceField.Start, SequenceField.End},
        {"Reserved", ReservedField.Start, ReservedField.End},
        {"LastFlag", LastFlag.Position, LastFlag.Position},
    }

    const totalBits = 32
    coverage := make([]bool, totalBits)

    // Проверяем каждое поле и отмечаем покрытые биты
    for _, f := range fields {
        t.Run(f.name, func(t *testing.T) {
            // Валидация границ поля
            if f.start > f.end || f.end >= totalBits {
                t.Errorf("поле %s выходит за пределы структуры: [%d:%d]", 
                    f.name, f.start, f.end)
            }

            // Проверка на пересечение с уже покрытыми битами
            for pos := f.start; pos <= f.end; pos++ {
                if coverage[pos] {
                    t.Errorf("бит %d уже занят другим полем (конфликт в поле %s)", 
                        pos, f.name)
                }
                coverage[pos] = true
            }
        })
    }

    // Проверка полного покрытия
    missing := 0
    for i, covered := range coverage {
        if !covered {
            t.Errorf("бит %d не назначен ни одному полю", i)
            missing++
        }
    }
    if missing > 0 {
        t.Fatalf("обнаружено %d непокрытых битов из %d", missing, totalBits)
    }

    t.Logf("Все %d бита корректно распределены между %d полями", 
        totalBits, len(fields))
}
```

## Конкурентность

**ВАЖНО:** Операции с полями не являются потоко безопасными. Необходима дополнительная синхронизация.


## Примеры использования

### Сетевой протокол

```go
type PacketHeader struct {
    version   bitpack.UIntBitField
    flags     bitpack.UIntBitField
    priority  bitpack.UIntBitField
    encrypted bitpack.BoolBitField
    sequence  bitpack.UIntBitField
}

func NewPacketHeader() PacketHeader {
    return PacketHeader{
        version:   bitpack.MustNewUIntBitField(0, 3, 15),
        flags:     bitpack.MustNewUIntBitField(4, 7, 15),
        priority:  bitpack.MustNewUIntBitField(8, 10, 7),
        encrypted: bitpack.MustNewBoolBitField(11),
        sequence:  bitpack.MustNewUIntBitField(12, 31, 1048575),
    }
}
```

### Метеорологическая станция

```go
var (
    ErrorFlag     = bitpack.MustNewBoolBitField(31)
    Reserved      = bitpack.MustNewUIntBitField(24, 30, 127)
    PressureDelta = bitpack.MustNewIntBitField(16, 23, -128, 127)
    Humidity      = bitpack.MustNewUIntBitField(8, 15, 100)
    Temperature   = bitpack.MustNewIntBitField(0, 7, -50, 100)
)

func EncodeWeatherData(temp, humidity int, pressure int, error bool) bitpack.Packed32 {
    var data bitpack.Packed32
    bitpack.SetIntFieldUncheckedAs[int8](data[:], Temperature, int8(temp))
    bitpack.SetUIntFieldUncheckedAs[uint8](data[:], Humidity, uint8(humidity))
    bitpack.SetIntFieldUncheckedAs[int8](data[:], PressureDelta, int8(pressure))
    bitpack.SetBoolFieldUnchecked(data[:], ErrorFlag, error)
    return data
}
```

### Игровой персонаж (компактное хранение)

```go
// 48 бит (6 байт) для хранения атрибутов персонажа
var (
    HealthField   = bitpack.MustNewUIntBitField(0, 9, 1000)   // 10 бит: 0-1000
    ManaField     = bitpack.MustNewUIntBitField(10, 19, 1000) // 10 бит: 0-1000
    LevelField    = bitpack.MustNewUIntBitField(20, 23, 10)   // 4 бит: 0-10
    StrengthField = bitpack.MustNewUIntBitField(24, 27, 10)   // 4 бит: 0-10
    HasHouseFlag  = bitpack.MustNewBoolBitField(28)           // 1 бит
    HasWeaponFlag = bitpack.MustNewBoolBitField(29)           // 1 бит
)

type Character struct {
    name   [32]byte            // имя персонажа
    attrs  bitpack.Packed48    // упакованные атрибуты
}

func (c *Character) SetHealth(health uint32) error {
    return bitpack.SetUIntFieldAs[uint32](c.attrs[:], HealthField, health)
}

func (c *Character) Health() uint32 {
    return bitpack.GetUIntFieldAs[uint32](c.attrs[:], HealthField)
}
```


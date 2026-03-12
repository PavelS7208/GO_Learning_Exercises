package task

import "fmt"

// Priority Приоритет задачи со значениями [0, 100].
type Priority int

const (
	priorityMin = 0
	priorityMax = 100
)

// Предустановленные уровни
const (
	LowPriority      Priority = 0
	NormalPriority   Priority = 25
	HighPriority     Priority = 50
	CriticalPriority Priority = 75
	UrgentPriority   Priority = 100
)

// NewPriority конструктор
func NewPriority(v int) (Priority, error) {
	if v < priorityMin || v > priorityMax {
		return 0, fmt.Errorf("priority %d out of range [%d, %d]", v, priorityMin, priorityMax)
	}
	return Priority(v), nil
}

// MustPriority специфичный конструктор для констант
// Паникуем если напутали в коде
func MustPriority(v int) Priority {
	p, err := NewPriority(v)
	if err != nil {
		panic(err)
	}
	return p
}

//--------------------  Task ---------------------------

// TaskID is a unique identifier for a task.
type TaskID int

type MetaData struct{}

// Task - unit of work с уникальным идентификатором и приоритетом
type Task struct {
	Identifier TaskID
	Priority   Priority
	metaData   *MetaData
}

// Реализуем интерфейс с методами Key() и Order() для хранения Task-ов в приоритетной очереди Планировщика (по Order порядку)

func (t Task) Key() int   { return int(t.Identifier) }
func (t Task) Order() int { return int(t.Priority) }

package scheduler

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/example/scheduler/task"
)

// ShardID unique identifier for a shard.
type ShardID uint64

// ShardInfo справочная структура (read-only)
type ShardInfo struct {
	Index int     // Текущая позиция Шарды в пуле Планировщика
	ID    ShardID // Глобальный идентификатор Шарды (не меняется)
}

// shardState
// Пока из статусов только то, что процесс перераспределения задач между шардами идет,
// но вдруг понадобится расширить статусы
type shardState struct {
	rebalancing atomic.Bool
}

// ----------------- shardEntry ---------------
// Полноценная шарда Планировщика/
// Хранит задачи удовлетворяющие минимальному контракту для хранения в node (Order() - приоритет и Key() - ID)
type shardEntry struct {
	id    ShardID
	node  *node[task.Task]
	state *shardState
}

// ============== Scheduler Pool =========================
//  Планировщик
//  Динамически изменяемый набор (pool) шард shardEntry
// (mu.Lock):  AddShard, RemoveShard, ForceRemoveShard
// (mu.RLock): остальное
//
// ShardIDs уникальный ID, номер в списке - это другое

type Pool struct {
	mu      sync.RWMutex
	shards  []*shardEntry   // dense slice, order may change on remove
	index   map[ShardID]int // ShardID → current position in shards
	nextID  atomic.Uint64   // monotonic ID generator
	counter atomic.Uint64   // round-robin указатель для AddTask
}

// New Создаем Планировщик. Минимум 1 нода
func New(n int) (*Pool, error) {
	if n < 1 {
		return nil, fmt.Errorf("invalid number of shards %d: must be at least 1", n)
	}
	p := &Pool{
		index: make(map[ShardID]int, n),
	}
	for range n {
		p.appendShard()
	}
	return p, nil
}

// ------  Конфигурация -----------------------------------------

// ShardCount Кол-во шардов, так как может меняться - то в моменте лочим
func (p *Pool) ShardCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.shards)
}

// Shards Информация о шардах в моменте
func (p *Pool) Shards() []ShardInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]ShardInfo, len(p.shards))
	for i, e := range p.shards {
		out[i] = ShardInfo{Index: i, ID: e.id}
	}
	return out
}

// IsShardEmpty есть ли задачи в шарде планировщика - в моменте запрос
func (p *Pool) IsShardEmpty(id ShardID) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	e, err := p.lookupLocked(id)
	if err != nil {
		return false, err
	}
	return e.node.Count() == 0, nil
}

// ------------ Изменения Конфигурации -----------------

// AddShard
// AddShard добавляет шарду в пул Планировщика и асинхронно запускает
// ребалансировку: крадёт половину низкоприоритетных задач у самой загруженной шарды.
// Ребалансировка выполняется в фоне — возврат ShardInfo не ждёт её завершения.
// Returns ShardInfo новой шарды.
func (p *Pool) AddShard() ShardInfo {
	p.mu.Lock()
	e := p.appendShard()
	info := ShardInfo{Index: p.index[e.id], ID: e.id}
	p.mu.Unlock()

	e.state.rebalancing.Store(true)
	go p.rebalance(e.id) // как в GetTaskFromShard
	return info
}

// RemoveShard удаляем первую по списку пустую шарду из пула Планировщика
// Returns false - если нет пустых шард или последняя в конфигурации осталась
func (p *Pool) RemoveShard() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.shards) <= 1 {
		return false
	}
	for i, e := range p.shards {
		if e.node.Count() == 0 {
			p.removeAtLocked(i)
			return true
		}
	}
	return false
}

// ForceRemoveShard Жестко удаляем шарду с ID заданным.
// Задачи раскидываются на остальные работающие шарды перед удалением
// Returns an error - если ID кривой или одна осталась в пулле шарда
func (p *Pool) ForceRemoveShard(id ShardID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.shards) <= 1 {
		return fmt.Errorf("cannot remove the last shard")
	}
	e, err := p.lookupLocked(id)
	if err != nil {
		return err
	}

	// Раскидываем все задачи по другим шардам
	var cursor int
	for {
		t, ok := e.node.get()
		if !ok {
			break
		}
		// в каждую шарду (кроме себя) кидаем по задачке. cursor обеспечивает замыкание круговое (алгоритм round-robin)
		target := p.nextRoundRobinExcludingLocked(id, &cursor)
		target.node.add(t)
	}
	p.removeAtLocked(p.index[id])
	return nil
}

// --------------- Task routing ------

// AddTask Добавляем задачу в Планировщик. Шарду выбираем via round-robin (по остатку от деления от счетчика)
// Panics if the pool has no shards.
func (p *Pool) AddTask(t task.Task) {
	p.mu.RLock()
	n := uint64(len(p.shards))
	if n == 0 { // Хоть такого и быть не может, но проверка
		p.mu.RUnlock()
		panic("scheduler: AddTask called on empty pool")
	}
	idx := int(p.counter.Add(1)-1) % int(n)
	e := p.shards[idx]
	p.mu.RUnlock()
	e.node.add(t)
}

// AddTaskToShard Добавляем задачу в конкретную шарду по ID.
func (p *Pool) AddTaskToShard(id ShardID, t task.Task) error {
	p.mu.RLock()
	e, err := p.lookupLocked(id)
	p.mu.RUnlock()
	if err != nil {
		return err
	}
	e.node.add(t)
	return nil
}

// GetTaskFromShard Получаем самую приоритетную задачу с заданной шарды Планировщика.
// Параметр ok is false если шарда пустая (без задач)
// Параметр err is non-nil когда ID шарды кривой
// Если забрали последнюю задачу - запускаем перебалансировку задач в фоне
func (p *Pool) GetTaskFromShard(id ShardID) (t task.Task, ok bool, err error) {
	p.mu.RLock()
	e, err := p.lookupLocked(id)
	p.mu.RUnlock()
	if err != nil {
		return task.Task{}, false, err
	}

	t, ok = e.node.get()
	if !ok {
		if e.state.rebalancing.CompareAndSwap(false, true) {
			go p.rebalance(id)
		}
		return task.Task{}, false, nil
	}
	return t, true, nil
}

// ------- Задача может динамически менять приоритет после попадания в планировщик ----------------

// ChangeTaskPriority Ищем по всем шардам
// Returns false if the task is not found.
func (p *Pool) ChangeTaskPriority(id task.TaskID, pr task.Priority) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	updated := task.Task{Identifier: id, Priority: pr}
	for _, e := range p.shards {
		if e.node.change(int(id), updated) {
			return true
		}
	}
	return false
}

// ChangeTaskPriorityInShard Ищем в конкретной шарде задачу на изменение
func (p *Pool) ChangeTaskPriorityInShard(id ShardID, taskID task.TaskID, pr task.Priority) (bool, error) {
	p.mu.RLock()
	e, err := p.lookupLocked(id)
	p.mu.RUnlock()
	if err != nil {
		return false, err
	}
	return e.node.change(int(taskID), task.Task{Identifier: taskID, Priority: pr}), nil
}

// ----- Сервисные функции -----------------------------

// appendShard creates a new shardEntry, appends it, and updates the index.
// Must be called with mu.Lock held (or during New before the pool is shared).
func (p *Pool) appendShard() *shardEntry {
	id := ShardID(p.nextID.Add(1) - 1)
	e := &shardEntry{
		id:    id,
		node:  newNode[task.Task](),
		state: &shardState{},
	}
	p.index[id] = len(p.shards)
	p.shards = append(p.shards, e)
	return e
}

// removeAtLocked Удаляет пустую шарду по индексу
// Updates the index map accordingly.
// Must be called with mu.Lock held.
func (p *Pool) removeAtLocked(i int) {
	last := len(p.shards) - 1
	removed := p.shards[i]
	delete(p.index, removed.id)

	if i != last {
		p.shards[i] = p.shards[last]
		p.index[p.shards[i].id] = i
	}
	p.shards[last] = nil
	p.shards = p.shards[:last]
}

// lookupLocked
// Must be called with at least mu.RLock held.
func (p *Pool) lookupLocked(id ShardID) (*shardEntry, error) {
	idx, ok := p.index[id]
	if !ok {
		return nil, fmt.Errorf("shard %d not found", id)
	}
	return p.shards[idx], nil
}

// nextRoundRobinExcludingLocked возвращает следующую шарду, проверяя что она не та что передали
// cursor нужен для кругового обхода.
// Panics если какая-то серьезная ошибка конфигурации пула. Не должно быь если все верно сконфигурировано.
// Must be called with mu.Lock held.
func (p *Pool) nextRoundRobinExcludingLocked(exclude ShardID, cursor *int) *shardEntry {
	n := len(p.shards)
	for range n {
		e := p.shards[*cursor%n]
		*cursor++
		if e.id != exclude {
			return e
		}
	}
	panic("scheduler: nextRoundRobinExcludingLocked found no eligible shard")
}

// rebalance()
// Используем алгоритм stealing tasks
// Воруем половину низкоприоритетных задач у самой загруженной шарды
// и передаем их указанной шарде
func (p *Pool) rebalance(id ShardID) {
	defer func() {
		p.mu.RLock()
		e, err := p.lookupLocked(id)
		p.mu.RUnlock()
		if err == nil {
			e.state.rebalancing.Store(false)
		}
	}()

	// п.2: single RLock/RUnlock scope — no manual unlocks inside.
	p.mu.RLock()
	e, err := p.lookupLocked(id)
	if err != nil {
		p.mu.RUnlock()
		return
	}
	if len(p.shards) < 2 {
		p.mu.RUnlock()
		return
	}
	victim := p.mostLoadedShardLocked(id)
	p.mu.RUnlock()

	if victim == nil {
		return
	}

	stolen := victim.node.stealHalf()
	for _, t := range stolen {
		e.node.add(t)
	}
}

// mostLoadedShardLocked возвращает самую загруженную шарду по Count(),
// исключая ту что передали при поиске
// Используем atomic counters — no node locks acquired.
// Must be called with at least mu.RLock held.
func (p *Pool) mostLoadedShardLocked(exclude ShardID) *shardEntry {
	var best *shardEntry
	bestCount := 0
	for _, e := range p.shards {
		if e.id == exclude {
			continue
		}
		if c := e.node.Count(); c > bestCount {
			bestCount = c
			best = e
		}
	}
	return best
}

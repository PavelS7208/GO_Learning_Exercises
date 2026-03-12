package scheduler_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/scheduler/scheduler"
	"github.com/example/scheduler/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTask(id task.TaskID, p task.Priority) task.Task {
	return task.Task{Identifier: id, Priority: p}
}

// mustNew is a test helper that wraps scheduler.New and fails the test on error.
func mustNew(t *testing.T, n int) *scheduler.Pool {
	t.Helper()
	p, err := scheduler.New(n)
	require.NoError(t, err)
	return p
}

// shardID returns the ShardID at position idx in the pool's Shards() snapshot.
func shardID(p *scheduler.Pool, idx int) scheduler.ShardID {
	return p.Shards()[idx].ID
}

// -----------------------------------------------------------------------
// New — validation
// -----------------------------------------------------------------------

func TestNewZeroShards(t *testing.T) {
	_, err := scheduler.New(0)
	assert.Error(t, err)
}

func TestNewNegativeShards(t *testing.T) {
	_, err := scheduler.New(-1)
	assert.Error(t, err)
}

func TestNewOneShardIsValid(t *testing.T) {
	p, err := scheduler.New(1)
	require.NoError(t, err)
	assert.Equal(t, 1, p.ShardCount())
}

// -----------------------------------------------------------------------
// ShardCount / Shards / IsShardEmpty
// -----------------------------------------------------------------------

func TestShardCount(t *testing.T) {
	p := mustNew(t, 3)
	assert.Equal(t, 3, p.ShardCount())
}

func TestShardsSnapshot(t *testing.T) {
	p := mustNew(t, 3)
	infos := p.Shards()
	require.Len(t, infos, 3)
	seen := make(map[scheduler.ShardID]bool)
	for _, info := range infos {
		assert.False(t, seen[info.ID], "duplicate ShardID in snapshot")
		seen[info.ID] = true
		assert.GreaterOrEqual(t, info.Index, 0)
		assert.Less(t, info.Index, 3)
	}
}

func TestIsShardEmpty(t *testing.T) {
	p := mustNew(t, 2)
	id0 := shardID(p, 0)

	empty, err := p.IsShardEmpty(id0)
	require.NoError(t, err)
	assert.True(t, empty)

	require.NoError(t, p.AddTaskToShard(id0, newTask(1, task.NormalPriority)))

	empty, err = p.IsShardEmpty(id0)
	require.NoError(t, err)
	assert.False(t, empty)
}

func TestIsShardEmptyUnknownID(t *testing.T) {
	p := mustNew(t, 2)
	_, err := p.IsShardEmpty(scheduler.ShardID(999))
	assert.Error(t, err)
}

// -----------------------------------------------------------------------
// AddTask — round-robin distribution
// -----------------------------------------------------------------------

func TestAddTaskRoundRobin(t *testing.T) {
	p := mustNew(t, 2)
	for i := task.TaskID(1); i <= 4; i++ {
		p.AddTask(newTask(i, task.NormalPriority))
	}
	id0 := shardID(p, 0)
	id1 := shardID(p, 1)

	t1, ok, err := p.GetTaskFromShard(id0)
	require.NoError(t, err)
	require.True(t, ok)
	_ = t1

	t2, ok, err := p.GetTaskFromShard(id1)
	require.NoError(t, err)
	require.True(t, ok)
	_ = t2
}

// -----------------------------------------------------------------------
// AddTaskToShard
// -----------------------------------------------------------------------

func TestAddTaskToShardUnknownID(t *testing.T) {
	p := mustNew(t, 3)
	err := p.AddTaskToShard(scheduler.ShardID(999), newTask(1, task.NormalPriority))
	assert.Error(t, err)
}

func TestAddTaskToShardValid(t *testing.T) {
	p := mustNew(t, 3)
	id := shardID(p, 2)

	require.NoError(t, p.AddTaskToShard(id, newTask(1, task.NormalPriority)))

	got, ok, err := p.GetTaskFromShard(id)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, task.TaskID(1), got.Identifier)
}

// -----------------------------------------------------------------------
// GetTaskFromShard
// -----------------------------------------------------------------------

func TestGetTaskFromShardUnknownID(t *testing.T) {
	p := mustNew(t, 2)
	_, ok, err := p.GetTaskFromShard(scheduler.ShardID(999))
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestGetTaskFromShardEmptyShard(t *testing.T) {
	p := mustNew(t, 2)
	_, ok, err := p.GetTaskFromShard(shardID(p, 0))
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestGetTaskFromShardReturnsNilOnlyWhenTrulyEmpty(t *testing.T) {
	p := mustNew(t, 1)
	id := shardID(p, 0)
	p.AddTask(newTask(1, task.NormalPriority))

	first, ok, err := p.GetTaskFromShard(id)
	require.NoError(t, err)
	require.True(t, ok)
	_ = first

	_, ok, err = p.GetTaskFromShard(id)
	require.NoError(t, err)
	assert.False(t, ok)
}

// -----------------------------------------------------------------------
// ChangeTaskPriority
// -----------------------------------------------------------------------

func TestChangeTaskPriorityFound(t *testing.T) {
	p := mustNew(t, 3)
	ids := make([]scheduler.ShardID, 3)
	for i := range 3 {
		ids[i] = shardID(p, i)
	}
	require.NoError(t, p.AddTaskToShard(ids[0], newTask(1, task.LowPriority)))
	require.NoError(t, p.AddTaskToShard(ids[1], newTask(2, task.NormalPriority)))
	require.NoError(t, p.AddTaskToShard(ids[2], newTask(3, task.HighPriority)))

	assert.True(t, p.ChangeTaskPriority(2, task.UrgentPriority))

	got, ok, err := p.GetTaskFromShard(ids[1])
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, task.Task{Identifier: 2, Priority: task.UrgentPriority}, got)
}

func TestChangeTaskPriorityNotFound(t *testing.T) {
	p := mustNew(t, 3)
	assert.False(t, p.ChangeTaskPriority(999, task.HighPriority))
}

func TestChangeTaskPriorityAlreadyConsumed(t *testing.T) {
	p := mustNew(t, 1)
	id := shardID(p, 0)
	p.AddTask(newTask(1, task.NormalPriority))

	got, ok, err := p.GetTaskFromShard(id)
	require.NoError(t, err)
	require.True(t, ok)
	_ = got

	assert.False(t, p.ChangeTaskPriority(1, task.UrgentPriority))
}

func TestChangeTaskPriorityInShardUnknownID(t *testing.T) {
	p := mustNew(t, 2)
	ok, err := p.ChangeTaskPriorityInShard(scheduler.ShardID(99), 1, task.HighPriority)
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestChangeTaskPriorityInShardNotFound(t *testing.T) {
	p := mustNew(t, 2)
	ok, err := p.ChangeTaskPriorityInShard(shardID(p, 0), 999, task.HighPriority)
	assert.NoError(t, err)
	assert.False(t, ok)
}

// -----------------------------------------------------------------------
// Priority ordering
// -----------------------------------------------------------------------

func TestSingleShard(t *testing.T) {
	p := mustNew(t, 1)
	id := shardID(p, 0)

	p.AddTask(newTask(1, task.MustPriority(10)))
	p.AddTask(newTask(2, task.MustPriority(20)))
	p.AddTask(newTask(3, task.MustPriority(30)))
	p.AddTask(newTask(4, task.MustPriority(40)))
	p.AddTask(newTask(5, task.MustPriority(50)))

	got, ok, err := p.GetTaskFromShard(id)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, task.TaskID(5), got.Identifier)

	got, ok, err = p.GetTaskFromShard(id)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, task.TaskID(4), got.Identifier)

	p.ChangeTaskPriority(1, task.UrgentPriority)

	got, ok, err = p.GetTaskFromShard(id)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, task.Task{Identifier: 1, Priority: task.UrgentPriority}, got)
}

func TestPredefinedPriorityLevels(t *testing.T) {
	p := mustNew(t, 1)
	id := shardID(p, 0)

	p.AddTask(newTask(1, task.LowPriority))
	p.AddTask(newTask(2, task.NormalPriority))
	p.AddTask(newTask(3, task.HighPriority))
	p.AddTask(newTask(4, task.CriticalPriority))
	p.AddTask(newTask(5, task.UrgentPriority))

	got, ok, err := p.GetTaskFromShard(id)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, task.TaskID(5), got.Identifier)
}

// -----------------------------------------------------------------------
// Rebalance — steals from the most loaded shard
// -----------------------------------------------------------------------

func TestRebalanceStealsFromMostLoaded(t *testing.T) {
	p := mustNew(t, 3)
	id0 := shardID(p, 0)
	id1 := shardID(p, 1)

	// Shard 1 gets 8 tasks (most loaded).
	for i := task.TaskID(1); i <= 8; i++ {
		require.NoError(t, p.AddTaskToShard(id1, newTask(i, task.MustPriority(10))))
	}

	// Reading from empty shard 0 triggers rebalance.
	_, ok, err := p.GetTaskFromShard(id0)
	require.NoError(t, err)
	assert.False(t, ok)

	time.Sleep(20 * time.Millisecond)

	_, ok, err = p.GetTaskFromShard(id0)
	require.NoError(t, err)
	assert.True(t, ok, "rebalance should have stolen tasks from most loaded shard into shard 0")
}

func TestRebalanceNotTriggeredTwice(t *testing.T) {
	p := mustNew(t, 2)
	id0 := shardID(p, 0)

	for i := task.TaskID(1); i <= 8; i++ {
		p.AddTask(newTask(i, task.MustPriority(10)))
	}

	for {
		_, ok, err := p.GetTaskFromShard(id0)
		require.NoError(t, err)
		if !ok {
			break
		}
	}

	_, ok, err := p.GetTaskFromShard(id0)
	require.NoError(t, err)
	assert.False(t, ok)

	time.Sleep(20 * time.Millisecond)
}

func TestRebalanceNoVictimBelowMinStealSize(t *testing.T) {
	p := mustNew(t, 2)
	id0 := shardID(p, 0)
	id1 := shardID(p, 1)

	require.NoError(t, p.AddTaskToShard(id1, newTask(1, task.NormalPriority)))
	require.NoError(t, p.AddTaskToShard(id1, newTask(2, task.NormalPriority)))

	_, ok, err := p.GetTaskFromShard(id0)
	require.NoError(t, err)
	assert.False(t, ok)

	time.Sleep(20 * time.Millisecond)

	_, ok, err = p.GetTaskFromShard(id0)
	require.NoError(t, err)
	assert.False(t, ok, "rebalance must not steal when victim is below minStealSize")

	_, ok, err = p.GetTaskFromShard(id1)
	require.NoError(t, err)
	assert.True(t, ok, "shard 1 tasks must be intact")
}

// -----------------------------------------------------------------------
// AddShard
// -----------------------------------------------------------------------

func TestAddShardIncreasesCount(t *testing.T) {
	p := mustNew(t, 2)
	info := p.AddShard()
	assert.Equal(t, 3, p.ShardCount())
	assert.NotEqual(t, scheduler.ShardID(0), info.ID) // ID must be fresh (not reused)
}

func TestAddShardIDIsStable(t *testing.T) {
	p := mustNew(t, 2)
	info := p.AddShard()
	found := false
	for _, s := range p.Shards() {
		if s.ID == info.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "new shard ID must appear in Shards() snapshot")
}

func TestAddShardStealsFromMostLoaded(t *testing.T) {
	p := mustNew(t, 1)
	id0 := shardID(p, 0)

	for i := task.TaskID(1); i <= 8; i++ {
		require.NoError(t, p.AddTaskToShard(id0, newTask(i, task.MustPriority(10))))
	}

	info := p.AddShard() // rebalance запускается асинхронно

	time.Sleep(20 * time.Millisecond)

	_, ok, err := p.GetTaskFromShard(info.ID)
	require.NoError(t, err)
	assert.True(t, ok, "new shard should have received stolen tasks after async rebalance")
}

// -----------------------------------------------------------------------
// RemoveShard
// -----------------------------------------------------------------------

func TestRemoveShardEmptyExists(t *testing.T) {
	p := mustNew(t, 3)
	id0 := shardID(p, 0)

	require.NoError(t, p.AddTaskToShard(shardID(p, 1), newTask(1, task.NormalPriority)))
	require.NoError(t, p.AddTaskToShard(shardID(p, 2), newTask(2, task.NormalPriority)))

	ok := p.RemoveShard()
	assert.True(t, ok)
	assert.Equal(t, 2, p.ShardCount())

	_, err := p.IsShardEmpty(id0)
	assert.Error(t, err)
}

func TestRemoveShardNoEmptyShards(t *testing.T) {
	p := mustNew(t, 2)
	require.NoError(t, p.AddTaskToShard(shardID(p, 0), newTask(1, task.NormalPriority)))
	require.NoError(t, p.AddTaskToShard(shardID(p, 1), newTask(2, task.NormalPriority)))

	ok := p.RemoveShard()
	assert.False(t, ok)
	assert.Equal(t, 2, p.ShardCount())
}

func TestRemoveShardLastShardProtected(t *testing.T) {
	p := mustNew(t, 1)
	ok := p.RemoveShard()
	assert.False(t, ok)
	assert.Equal(t, 1, p.ShardCount())
}

// -----------------------------------------------------------------------
// ForceRemoveShard
// -----------------------------------------------------------------------

func TestForceRemoveShardRedistributesTasks(t *testing.T) {
	p := mustNew(t, 3)
	id0 := shardID(p, 0)
	id1 := shardID(p, 1)
	id2 := shardID(p, 2)

	for i := task.TaskID(1); i <= 4; i++ {
		require.NoError(t, p.AddTaskToShard(id0, newTask(i, task.MustPriority(10))))
	}

	require.NoError(t, p.ForceRemoveShard(id0))
	assert.Equal(t, 2, p.ShardCount())

	_, err := p.IsShardEmpty(id0)
	assert.Error(t, err)

	total := 0
	for _, id := range []scheduler.ShardID{id1, id2} {
		for {
			_, ok, err := p.GetTaskFromShard(id)
			require.NoError(t, err)
			if !ok {
				break
			}
			total++
		}
	}
	assert.Equal(t, 4, total, "all tasks must be redistributed to remaining shards")
}

func TestForceRemoveShardEmptyShard(t *testing.T) {
	p := mustNew(t, 2)
	id0 := shardID(p, 0)

	// shard 0 is empty — ForceRemoveShard must still succeed
	require.NoError(t, p.ForceRemoveShard(id0))
	assert.Equal(t, 1, p.ShardCount())

	_, err := p.IsShardEmpty(id0)
	assert.Error(t, err)
}

func TestForceRemoveShardLastShardProtected(t *testing.T) {
	p := mustNew(t, 1)
	err := p.ForceRemoveShard(shardID(p, 0))
	assert.Error(t, err)
}

func TestForceRemoveShardUnknownID(t *testing.T) {
	p := mustNew(t, 2)
	err := p.ForceRemoveShard(scheduler.ShardID(999))
	assert.Error(t, err)
}

// ShardIDs remain stable after ForceRemoveShard (swap-with-last).
func TestShardIDsStableAfterRemove(t *testing.T) {
	p := mustNew(t, 3)
	before := p.Shards()

	mid := before[1].ID
	require.NoError(t, p.ForceRemoveShard(mid))

	after := p.Shards()
	require.Len(t, after, 2)

	for _, info := range after {
		assert.NotEqual(t, mid, info.ID, "removed shard ID must not appear after removal")
	}
	remaining := map[scheduler.ShardID]bool{before[0].ID: false, before[2].ID: false}
	for _, info := range after {
		if _, ok := remaining[info.ID]; ok {
			remaining[info.ID] = true
		}
	}
	for id, found := range remaining {
		assert.True(t, found, "shard %d should still be present", id)
	}
}

// -----------------------------------------------------------------------
// Concurrent stress
// -----------------------------------------------------------------------

func TestConcurrentAccess(t *testing.T) {
	const numShards = 4
	const numTasks = 1000

	p := mustNew(t, numShards)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numTasks; i++ {
			p.AddTask(newTask(task.TaskID(i), task.MustPriority(i%100)))
		}
	}()
	wg.Wait()

	infos := p.Shards()

	var (
		got int
		mu  sync.Mutex
	)

	for _, info := range infos {
		id := info.ID
		wg.Add(1)
		go func() {
			defer wg.Done()
			idle := 0
			for idle < 5 {
				_, ok, err := p.GetTaskFromShard(id)
				if err != nil {
					return
				}
				if !ok {
					idle++
					time.Sleep(5 * time.Millisecond)
					continue
				}
				idle = 0
				mu.Lock()
				got++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	time.Sleep(50 * time.Millisecond)

	for _, info := range p.Shards() {
		for {
			_, ok, err := p.GetTaskFromShard(info.ID)
			if err != nil || !ok {
				break
			}
			mu.Lock()
			got++
			mu.Unlock()
		}
	}

	assert.Equal(t, numTasks, got, "all tasks must be retrieved exactly once")
}

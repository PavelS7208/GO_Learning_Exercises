package main

import (
	"fmt"

	"github.com/example/scheduler/scheduler"
	"github.com/example/scheduler/task"
)

func main() {
	p, _ := scheduler.New(4)

	fmt.Printf("Создали pool: %d shards\n", p.ShardCount())
	for _, info := range p.Shards() {
		fmt.Printf("  index=%-2d id=%d\n", info.Index, info.ID)
	}

	tasks := []task.Task{
		{Identifier: 1, Priority: task.LowPriority},
		{Identifier: 2, Priority: task.NormalPriority},
		{Identifier: 3, Priority: task.HighPriority},
		{Identifier: 4, Priority: task.CriticalPriority},
		{Identifier: 5, Priority: task.UrgentPriority},
	}
	for _, t := range tasks {
		p.AddTask(t)
	}

	p.ChangeTaskPriority(task.TaskID(1), task.MustPriority(90))

	// Add a new shard — it will steal from the most loaded.
	info := p.AddShard()
	fmt.Printf("\nДобавили shard: index=%-2d id=%d\n", info.Index, info.ID)
	fmt.Printf("Pool сейчас имеет %d shards\n", p.ShardCount())

	//
	fmt.Println("\nВсе шарды содержали:")
	for _, s := range p.Shards() {
		for {
			t, ok, err := p.GetTaskFromShard(s.ID)
			if err != nil || !ok {
				break
			}
			fmt.Printf("  shard id=%-2d  task id=%-3d priority=%d\n", s.ID, t.Identifier, t.Priority)
		}
	}

	// Remove empty shards.
	removed := p.RemoveShard()
	fmt.Printf("\nRemoveShard (empty): %v — pool сейчас имеет %d shards\n", removed, p.ShardCount())
	removed = p.RemoveShard()
	fmt.Printf("\nRemoveShard (empty): %v — pool сейчас имеет %d shards\n", removed, p.ShardCount())
}

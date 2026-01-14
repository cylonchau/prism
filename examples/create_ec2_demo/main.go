// Example demonstrates the executor layer usage with LocalStack.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/cylonchau/prism/pkg/dao"
	"github.com/cylonchau/prism/pkg/executor"
	"github.com/cylonchau/prism/pkg/executor/lock"
	"github.com/cylonchau/prism/pkg/executor/terraform"
	"github.com/cylonchau/prism/pkg/executor/ws"
)

func main() {
	// LocalStack terraform config path
	workDir := "./examples/localstack"

	// Check if localstack directory exists
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		log.Fatalf("LocalStack config not found at %s", workDir)
	}

	// 1. Setup database
	db, err := gorm.Open(sqlite.Open("create_ec2_demo.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	fmt.Println("✓ Database connected")

	// 2. Create components
	locker := lock.NewMemoryLocker(nil)
	taskDAO := dao.NewExecutionTaskDAO(db)
	hub := ws.NewHub()
	fmt.Println("✓ Components initialized")

	// 3. Create Terraform executor
	config := &terraform.Config{
		BinaryPath: "terraform",
		BasePath:   "/opt/homebrew/bin/terraform",
		Timeout:    5 * time.Minute,
	}
	exec := terraform.New(config, locker, taskDAO, hub)
	fmt.Println("✓ Terraform executor created")
	fmt.Printf("✓ Using config: %s\n", workDir)

	ctx := context.Background()

	// 4. Execute terraform plan
	fmt.Println("\n========== TERRAFORM PLAN ==========")
	planTaskID := fmt.Sprintf("plan-%d", time.Now().UnixNano())
	planReq := &executor.ExecuteRequest{
		TaskID:     planTaskID,
		ResourceID: 1,
		Action:     executor.ActionPlan,
		WorkDir:    workDir,
	}

	result, err := exec.Execute(ctx, planReq)
	printResult("Plan", result, err)

	// Show task record
	task, _ := taskDAO.Get(planTaskID)
	if task != nil {
		fmt.Printf("\n--- Task Record ---\n")
		fmt.Printf("TaskID:     %s\n", task.TaskID)
		fmt.Printf("Status:     %d (2=success, 3=failed)\n", task.Status)
		fmt.Printf("Duration:   %dms\n", task.Duration)
		if task.Error != "" {
			fmt.Printf("Error:      %s\n", task.Error)
		}
		fmt.Printf("\n--- Task Output (first 2000 chars) ---\n")
		output := task.Output
		if len(output) > 2000 {
			output = output[:2000] + "...(truncated)"
		}
		fmt.Println(output)
	}

	// 5. Test lock
	fmt.Println("\n========== LOCK TEST ==========")
	fmt.Printf("IsLocked before: %v\n", locker.IsLocked(1))
	locker.Acquire(ctx, 1, "test")
	fmt.Printf("IsLocked after acquire: %v\n", locker.IsLocked(1))
	locker.Release(1)
	fmt.Printf("IsLocked after release: %v\n", locker.IsLocked(1))

	// 6. List tasks
	fmt.Println("\n========== TASK HISTORY ==========")
	tasks, _ := taskDAO.ListByResource(1)
	for _, t := range tasks {
		fmt.Printf("- %s: action=%s, status=%d, duration=%dms\n",
			t.TaskID[:20]+"...", t.Action, t.Status, t.Duration)
	}

	// Database location
	fmt.Println("\n========== DATABASE ==========")
	fmt.Println("Database file: ./create_ec2_demo.db")
	fmt.Println("To view: sqlite3 create_ec2_demo.db 'SELECT * FROM execution_task;'")

	fmt.Println("\n✓ Demo completed")
}

func printResult(action string, result *executor.ExecuteResult, err error) {
	if err != nil {
		fmt.Printf("✗ %s failed: %v\n", action, err)
	} else {
		fmt.Printf("✓ %s completed: status=%s, duration=%dms\n",
			action, result.Status, result.Duration)
	}
}

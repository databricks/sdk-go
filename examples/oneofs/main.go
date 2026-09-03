//go:build examples

package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/databricks/sdk-go/jobs/v2"
)

func main() {
	// Oneof fields use a wrapper for each variant. Choose a variant by assigning
	// its wrapper to the interface field.
	task := jobs.TaskSettings{
		TaskKey: new("hello_world"),
		Task: &jobs.TaskSettings_Task_SparkPythonTask{
			SparkPythonTask: jobs.SparkPythonTask{
				PythonFile: new("/Workspace/Users/user@example.com/hello.py"),
				Source:     jobs.Source_Workspace,
			},
		},
	}

	if err := describeTask(task); err != nil {
		log.Fatal(err)
	}
}

func describeTask(task jobs.TaskSettings) error {
	// A type switch identifies the configured variant and provides typed access
	// to its fields.
	switch configuredTask := task.Task.(type) {
	case *jobs.TaskSettings_Task_NotebookTask:
		fmt.Printf("Notebook: %s\n", *configuredTask.NotebookTask.NotebookPath)
	case *jobs.TaskSettings_Task_SparkPythonTask:
		fmt.Printf("Python file: %s\n", *configuredTask.SparkPythonTask.PythonFile)
	case nil:
		return errors.New("task is missing its configuration")
	default:
		return fmt.Errorf("unsupported task configuration %T", configuredTask)
	}
	return nil
}

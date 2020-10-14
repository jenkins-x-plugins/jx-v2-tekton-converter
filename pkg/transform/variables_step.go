package transform

import (
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

// CreateVariablesTask creates a step to setup variables
func CreateVariablesTask(jxImage string) *v1alpha1.Task {
	s := v1alpha1.Step{}
	s.Command = []string{"jx"}
	s.Args = []string{"gitops", "variables"}
	s.Name = "jx-variables"
	s.Image = jxImage

	task := &v1alpha1.Task{}
	task.Spec.Steps = append(task.Spec.Steps, s)
	return task
}

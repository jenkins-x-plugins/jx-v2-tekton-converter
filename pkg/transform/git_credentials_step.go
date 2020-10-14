package transform

import (
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// CreateGitCredentialsTask creates a step to setup credentials
func CreateGitCredentialsTask(jxImage string) *v1alpha1.Task {
	s := v1alpha1.Step{}
	s.Command = []string{"jx"}
	s.Args = []string{"gitops", "git", "setup", "--namespace", "jx-git-operator"}
	s.Name = "git-setup"
	s.Image = jxImage

	s.VolumeMounts = []corev1.VolumeMount{
		{
			Name:      "workspace-volume",
			MountPath: "/home/jenkins",
		},
	}
	s.WorkingDir = "/workspace"
	task := &v1alpha1.Task{}
	task.Spec.Steps = append(task.Spec.Steps, s)
	return task
}

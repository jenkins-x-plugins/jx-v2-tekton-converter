package transform

import (
	"strings"

	"github.com/jenkins-x/jx/v2/pkg/config"
	"github.com/jenkins-x/jx/v2/pkg/jenkinsfile"
	"github.com/jenkins-x/jx/v2/pkg/tekton"
	"github.com/jenkins-x/jx/v2/pkg/tekton/syntax"
	"github.com/jenkins-x/jx/v2/pkg/util"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

var (
	initialGitStepNames = []string{"git-merge", "setup-jx-git-credentials", "git-credentials"}
)

// AddSetVersionSteps for a Jenkins X pipeline with a set version set of steps lets add the steps into the task
func AddSetVersionSteps(o *Options, results *tekton.CRDWrapper, projectConfig *config.ProjectConfig) (*tekton.CRDWrapper, error) {
	var sv *jenkinsfile.PipelineLifecycle
	pipelineConfig := projectConfig.PipelineConfig
	if pipelineConfig != nil {
		release := pipelineConfig.Pipelines.Release
		if release != nil {
			sv = release.SetVersion
		}

	}
	if sv == nil {
		command := "jx step next-version --use-git-tag-only --tag"
		if o.SemanticRelease {
			command = "jx step next-version --semantic-release --tag"
		}
		// lets create a default set version pipeline
		sv = &jenkinsfile.PipelineLifecycle{
			Steps: []*syntax.Step{
				{
					Command: command,
					Name:    "next-version",
					Comment: "tags git with the next version",
				},
			},
		}
	}

	tasks := results.Tasks()
	for _, task := range tasks {
		newSteps := addSteps(o, sv)

		// lets add the git steps first, then the set version, then the rest
		currentSteps := task.Spec.Steps
		task.Spec.Steps = nil
		for {
			if len(currentSteps) == 0 {
				break
			}
			s := currentSteps[0]
			if util.StringArrayIndex(initialGitStepNames, s.Name) < 0 {
				break
			}
			currentSteps = currentSteps[1:]
			task.Spec.Steps = append(task.Spec.Steps, s)
		}
		task.Spec.Steps = append(task.Spec.Steps, newSteps...)
		task.Spec.Steps = append(task.Spec.Steps, currentSteps...)
	}
	return results, nil
}

func addSteps(o *Options, sv *jenkinsfile.PipelineLifecycle) []v1alpha1.Step {
	newSteps := []v1alpha1.Step{}
	for _, ss := range sv.Steps {
		ssSteps := ss.Steps
		if len(ssSteps) == 0 {
			ssSteps = append(ssSteps, ss)
		}
		for _, s := range ssSteps {
			args := s.Arguments
			command := s.Command
			if command == "" {
				sh := s.Sh
				if sh != "" {
					command = "/bin/sh"
					sh = strings.ReplaceAll(sh, "\\$", "$")
					args = []string{"-c", sh}
				}
			}
			image := s.Image
			if image == "" {
				image = o.DefaultJXImage
			}
			dir := s.Dir
			if dir == "" {
				dir = "/workspace/source"
			}
			newSteps = append(newSteps, v1alpha1.Step{
				Container: corev1.Container{
					Name:       s.Name,
					Image:      image,
					Command:    []string{command},
					Args:       args,
					WorkingDir: dir,
					Env:        s.Env,
				},
				Script: "",
			})
		}
	}
	return newSteps
}

package transform

import (
	"github.com/jenkins-x-plugins/jx-v2-tekton-converter/pkg/assets"
	"github.com/pkg/errors"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

// CreateUsesGitClone creates a git clone step using the uses clause
func CreateUsesGitClone(release bool) *v1alpha1.Task {
	task := &v1alpha1.Task{}
	image := "uses:jenkins-x/jx3-pipeline-catalog/tasks/git-clone/git-clone-pr.yaml@versionStream"
	if release {
		image = "uses:jenkins-x/jx3-pipeline-catalog/tasks/git-clone/git-clone.yaml@versionStream"

	}
	task.Spec.Steps = []v1alpha1.Step{
		{
			Container: corev1.Container{
				Image: image,
			},
		},
	}
	return task
}

// CreateGitCloneTask creates the tekton catalog git clone
func CreateGitCloneTask(o *Options) (*v1alpha1.Task, error) {
	task := &v1alpha1.Task{}
	asset := "resources/git/git-clone.yaml"
	data, err := assets.Asset(asset)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load asset %s", asset)
	}

	err = yaml.Unmarshal(data, task)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse YAML %s", asset)
	}

	// lets set the default source dir
	sourceSubDir := o.CreateTaskOptions.SourceName
	if sourceSubDir != "" {
		in := task.Spec.Inputs
		if in != nil {
			for i, p := range in.Params {
				if p.Name == "subdirectory" {
					in.Params[i].Default = &v1alpha1.ArrayOrString{
						Type:      v1alpha1.ParamTypeString,
						StringVal: sourceSubDir,
					}
					break
				}
			}
		}
	}
	return task, nil
}

package transform

import (
	"github.com/jenkins-x-plugins/jx-v2-tekton-converter/pkg/assets"
	"github.com/pkg/errors"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"sigs.k8s.io/yaml"
)

// CreateGitCloneTask creates the tekton catalog git clone
func CreateGitCloneTask(o *Options) (*v1alpha1.Task, error) {
	asset := "resources/git/git-clone.yaml"
	data, err := assets.Asset(asset)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load asset %s", asset)
	}

	task := &v1alpha1.Task{}
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

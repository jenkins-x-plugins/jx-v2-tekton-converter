package transform

import (
	"fmt"
	"strings"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/jenkins-x/jx/v2/pkg/tekton"
	"github.com/jenkins-x/jx/v2/pkg/util"
	"github.com/pkg/errors"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TransformResources(o *Options, results *tekton.CRDWrapper) (*tekton.CRDWrapper, error) {
	pipeline := results.Pipeline()
	tasks := results.Tasks()
	pipelineRun := results.PipelineRun()

	envMap := map[string]string{}

	if o.JXImage == "" {
		// TODO use version stream
		o.JXImage = "gcr.io/jenkinsxio/jx-cli:latest"
	}

	for _, e := range replaceEnvVars {
		if e.Name != "" && e.Value != "" {
			envMap[e.Name] = e.Value
		}
	}

	var gitCloneTask *v1alpha1.Task
	var err error
	if o.UseCatalogGitClone {
		gitCloneTask, err = CreateGitCloneTask(o)
	} else {
		gitCloneTask, err = CreateSimpleGitCloneTask(o)
	}
	gitCredsTask := CreateGitCredentialsTask(o.JXImage)
	variablesTask := CreateVariablesTask(o.JXImage)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create git clone step")
	}
	if len(gitCloneTask.Spec.Steps) == 0 {
		return nil, errors.Errorf("git clone task has no steps")
	}

	removeEnvVarList := append([]string{}, removeEnvVars...)
	for _, p := range defaultTaskParams {
		removeEnvVarList = append(removeEnvVarList, p.Name)
	}

	// lets remove the pipeline resource input
	for _, task := range tasks {
		in := task.Spec.Inputs
		if in != nil {
			for i, r := range in.Resources {
				if r.Name == "workspace" && r.Type == v1alpha1.PipelineResourceTypeGit {
					in.Resources = append(in.Resources[0:i], in.Resources[i+1:]...)
					break
				}
			}
		}
		if in == nil {
			in = &v1alpha1.Inputs{}
			task.Spec.Inputs = in
		}

		// lets add any missing parameters from the git clone step
		task.Spec.Params = AddMissingParamSpecs(task.Spec.Params, gitCloneTask.Spec.Params)

		gitCloneSteps := []v1alpha1.Step{
			gitCloneTask.Spec.Steps[0],
			gitCredsTask.Spec.Steps[0],
		}
		if len(task.Spec.Steps) > 0 {
			s0 := task.Spec.Steps[0]
			if s0.WorkingDir == "/home/jenkins/go/src/REPLACE_ME_GIT_PROVIDER/REPLACE_ME_ORG/REPLACE_ME_APP_NAME" &&
				len(s0.Command) == 1 && s0.Command[0] == "" {
				// lets remove a bogus step
				task.Spec.Steps = task.Spec.Steps[1:]
			}
		}

		task.Spec.Steps = append(gitCloneSteps, task.Spec.Steps...)

		idx := -1
		for i, s := range task.Spec.Steps {
			if s.Name == "git-merge" {
				idx = i + 1
				break
			}
		}
		if idx < 0 {
			log.Logger().Warnf("failed to find git-merge step")
		} else {
			// lets insert the variables step at the right place
			rest := append([]v1beta1.Step{}, variablesTask.Spec.Steps[0])
			rest = append(rest, task.Spec.Steps[idx:]...)
			task.Spec.Steps = append(task.Spec.Steps[:idx], rest...)
		}

		// add the parameters to the task spec
		if task.Spec.StepTemplate == nil {
			task.Spec.StepTemplate = &corev1.Container{}
		}
		task.Spec.StepTemplate.Env = AddMissingParamEnv(task.Spec.StepTemplate.Env, defaultTaskParams)

		for i, s := range task.Spec.Steps {
			// filter out env vars
			var envs2 []corev1.EnvVar
			for _, e := range s.Env {
				if util.StringArrayIndex(removeEnvVarList, e.Name) < 0 {
					envs2 = append(envs2, e)
				}
			}
			task.Spec.Steps[i].Env = envs2
			for j, e := range s.Env {
				value := envMap[e.Name]
				if value != "" {
					task.Spec.Steps[i].Env[j].Value = value
					task.Spec.Steps[i].Env[j].ValueFrom = nil
				}
			}
			task.Spec.Steps[i].Args = transformStepArgs(s.Args)
		}

		if in == nil {
			in = &v1alpha1.Inputs{}
			task.Spec.Inputs = in
		}
		in.Params = RemoveOldParamSpecs(in.Params, "version")
		task.Spec.Params = AddMissingParamSpecs(task.Spec.Params, defaultTaskParams)
		task.Spec.Params = RemoveOldParamSpecs(task.Spec.Params, "version")
	}

	// lets zap old parameters
	pipeline.Spec.Params = AddMissingParamSpecs(pipeline.Spec.Params, defaultTaskParams)
	pipeline.Spec.Params = RemoveOldParamSpecs(pipeline.Spec.Params, "version")
	pipelineRun.Spec.Params = nil

	for i := range pipeline.Spec.Tasks {
		task := &pipeline.Spec.Tasks[i]
		task.Params = nil
		task.Params = AddMissingParams(task.Params, defaultTaskParams)
		task.Params = RemoveOldParams(task.Params, "version")
	}

	for i, r := range pipeline.Spec.Resources {
		if r.Type == v1alpha1.PipelineResourceTypeGit {
			pipeline.Spec.Resources = append(pipeline.Spec.Resources[0:i], pipeline.Spec.Resources[i+1:]...)
			break
		}
	}

	// lets remove old steps
	for _, t := range tasks {
		for i := range t.Spec.Steps {
			step := &t.Spec.Steps[i]
			args := step.Args
			if len(args) > 0 && args[0] == "jx step git credentials" {
				t.Spec.Steps = append(t.Spec.Steps[:i], t.Spec.Steps[i+1:]...)
				break
			}
		}
	}
	for _, t := range tasks {
		for i := range t.Spec.Steps {
			step := &t.Spec.Steps[i]
			if strings.HasPrefix(step.Image, "gcr.io/jenkinsxio/jx-cli:") {
				step.Image = o.JXImage
			}
		}
	}

	// convert tasks
	for _, t := range tasks {
		for i := range t.Spec.Steps {
			step := &t.Spec.Steps[i]
			args := step.Args
			if len(args) > 0 {
				if strings.Contains(args[0], "jx step helm release") {
					args[0] = "jx gitops helm release"
					step.Image = o.JXImage
					continue
				}
				if strings.HasPrefix(args[0], "/kaniko") {
					args[0] = "source .jx/variables.sh && cp /tekton/creds-secrets/tekton-container-registry-auth/.dockerconfigjson /kaniko/.docker/config.json && " + args[0]
					continue
				}

				if step.Command[0] == "jx" && len(args) > 2 && args[0] == "step" && args[1] == "git" && args[2] == "merge" {
					step.Args = []string{
						"step",
						"git",
						"merge",
						"--verbose",
						"--baseSHA",
						"$(params.PULL_BASE_SHA)",
						"--sha",
						"$(params.PULL_PULL_SHA)",
						"--baseBranch",
						"$(params.PULL_BASE_REF)",
					}
					continue
				}
				if strings.HasPrefix(args[0], "jx preview") {
					step.Command = []string{"/bin/bash", "-c"}
					args[0] = "source /workspace/source/.jx/variables.sh && jx preview create"
					continue
				}
				for _, cmdText := range avoidCDs {
					idx := strings.Index(args[0], cmdText)
					if idx > 0 {
						// strip the cd prefix
						args[0] = args[0][idx:]

						// lets add no poll for promote
						if strings.HasPrefix(args[0], "jx step changelog") {
							args[0] = "source /workspace/source/.jx/variables.sh && " + args[0]
							step.Command = []string{"/bin/bash", "-c"}
							continue
						}
						if strings.HasPrefix(args[0], "jx promote") {
							if !strings.Contains(args[0], " --no-poll") {
								args[0] = args[0] + " --no-poll"
							}
							args[0] = "source /workspace/source/.jx/variables.sh && " + args[0]
							step.Command = []string{"/bin/bash", "-c"}
							continue
						}
					}
				}
			}
		}

		// lets make sure that the git-merge step is before the set version (otherwise we lose the VERSION file)
		gitMergeIdx := -1
		nextVersionIdx := -1

		for i, step := range t.Spec.Steps {
			switch step.Name {
			case "git-merge":
				gitMergeIdx = i
			case "next-version":
				nextVersionIdx = i
			}
		}
		if gitMergeIdx >= 0 && nextVersionIdx >= 0 && gitMergeIdx > nextVersionIdx {
			// lets move the git merge to before the next-version
			steps := t.Spec.Steps
			s := append([]v1beta1.Step{}, steps[0:nextVersionIdx]...)
			s = append(s, steps[gitMergeIdx])
			s = append(s, steps[nextVersionIdx:gitMergeIdx]...)
			s = append(s, steps[gitMergeIdx+1:]...)
			t.Spec.Steps = s
		}
	}

	for _, t := range pipeline.Spec.Tasks {
		r := t.Resources
		if r != nil {
			for i, in := range r.Inputs {
				if in.Name == "workspace" {
					r.Inputs = append(r.Inputs[0:i], r.Inputs[i+1:]...)
					break
				}
			}
		}
	}

	// lets zap the first resource
	if len(pipelineRun.Spec.Resources) > 0 {
		pipelineRun.Spec.Resources = pipelineRun.Spec.Resources[1:]
	}

	// TODO until jx is on tekton 0.10.x lets hack around the validation change that lets scripts avoid shebangs
	scripts := []string{}
	for _, t := range tasks {
		script := ""
		if len(t.Spec.Steps) > 0 {
			script = t.Spec.Steps[0].Script
			t.Spec.Steps[0].Script = ""
		}
		scripts = append(scripts, script)
	}

	answer, err := tekton.NewCRDWrapper(pipeline, tasks, nil, nil, pipelineRun)
	if answer == nil || err != nil {
		return answer, err
	}

	// remove any template env vars
	tasks = answer.Tasks()
	for _, task := range tasks {
		if task.Spec.StepTemplate == nil {
			task.Spec.StepTemplate = &corev1.Container{}
		}
		if task.Spec.StepTemplate.WorkingDir == "" {
			task.Spec.StepTemplate.WorkingDir = "/workspace/source"
		}
		for i := range task.Spec.Steps {
			s := &task.Spec.Steps[i]

			s.SecurityContext = nil
			if s.WorkingDir == "/workspace/source" {
				s.WorkingDir = ""
			}

			// share common volume mounts
			for _, vm := range s.VolumeMounts {
				found := false
				for _, vm2 := range task.Spec.StepTemplate.VolumeMounts {
					if vm2.Name == vm.Name {
						found = true
						break
					}
				}
				if !found {
					task.Spec.StepTemplate.VolumeMounts = append(task.Spec.StepTemplate.VolumeMounts, vm)
				}
			}
			s.VolumeMounts = nil

			// share resources
			for k, q := range s.Resources.Limits {
				if !q.IsZero() {
					if task.Spec.StepTemplate.Resources.Limits == nil {
						task.Spec.StepTemplate.Resources.Limits = map[corev1.ResourceName]resource.Quantity{}
					}
					task.Spec.StepTemplate.Resources.Limits[k] = q
				}
			}
			for k, q := range s.Resources.Requests {
				if !q.IsZero() {
					if task.Spec.StepTemplate.Resources.Requests == nil {
						task.Spec.StepTemplate.Resources.Requests = map[corev1.ResourceName]resource.Quantity{}
					}
					task.Spec.StepTemplate.Resources.Requests[k] = q
				}
			}
			s.Resources.Limits = map[corev1.ResourceName]resource.Quantity{}
			s.Resources.Requests = map[corev1.ResourceName]resource.Quantity{}

			// filter out env vars
			var envs2 []corev1.EnvVar
			for _, e := range s.Env {
				if util.StringArrayIndex(templateEnvVars, e.Name) >= 0 || util.StringArrayIndex(removeEnvVarList, e.Name) >= 0 {
					found := false
					for _, et := range task.Spec.StepTemplate.Env {
						if et.Name == e.Name {
							found = true
							break
						}
					}
					if !found {
						task.Spec.StepTemplate.Env = append(task.Spec.StepTemplate.Env, e)
					}
					continue
				}
				envs2 = append(envs2, e)
			}
			s.Env = envs2
		}

		// lets remove old volume mounts
		stepTemplate := task.Spec.StepTemplate
		if stepTemplate != nil {
			for i, vm := range stepTemplate.VolumeMounts {
				if vm.Name == "workspace-volume" {
					old := stepTemplate.VolumeMounts
					stepTemplate.VolumeMounts = old[0:i]
					if i+1 < len(old) {
						stepTemplate.VolumeMounts = append(stepTemplate.VolumeMounts, old[i+1:]...)
					}
					break
				}
			}
		}
	}

	for i, t := range tasks {
		if len(t.Spec.Steps) > 0 {
			t.Spec.Steps[0].Script = scripts[i]
		}
	}

	return answer, nil
}

// CombinePipelinesAndTasksIntoRun combines the pipelines to a single PipelineRun
func CombinePipelinesAndTasksIntoRun(results *tekton.CRDWrapper) *v1alpha1.PipelineRun {
	run := results.PipelineRun()
	pipeline := results.Pipeline()
	tasks := results.Tasks()

	if pipeline == nil {
		return run
	}
	run.Spec.PipelineSpec = &pipeline.Spec
	run.Spec.PipelineRef = nil

	oldTasks := run.Spec.PipelineSpec.Tasks
	run.Spec.PipelineSpec.Tasks = nil
	for i, t := range tasks {
		pt := v1alpha1.PipelineTask{}
		if len(oldTasks) > i {
			pt = oldTasks[i]
		}
		pt.TaskRef = nil
		pt.TaskSpec = &t.Spec

		if t.Spec.Inputs != nil {
			for j := range t.Spec.Inputs.Params {
				ps := &t.Spec.Inputs.Params[j]

				// lets make sure there's a param spec in the pipeline spec
				found := false
				for _, pps := range run.Spec.PipelineSpec.Params {
					if pps.Name == ps.Name {
						found = true
						break
					}
				}

				if !found {
					run.Spec.PipelineSpec.Params = append(run.Spec.PipelineSpec.Params, *ps)
				}

				found = false
				for _, pts := range pt.Params {
					if pts.Name == ps.Name {
						found = true
						break
					}
				}
				if !found {
					pt.Params = append(pt.Params, v1alpha1.Param{
						Name: ps.Name,
						Value: v1beta1.ArrayOrString{
							Type:      v1beta1.ParamTypeString,
							StringVal: fmt.Sprintf("$(params.%s)", ps.Name),
						},
					})
				}
			}
		}

		run.Spec.PipelineSpec.Tasks = append(run.Spec.PipelineSpec.Tasks, pt)
	}
	run.Spec.PipelineRef = nil
	return run
}

func transformStepArgs(args []string) []string {
	answer := []string{}
	for _, arg := range args {
		arg = strings.ReplaceAll(arg, "/workspace/source/charts/myrepo", "/workspace/source/charts/$(inputs.params.repo_name)")
		answer = append(answer, arg)
	}
	return answer
}

// AddMissingParamSpecs adds any missing input parameters
func AddMissingParamSpecs(params []v1alpha1.ParamSpec, defaults []v1alpha1.ParamSpec) []v1alpha1.ParamSpec {
	for _, p := range defaults {
		found := false
		for _, param := range params {
			if param.Name == p.Name {
				found = true
			}
		}
		if !found {
			params = append(params, p)
		}
	}
	return params
}

// AddMissingParamSpecs adds any missing input parameters
func AddMissingParams(params []v1alpha1.Param, defaults []v1alpha1.ParamSpec) []v1alpha1.Param {
	for _, p := range defaults {
		found := false
		for _, param := range params {
			if param.Name == p.Name {
				found = true
			}
		}
		if !found {
			params = append(params, v1alpha1.Param{
				Name: p.Name,
				Value: v1beta1.ArrayOrString{
					Type:      v1beta1.ParamTypeString,
					StringVal: fmt.Sprintf("$(params.%s)", p.Name),
				},
			})
		}
	}
	return params
}

// AddDefaultParams adds any missing input parameters
func AddDefaultParams(params []v1alpha1.Param, defaults []v1alpha1.ParamSpec) []v1alpha1.Param {
	for _, p := range defaults {
		found := false
		for _, param := range params {
			if param.Name == p.Name {
				found = true
			}
		}
		if !found {
			params = append(params, v1alpha1.Param{
				Name: p.Name,
				Value: v1beta1.ArrayOrString{
					Type:      v1beta1.ParamTypeString,
					StringVal: "",
				},
			})
		}
	}
	return params
}

func AddMissingParamEnv(env []corev1.EnvVar, defaults []v1alpha1.ParamSpec) []corev1.EnvVar {
	for _, p := range defaults {
		found := false
		for _, e := range env {
			if e.Name == p.Name {
				found = true
			}
		}
		if !found {
			env = append(env, corev1.EnvVar{
				Name:  p.Name,
				Value: fmt.Sprintf("$(params.%s)", p.Name),
			})
		}
	}
	return env

}

// RemoveOldParamSpecs removes any old parameters
func RemoveOldParamSpecs(params []v1alpha1.ParamSpec, names ...string) []v1alpha1.ParamSpec {
	for _, name := range names {
		// lets handle duplicate parameters
		for {
			found := false
			for i, param := range params {
				if param.Name == name {
					params = append(params[0:i], params[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				break
			}
		}
	}
	return params
}

// RemoveOldParams removes old params
func RemoveOldParams(params []v1alpha1.Param, names ...string) []v1alpha1.Param {
	for _, name := range names {
		// lets handle duplicate parameters
		for {
			found := false
			for i, param := range params {
				if param.Name == name {
					params = append(params[0:i], params[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				break
			}
		}
	}
	return params
}

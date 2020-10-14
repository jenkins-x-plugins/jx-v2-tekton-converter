package transform

import (
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// CreateSimpleGitCloneTask creates a simple git clone task for pre-0.9.0 tekton
func CreateSimpleGitCloneTask(o *Options) (*v1alpha1.Task, error) {
	s := v1alpha1.Step{}
	s.Command = []string{"/bin/sh"}
	s.Args = []string{"-c", "mkdir -p $HOME; git config --global --add user.name ${GIT_AUTHOR_NAME:-jenkins-x-bot}; git config --global --add user.email ${GIT_AUTHOR_EMAIL:-jenkins-x@googlegroups.com}; git config --global credential.helper store; git clone $(params.REPO_URL) $(params.subdirectory); echo cloned url: $(params.REPO_URL) to dir: $(params.subdirectory); cd $(params.subdirectory); git checkout $(params.PULL_PULL_SHA); echo checked out revision: $(params.PULL_PULL_SHA) to dir: $(params.subdirectory)"}
	s.Name = "git-clone"

	// TODO use version stream
	s.Image = "gcr.io/jenkinsxio/builder-jx:2.1.32-662"

	s.VolumeMounts = []corev1.VolumeMount{
		{
			Name:      "workspace-volume",
			MountPath: "/home/jenkins",
		},
	}
	s.WorkingDir = "/workspace"

	s.Env = []corev1.EnvVar{
		/*
			{
				Name:  "GIT_AUTHOR_EMAIL",
				Value: "$(params.git_author_email)",
			},
			{
				Name:  "GIT_AUTHOR_NAME",
				Value: "$(params.git_author_name)",
			},
			{
				Name:  "GIT_COMMITTER_EMAIL",
				Value: "$(params.git_committer_email)",
			},
			{
				Name:  "GIT_COMMITTER_NAME",
				Value: "$(params.git_committer_name)",
			},
			{
				Name:  "XDG_CONFIG_HOME",
				Value: "/workspace/xdg_config",
			},
			{
				Name:  "REPO_OWNER",
				Value: "$(params.repo_owner)",
			},
			{
				Name:  "REPO_NAME",
				Value: "$(params.repo_name)",
			},
			{
				Name:  "BRANCH_NAME",
				Value: "$(params.revision)",
			},
		*/
	}
	task := &v1alpha1.Task{}
	task.Spec.Steps = append(task.Spec.Steps, s)
	task.Spec.Params = []v1alpha1.ParamSpec{
		{
			Name:        "REPO_URL",
			Type:        v1alpha1.ParamTypeString,
			Description: "git url to clone",
		},
		{
			Name:        "PULL_PULL_SHA",
			Type:        v1alpha1.ParamTypeString,
			Description: "git revision to checkout (branch, tag, sha, ref…)",
			Default: &v1alpha1.ArrayOrString{
				Type:      v1alpha1.ParamTypeString,
				StringVal: "master",
			},
		},
		{
			Name:        "subdirectory",
			Type:        v1alpha1.ParamTypeString,
			Description: "subdirectory inside of /workspace to clone the git repo",
			Default: &v1alpha1.ArrayOrString{
				Type:      v1alpha1.ParamTypeString,
				StringVal: o.CreateTaskOptions.SourceName,
			},
		},
	}
	return task, nil
}

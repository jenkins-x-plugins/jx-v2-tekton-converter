package transform

import (
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

var (
	// avoidCDs avoid the cd command before these commands
	avoidCDs = []string{
		"jx step changelog",
		"jx step helm release",
		"jx promote",
	}

	removeEnvVars = []string{
		"GIT_AUTHOR_EMAIL", "GIT_AUTHOR_NAME", "GIT_COMMITTER_EMAIL", "GIT_COMMITTER_NAME",
		"REPO_OWNER", "REPO_NAME", "APP_NAME",
		"BRANCH_NAME", "DOCKER_REGISTRY", "DOCKER_REGISTRY_ORG",
		"PREVIEW_VERSION",
		"JOB_NAME",
		"XDG_CONFIG_HOME",
		"BUILD_NUMBER",
		"NO_GOOGLE_APPLICATION_CREDENTIALS",
		"GOOGLE_APPLICATION_CREDENTIALS",
		"VERSION",
		"JX_BATCH_MODE",
		"PIPELINE_KIND",
	}

	templateEnvVars = []string{
		"JX_BATCH_MODE",
		"PIPELINE_KIND",
	}

	replaceEnvVars = []corev1.EnvVar{}

	defaultTaskParams = []v1alpha1.ParamSpec{
		/*
			{
				Name:        "BUILD_ID",
				Type:        v1alpha1.ParamTypeString,
				Description: "the unique build number",
			},
			{
				Name:        "JOB_NAME",
				Type:        v1alpha1.ParamTypeString,
				Description: "the name of the job which is the trigger context name",
			},
			{
				Name:        "JOB_SPEC",
				Type:        v1alpha1.ParamTypeString,
				Description: "the specification of the job",
			},
			{
				Name:        "JOB_TYPE",
				Type:        v1alpha1.ParamTypeString,
				Description: "the kind of job: postsubmit or presubmit",
			},
			{
				Name:        "PULL_BASE_REF",
				Type:        v1alpha1.ParamTypeString,
				Description: "the base git reference of the pull request",
				Default: &v1alpha1.ArrayOrString{
					Type:      v1alpha1.ParamTypeString,
					StringVal: "master",
				},
			},
			{
				Name:        "PULL_BASE_SHA",
				Type:        v1alpha1.ParamTypeString,
				Description: "the git sha of the base of the pull request",
			},
			{
				Name:        "PULL_NUMBER",
				Type:        v1alpha1.ParamTypeString,
				Description: "git pull request number",
				Default: &v1alpha1.ArrayOrString{
					Type:      v1alpha1.ParamTypeString,
					StringVal: "",
				},
			},
			{
				Name:        "PULL_PULL_REF",
				Type:        v1alpha1.ParamTypeString,
				Description: "git pull request ref in the form 'refs/pull/$PULL_NUMBER/head'",
				Default: &v1alpha1.ArrayOrString{
					Type:      v1alpha1.ParamTypeString,
					StringVal: "",
				}},
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
				Name:        "PULL_REFS",
				Type:        v1alpha1.ParamTypeString,
				Description: "git pull reference strings of base and latest in the form 'master:$PULL_BASE_SHA,$PULL_NUMBER:$PULL_PULL_SHA:refs/pull/$PULL_NUMBER/head'",
			},
			{
				Name:        "REPO_NAME",
				Type:        v1alpha1.ParamTypeString,
				Description: "git repository name",
			},
			{
				Name:        "REPO_OWNER",
				Type:        v1alpha1.ParamTypeString,
				Description: "git repository owner (user or organisation)",
			},
			{
				Name:        "REPO_URL",
				Type:        v1alpha1.ParamTypeString,
				Description: "git url to clone",
			},
		*/
	}
)

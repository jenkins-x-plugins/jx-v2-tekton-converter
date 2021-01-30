package transform

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/jenkins-x/jx/v2/pkg/config"
	"github.com/jenkins-x/jx/v2/pkg/jenkinsfile"
	"github.com/jenkins-x/jx/v2/pkg/tekton"
	"github.com/jenkins-x/jx/v2/pkg/util"
	"github.com/pkg/errors"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/yaml"
)

const (
	triggersYaml = `apiVersion: config.lighthouse.jenkins-x.io/v1alpha1
kind: TriggerConfig
spec:
  presubmits:
  - name: pr
    context: "pr"
    always_run: true
    optional: false
    trigger: "/test"
    rerun_command: "/retest"
    source: "pullrequest.yaml"
  postsubmits:
  - name: release
    context: "release"
    source: "release.yaml"
    branches:
    - main
    - master
`
)

func CreateCatalogForPackDir(o *Options, rootTmpDir string, pack string) error {
	outDir := o.OutDir
	appName := "todo"

	dummySourceGitURL := "https://github.com/jenkins-x-quickstarts/node-http"
	dummySourceDir := filepath.Join(rootTmpDir, "projects", pack, appName)
	err := os.MkdirAll(dummySourceDir, util.DefaultWritePermissions)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", dummySourceDir)
	}

	c := util.Command{
		Name: "git",
		Args: []string{"clone", dummySourceGitURL, dummySourceDir},
	}
	_, err = c.RunWithoutRetry()
	if err != nil {
		return errors.Wrapf(err, "failed to git %s", strings.Join(c.Args, " "))
	}
	log.Logger().Infof("cloned dummy source repo %s to %s", dummySourceGitURL, dummySourceDir)

	// lets create a project config file
	projectConfig := &config.ProjectConfig{
		BuildPack:        pack,
		BuildPackGitURL:  o.BuildPackURL,
		BuildPackGitURef: o.BuildPackRef,
	}

	projectConfigFile := filepath.Join(dummySourceDir, config.ProjectConfigFileName)
	err = projectConfig.SaveConfig(projectConfigFile)
	if err != nil {
		return errors.Wrapf(err, "failed to save file %s", projectConfigFile)
	}

	packRootDir := filepath.Join(outDir, pack, ".lighthouse", "jenkins-x")

	return ConvertDirectory(o, pack, dummySourceGitURL, dummySourceDir, projectConfigFile, packRootDir)
}

// ConvertDirectory converts the given directory
func ConvertDirectory(o *Options, description string, gitURL string, sourceDir string, projectConfigFile string, outDir string) error {
	copy := o.CreateTaskOptions
	so := &copy
	so.Branch = "master"
	so.CloneDir = o.Dir

	pipelineKinds := []string{jenkinsfile.PipelineKindRelease, jenkinsfile.PipelineKindPullRequest}
	for _, kind := range pipelineKinds {
		so.PipelineKind = kind
		so.DryRun = true
		so.Pack = description
		so.DisableGitClone = true
		so.NoOutput = true
		so.CloneGitURL = gitURL

		so.CloneDir = sourceDir
		log.Logger().Debugf("using dummy source project at %s", util.ColorInfo(projectConfigFile))

		so.OutDir = outDir

		err := so.Run()
		if err != nil {
			return errors.Wrapf(err, "failed to generate tekton resources to %s", outDir)
		}

		results := &so.Results
		if kind == jenkinsfile.PipelineKindRelease {
			if so.EffectiveProjectConfig == nil {
				log.Logger().Warnf("no EffectiveProjectConfig after generating pipeline for release of %s", description)
			} else {
				results, err = AddSetVersionSteps(o, results, so.EffectiveProjectConfig)
				if err != nil {
					return errors.Wrapf(err, "failed to add release setVersion steps to tekton resources for %s kind %s", description, so.PipelineKind)
				}
			}
		}
		results, err = TransformResources(o, results, kind)
		if err != nil {
			return errors.Wrapf(err, "failed to transform tekton resources for %s kind %s", description, so.PipelineKind)
		}
		err = WriteToDisk(results, outDir, so.PipelineKind, true)
		if err != nil {
			return errors.Wrapf(err, "failed to save tekton resources to %s", outDir)
		}
		log.Logger().Debugf("generated tekton resources to %s", util.ColorInfo(outDir))
	}
	flowFile := filepath.Join(outDir, "triggers.yaml")
	err := ioutil.WriteFile(flowFile, []byte(triggersYaml), util.DefaultFileWritePermissions)
	if err != nil {
		return errors.Wrapf(err, "failed to save file %s", flowFile)
	}
	return nil
}

func WriteToDisk(results *tekton.CRDWrapper, packRootDir string, kind string, combineToPipelineRun bool) error {
	if !combineToPipelineRun {
		packDir := filepath.Join(packRootDir, kind)
		err := os.MkdirAll(packDir, util.DefaultWritePermissions)
		if err != nil {
			return errors.Wrapf(err, "failed to create directory %s", packDir)
		}
		return results.WriteToDisk(packDir, nil)
	}

	err := os.MkdirAll(packRootDir, util.DefaultWritePermissions)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", packRootDir)
	}

	pr := CombinePipelinesAndTasksIntoRun(results)

	// lets convert to a v1beta1
	pr2 := &v1beta1.PipelineRun{}

	ctx := context.Background()

	err = pr.ConvertTo(ctx, pr2)
	if err != nil {
		return errors.Wrapf(err, "failed to convert PipelineRun %s to v1beta1", pr.Name)
	}
	if pr2.APIVersion == "" {
		pr2.APIVersion = "tekton.dev/v1beta1"
	}
	if pr2.Kind == "" {
		pr2.Kind = "PipelineRun"
	}
	data, err := yaml.Marshal(pr2)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal PipelineRun YAML")
	}
	fileName := filepath.Join(packRootDir, kind+".yaml")
	err = ioutil.WriteFile(fileName, data, util.DefaultWritePermissions)
	if err != nil {
		return errors.Wrapf(err, "failed to save PipelineRun file %s", fileName)
	}
	log.Logger().Infof("generated pipelineRun at %s", util.ColorInfo(fileName))
	return nil
}

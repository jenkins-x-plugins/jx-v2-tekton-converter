package transform

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/jenkins-x/jx/v2/pkg/config"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Run implements the command
func (o *Options) Run() error {
	var err error
	if o.OutDir == "" {
		o.OutDir, err = ioutil.TempDir("", "jx-v2-tekton-converter-")
		if err != nil {
			return errors.Wrap(err, "failed to create temporary dir for the output")
		}
	}

	rootTmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return errors.Wrap(err, "failed to create temp dir")
	}

	if o.Pack != "" {
		return CreateCatalogForPackDir(o, rootTmpDir, o.Pack)
	}

	if o.BuildPack {
		return o.ConvertBuildPack(rootTmpDir)
	}

	projectConfig, projectConfigFile, err := config.LoadProjectConfig(o.Dir)
	if err != nil {
		return errors.Wrapf(err, "failed to load project configuration in dir %s", o.Dir)
	}
	if projectConfigFile == "" {
		return errors.Errorf("could not find jenkins-x.yml file in dir %s", o.Dir)
	}

	if (projectConfig.BuildPackGitURL == "" || projectConfig.BuildPackGitURL == "https://github.com/jenkins-x/jxr-packs-kubernetes.git") && projectConfig.BuildPack != "" && projectConfig.PipelineConfig == nil {
		// there is no local configuration so lets just replace with the current pipelines from the
		// pipeline catalog
		c := exec.Command("jx", "project", "enable", "--pack", projectConfig.BuildPack)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			return errors.Wrapf(err, "failed to run command: %s", c.String())
		}

		err = os.Remove(projectConfigFile)
		if err != nil {
			return errors.Wrapf(err, "failed to remove %s", projectConfigFile)
		}
		log.Logger().Infof("removed old file %s", projectConfigFile)
		return nil
	}

	if projectConfig.BuildPackGitURL == "" {
		projectConfig.BuildPackGitURL = o.BuildPackURL

		if projectConfig.BuildPackGitURef == "" && o.BuildPackRef != "" {
			projectConfig.BuildPackGitURef = o.BuildPackRef
		}

		err = projectConfig.SaveConfig(projectConfigFile)
		if err != nil {
			return errors.Wrapf(err, "failed to save %s", projectConfigFile)
		}
	}

	outDir := filepath.Join(o.Dir, ".lighthouse", "jenkins-x")
	err = ConvertDirectory(o, projectConfig.BuildPack, "", o.Dir, projectConfigFile, outDir)
	if err != nil {
		return errors.Wrapf(err, "failed to convert directory %s", o.Dir)
	}

	err = o.CreateTaskOptions.Git().Add(o.Dir, ".lighthouse")
	if err != nil {
		return errors.Wrapf(err, "failed to add new files to git")
	}

	err = os.Remove(projectConfigFile)
	if err != nil {
		return errors.Wrapf(err, "failed to remove %s", projectConfigFile)
	}
	log.Logger().Infof("removed old file %s", projectConfigFile)
	return nil
}

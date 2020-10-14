package transform

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/jenkins-x/jx/v2/pkg/errorutil"
	"github.com/jenkins-x/jx/v2/pkg/gits"
	"github.com/jenkins-x/jx/v2/pkg/util"
	"github.com/pkg/errors"
)

// ConvertBuildPack converts a build pack over to the new format
func (o *Options) ConvertBuildPack(rootTmpDir string) error {

	// lets clone the repository and iterate through each pack
	dir, err := ioutil.TempDir("", "jx-build-pack-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory")
	}

	gitter := gits.NewGitCLI()
	err = gitter.Clone(o.BuildPackURL, dir)
	if err != nil {
		return errors.Wrapf(err, "failed to clone %s to dir %s", o.BuildPackURL, dir)
	}

	errs := []error{}
	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if info.Name() == "pipeline.yaml" {
				packDir := filepath.Dir(path)
				if packDir == dir {
					return nil
				}
				_, pack := filepath.Split(packDir)
				if pack != "" {
					log.Logger().Infof("processing build pack dir %s\n", util.ColorInfo(pack))

					packRootDir := filepath.Join(o.OutDir, pack)
					err = os.MkdirAll(packRootDir, util.DefaultWritePermissions)
					if err != nil {
						return errors.Wrapf(err, "failed to make dir %s", packRootDir)
					}

					dockerfile := filepath.Join(packDir, "Dockerfile")
					exists, err := util.FileExists(dockerfile)
					if err != nil {
						return errors.Wrapf(err, "failed to check if file exists %s", dockerfile)
					}
					if exists {
						// lets use kpt to generate the charts folder...
						targetDockerfile := filepath.Join(packRootDir, "Dockerfile")
						err = util.CopyFile(dockerfile, targetDockerfile)
						if err != nil {
							return errors.Wrapf(err, "failed to copy %s to %s", dockerfile, targetDockerfile)
						}
					}

					chartsDir := filepath.Join(packDir, "charts")
					exists, err = util.DirExists(chartsDir)
					if err != nil {
						return errors.Wrapf(err, "failed to check if dir exists %s", chartsDir)
					}
					if exists {
						// lets use kpt to generate the charts folder...
						targetChartsDir := filepath.Join(packRootDir, "charts")
						exists, err := util.DirExists(targetChartsDir)
						if err != nil {
							return errors.Wrapf(err, "failed to check if dir exists %s", targetChartsDir)
						}
						if !exists {
							log.Logger().Infof("creating a charts dir to %s", targetChartsDir)

							c := util.Command{
								Dir:  packRootDir,
								Name: "kpt",
								Args: []string{"pkg", "get", "https://github.com/jenkins-x/jx3-gitops-template.git/charts/charts", "."},
							}
							_, err = c.RunWithoutRetry()
							if err != nil {
								log.Logger().Warnf("got kpt error: %s", err.Error())
							}

							c = util.Command{
								Dir:  packRootDir,
								Name: "kpt",
								Args: []string{"pkg", "get", "https://github.com/jenkins-x/jx3-gitops-template.git/charts/preview", "."},
							}
							_, err = c.RunWithoutRetry()
							if err != nil {
								log.Logger().Warnf("got kpt error: %s", err.Error())
							}
						}
					}

					err = CreateCatalogForPackDir(o, rootTmpDir, pack)
					if err != nil {
						errs = append(errs, errors.Wrapf(err, "failed to process build pack %s", pack))
					}
				}
			}
			return nil
		})
	if err != nil {
		errs = append(errs, err)
	}
	err = errorutil.CombineErrors(errs...)
	if err != nil {
		return errors.Wrapf(err, "failed to process build packs in dir %s", dir)
	}
	return nil
}

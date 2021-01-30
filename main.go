package main

import (
	"fmt"
	"os"

	"github.com/jenkins-x-plugins/jx-v2-tekton-converter/pkg/transform"
	"github.com/jenkins-x/jx/v2/pkg/cmd/clients"
	"github.com/jenkins-x/jx/v2/pkg/cmd/opts"
	"github.com/jenkins-x/jx/v2/pkg/tekton/syntax"
	"github.com/spf13/cobra"
)

func main() {
	cmd, _ := NewMain()
	cmd.Execute()
}

// NewMain creates a command object
func NewMain() (*cobra.Command, *transform.Options) {
	o := &transform.Options{}
	so := &o.CreateTaskOptions
	so.CommonOptions = opts.NewCommonOptionsWithTerm(clients.NewFactory(), os.Stdin, os.Stdout, os.Stderr)

	cmd := &cobra.Command{
		Short: "Generates Tekton Pipelines for a Jenkins X Build Pack or Jenkins X YAML file",
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run()
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringVarP(&o.BuildPackURL, "url", "u", "https://github.com/jenkins-x/jxr-packs-kubernetes.git", "the build pack git URL")
	cmd.Flags().StringVarP(&o.BuildPackRef, "ref", "r", "", "the build pack git reference")
	cmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "the directory to look for the build pack or defaults to the currect directory")
	cmd.Flags().StringVarP(&o.Pack, "pack", "p", "", "the specific build pack directory to generate. If not specified all the build pack directories are generated")
	cmd.Flags().StringVarP(&o.OutDir, "out-dir", "o", "", "the directory to save the resources to. Uses a temporary directory if none is specified")
	// TODO should use the current version stream to replace the version number
	cmd.Flags().StringVarP(&o.DefaultJXImage, "default-jx-image", "", "gcr.io/jenkinsxio-labs-private/jxl:0.0.136", "the default image used for jx steps in the release setVersion steps")
	cmd.Flags().BoolVarP(&o.Verbose, "verbose", "v", false, "enable verbose logging")
	cmd.Flags().BoolVarP(&o.GitCloneUsesSteps, "uses-git-clone", "", true, "use the uses:sourceURI git clone steps")
	cmd.Flags().BoolVarP(&o.UseCatalogGitClone, "catalog-git-clone", "", false, "if enabled uses the catalog git-clone task. Requires tekton 0.9.x or later as it requires the script tag on a step")

	// step flags
	cmd.Flags().StringVarP(&so.DefaultImage, "default-image", "", syntax.DefaultContainerImage, "Specify the docker image to use if there is no image specified for a step and there's no Pod Template")
	cmd.Flags().StringVarP(&so.SourceName, "source", "", "source", "The name of the source repository")
	cmd.Flags().StringVarP(&so.ServiceAccount, "service-account", "", "tekton-bot", "The Kubernetes ServiceAccount to use to run the pipeline")
	cmd.Flags().StringVarP(&so.KanikoImage, "kaniko-image", "", syntax.KanikoDockerImage, "The docker image for Kaniko")
	return cmd, o
}

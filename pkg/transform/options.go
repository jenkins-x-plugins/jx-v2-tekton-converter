package transform

import (
	"github.com/jenkins-x/jx/v2/pkg/cmd/step/create"
)

// Options the CLI Options for the command
type Options struct {
	CreateTaskOptions  create.StepCreateTaskOptions
	BuildPackURL       string
	BuildPackRef       string
	Dir                string
	Pack               string
	OutDir             string
	DefaultJXImage     string
	JXImage            string
	BuildPack          bool
	Verbose            bool
	SemanticRelease    bool
	UseCatalogGitClone bool
}

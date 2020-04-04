package archive

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/AnastasiaYarema/design-practice-2/build/archive")

	// Ninja rule to archive binary output file.
	makeArchive = pctx.StaticRule("makeArchive", blueprint.RuleParams{
		Command:     "cd $workDir && zip $outputPath $inputPath",
		Description: "make archive from $inputPath binary",
	}, "workDir", "outputPath", "inputPath")
)

type zipArchiveModule struct {
	blueprint.SimpleName

	properties struct {
		// Go binary filename to archive
		Binary string
	}
}

func (gb *zipArchiveModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	archiveName := name + ".zip"
	config := bood.ExtractConfig(ctx)
	inputPath := path.Join(config.BaseOutputDir, "bin", gb.properties.Binary)
	outputPath := path.Join(config.BaseOutputDir, "archives", archiveName)
	fmt.Println(inputPath)
	fmt.Println(outputPath)

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as zip archive", name),
		Rule:        makeArchive,
		Outputs:     []string{outputPath},
		Implicits:   []string{inputPath},
		Args: map[string]string{
			"workDir":    ctx.ModuleDir(),
			"outputPath": outputPath,
			"inputPath":  inputPath,
		},
	})
}

// SimpleArchiveFactory is a factory for go binary module type which supports Go command packages with running tests.
func SimpleArchiveFactory() (blueprint.Module, []interface{}) {
	mType := &zipArchiveModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
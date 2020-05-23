package gomodule

import (
	"testing"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"os"
	"os/exec"
)

func NewContext() *blueprint.Context {
	ctx := bood.PrepareContext()

	ctx.RegisterModuleType("go_binary", SimpleBinFactory)
	return ctx
}

// Hook up gocheck into "go test" runner
func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {}

var _ = Suite(&TestSuite{})

func (s * TestSuite) TestArchiveModule(c *C) {
	buildFile := []byte(
		"go_binary {\n" + 
  		"name: \"bood\",\n" + 
		"pkg: \"github.com/AnastasiaYarema/design-practice-2/build/cmd/bood\",\n" +
  		"testPkg: \"github.com/AnastasiaYarema/design-practice-2/build/cmd/bood\",\n" +
  		"srcs: [\"**/*.go\", \"../go.mod\"]\n" +
		"}")

	err := ioutil.WriteFile("./build.bood", buildFile, 0644)
	c.Assert(err, IsNil)

	config := bood.NewConfig()
	ctx := NewContext()
	ninjaBuildPath := bood.GenerateBuildFile(config, ctx)

	cmd := exec.Command("ninja", append([]string{"-f", ninjaBuildPath})...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		c.Error(err)
	}

	// Test whether specified binary was created
	if _, err := os.Stat("./out/bin/bood"); os.IsNotExist(err) {
		c.Error("Bood binary was not create")
	}

	// Test whether specified binary was created
	if _, err := os.Stat("./out/reports/bood.txt"); os.IsNotExist(err) {
		c.Error("Bood report file was not create")
	}

	// remove build.bood config
	buildBoodRemovingErr := os.Remove("./build.bood")
	c.Assert(buildBoodRemovingErr, IsNil)

	// remove out dir
	outRemovingErr := os.RemoveAll("./out")
	c.Assert(outRemovingErr, IsNil)
}
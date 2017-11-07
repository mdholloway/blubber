package docker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"phabricator.wikimedia.org/source/blubber/build"
	"phabricator.wikimedia.org/source/blubber/docker"
)

func TestRun(t *testing.T) {
	i := build.Run{"echo", []string{"hello"}}
	di, err := docker.NewInstruction(i)

	var dockerRun docker.Run

	assert.Nil(t, err)
	assert.IsType(t, dockerRun, di)
	assert.Equal(t, "RUN echo \"hello\"\n", di.Compile())
}

func TestRunAll(t *testing.T) {
	i := build.RunAll{[]build.Run{
		{"echo", []string{"hello"}},
		{"echo", []string{"yo"}},
	}}

	di, err := docker.NewInstruction(i)

	var dockerRun docker.Run

	assert.Nil(t, err)
	assert.IsType(t, dockerRun, di)
	assert.Equal(t, "RUN echo \"hello\" && echo \"yo\"\n", di.Compile())
}

func TestCopy(t *testing.T) {
	i := build.Copy{[]string{"foo1", "foo2"}, "bar"}

	di, err := docker.NewInstruction(i)

	var dockerCopy docker.Copy

	assert.Nil(t, err)
	assert.IsType(t, dockerCopy, di)
	assert.Equal(t, "COPY [\"foo1\", \"foo2\", \"bar\"]\n", di.Compile())
}

func TestCopyFrom(t *testing.T) {
	i := build.CopyFrom{"foo", build.Copy{[]string{"foo1", "foo2"}, "bar"}}

	di, err := docker.NewInstruction(i)

	var dockerCopyFrom docker.CopyFrom

	assert.Nil(t, err)
	assert.IsType(t, dockerCopyFrom, di)
	assert.Equal(t, "COPY --from=foo [\"foo1\", \"foo2\", \"bar\"]\n", di.Compile())
}

func TestEnv(t *testing.T) {
	i := build.Env{map[string]string{"foo": "bar", "bar": "foo"}}

	di, err := docker.NewInstruction(i)

	var dockerEnv docker.Env

	assert.Nil(t, err)
	assert.IsType(t, dockerEnv, di)
	assert.Equal(t, "ENV bar=\"foo\" foo=\"bar\"\n", di.Compile())
}

func TestLabel(t *testing.T) {
	i := build.Label{map[string]string{"foo": "bar", "bar": "foo"}}

	di, err := docker.NewInstruction(i)

	var dockerLabel docker.Label

	assert.Nil(t, err)
	assert.IsType(t, dockerLabel, di)
	assert.Equal(t, "LABEL bar=\"foo\" foo=\"bar\"\n", di.Compile())
}

func TestVolume(t *testing.T) {
	i := build.Volume{"/foo/dir"}

	di, err := docker.NewInstruction(i)

	var dockerVolume docker.Volume

	assert.Nil(t, err)
	assert.IsType(t, dockerVolume, di)
	assert.Equal(t, "VOLUME [\"/foo/dir\"]\n", di.Compile())
}

func TestEscapeRun(t *testing.T) {
	i := build.Run{"/bin/true\nRUN echo HACKED!", []string{}}
	dr, _ := docker.NewInstruction(i)

	assert.Equal(t, "RUN /bin/true\\nRUN echo HACKED!\n", dr.Compile())
}

func TestEscapeCopy(t *testing.T) {
	i := build.Copy{[]string{"file.a", "file.b"}, "dest"}
	dr, _ := docker.NewInstruction(i)

	assert.Equal(t, "COPY [\"file.a\", \"file.b\", \"dest\"]\n", dr.Compile())
}

func TestEscapeEnv(t *testing.T) {
	i := build.Env{map[string]string{"a": "b\nRUN echo HACKED!"}}
	dr, _ := docker.NewInstruction(i)

	assert.Equal(t, "ENV a=\"b\\nRUN echo HACKED!\"\n", dr.Compile())
}

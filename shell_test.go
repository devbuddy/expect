package expect

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

const DockerImage = "docker-test-image"

func Test_ShellExpect(t *testing.T) {
	tests := []struct {
		shell     string
		shellArgs []string
		env       []string
	}{
		{
			shell:     "zsh",
			shellArgs: []string{"--no-globalrcs", "--no-rcs", "--no-zle", "--no-promptcr"},
			env: []string{
				"PROMPT=##\n",
				"TESTVAR=foobar",
			},
		},
		{
			shell:     "bash",
			shellArgs: []string{"--noprofile", "--norc"},
			env: []string{
				"PS1=##\n",
				"TESTVAR=foobar",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			shellPath, err := exec.LookPath(tt.shell)
			if err != nil {
				t.Skipf("shell executable not found for %s (%s)", tt.shell, err)
			}

			ep := NewExpectWithEnv(shellPath, tt.shellArgs, tt.env)
			err = ep.Start()
			require.NoError(t, err)

			shell := NewShellExpect(ep, "##\n")

			t.Run("init", func(t *testing.T) {
				err = shell.Init()
				require.NoError(t, err)
			})

			t.Run("echo", func(t *testing.T) {
				output, err := shell.Run("echo $TESTVAR")
				require.NoError(t, err)
				require.Equal(t, []string{"foobar"}, output)
			})
		})
	}
}

func Test_ShellExpect_Docker_Bash(t *testing.T) {
	args := []string{
		"docker", "run", "-ti", "--rm",
		"-e", "PS1=##\n",
		"-e", "TESTVAR=foobar",
		"--entrypoint", "/bin/bash",
		DockerImage,
		"--noprofile", "--norc",
	}

	ep := NewExpect(args[0], args[1:]...)
	err := ep.Start()
	require.NoError(t, err)

	shell := NewShellExpect(ep, "##\n")

	err = shell.Init()
	require.NoError(t, err)

	output, err := shell.Run("stty -echo") // disable echo inside the container
	require.NoError(t, err)
	require.Equal(t, []string{"stty -echo"}, output)

	output, err = shell.Run("echo $TESTVAR")
	require.NoError(t, err)
	require.Equal(t, []string{"foobar"}, output)
}

func Test_ShellExpect_Docker_Zsh(t *testing.T) {
	args := []string{
		"docker", "run", "-ti", "--rm",
		"-e", "PROMPT=##\n",
		"-e", "TESTVAR=foobar",
		"--entrypoint", "/bin/zsh",
		DockerImage,
		"--no-globalrcs", "--no-rcs", "--no-zle", "--no-promptcr",
	}

	ep := NewExpect(args[0], args[1:]...)
	err := ep.Start()
	require.NoError(t, err)
	ep.Debug = true

	shell := NewShellExpect(ep, "##\n")

	err = shell.Init()
	require.NoError(t, err)

	output, err := shell.Run("stty -echo") // disable echo inside the container
	require.NoError(t, err)
	require.Equal(t, []string{"stty -echo"}, output)

	output, err = shell.Run("echo $TESTVAR")
	require.NoError(t, err)
	require.Equal(t, []string{"foobar"}, output)
}

package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new module-name",
	Short: "generate a project in current dir",
	Long:  "generate a project in current dir, do not use golang keywords or official package name as your module name.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
		}
		for _, a := range args {
			if strings.Contains(a, " ") {
				return fmt.Errorf("no spacese in module name")
			}
		}
		return nil
	},
	Run: newProject,
}

func initNewCmd() {
	rootCmd.AddCommand(newCmd)
}

func newProject(cmd *cobra.Command, args []string) {
	path := "./" + args[0]
	if _, err := os.Stat(path); err == nil || os.IsExist(err) {
		fmt.Println(path, "already exist.")
		return
	}

	if err := execShell("git", "clone", "--progress", "git@git.github.com:22km/katana-demo.git", args[0]); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := execShell("sh", path+"/replace.sh", args[0]); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func execShell(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	out, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	reader := bufio.NewReader(out)

	if err := cmd.Start(); err != nil {
		return err
	}

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Println(string(line))
	}

	return cmd.Wait()
}

package cmd

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	configFlag string
)

// appsCreateCmd represents the create command
var appsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		remote, err := ensureGitRemote()
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			return errors.WithStack(err)
		}
		fmt.Printf("Selected remote: %s", remote)
		return nil
	},
}

func init() {
	appsCmd.AddCommand(appsCreateCmd)
	appsCreateCmd.Flags().StringVarP(&configFlag, "config", "c", "", "Path of the bitrise.yml config file")
}

// AppsCreateResponseModel ...
type AppsCreateResponseModel struct {
	// TODO
}

func gitRemotes() ([]string, error) {
	cmd := command.New("git", "remote", "-v")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, err
	}

	out = strings.Replace(strings.Replace(out, "(fetch)", "", -1), "(push)", "", -1)
	return removeDuplicates(strings.Split(out, "\n")), nil
}

func ensureGitRemote() (string, error) {
	remotes, err := gitRemotes()
	if err != nil {
		return "", err
	}

	switch len(remotes) {
	case 0:
		return "", fmt.Errorf("failed to find the git remote")
	case 1:
		return remotes[0], nil
	default:
		fmt.Println()
		remote, err := goinp.SelectFromStringsWithDefault("Select the app which you want to upload the privisioning profiles", 1, remotes)
		if err != nil {
			return "", err
		}
		return remote, nil
	}
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	var result []string

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

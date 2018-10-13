package cmd

import (
	"fmt"
	"strings"

	"github.com/bitrise-core/bitrise-plugins-io/services"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	configFlag   string
	remoteFlag   string
	isPublicFlag string
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
		remote, err := ensureGitRemote(cmd)
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Printf("Selected remote: %s\n", remote)

		repoName, err := ensureRepoName(cmd)
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Printf("Repo name: %s\n", repoName)

		priv, err := appPrivacy(cmd)
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Printf("App privacy: %s\n", priv)

		params := map[string]interface{}{
			"is_public":     Privacy(isPublicFlag).isPublic(),
			"repo_url":      remote,
			"type":          "git",
			"git_repo_slug": repoName,
		}

		return errors.WithStack(create(params))
	},
}

func init() {
	appsCmd.AddCommand(appsCreateCmd)
	appsCreateCmd.Flags().StringVar(&configFlag, "config", "", "Path of the bitrise.yml config file")
	appsCreateCmd.Flags().StringVar(&configFlag, "remote", "", "Git remote URL of your repository")
	appsCreateCmd.Flags().StringVar(&isPublicFlag, "privacy", "", "Privacy of the app [private/public]")
}

// AppsCreateResponseModel ...
type AppsCreateResponseModel struct {
	Status string `json:"status"`
	Slug   string `json:"slug"`
}

// Pretty ...
func (respModel *AppsCreateResponseModel) Pretty() string {
	return fmt.Sprintf("%s / %s\n", respModel.Status, respModel.Status)
}

// Privacy ...
type Privacy string

// const ...
const (
	Private Privacy = "private"
	Public  Privacy = "public"
)

func (p Privacy) isPublic() bool {
	return p == Public
}

func gitRemotes() ([]string, error) {
	cmd := command.New("git", "config", "--get", "remote.origin.url")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	return strings.Split(out, "\n"), err
}

func ensureGitRemote(cmd *cobra.Command) (string, error) {
	if cmd.Flags().Changed("remote") {
		return cmd.Flag("remote").Value.String(), nil
	}

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

func ensureRepoName(c *cobra.Command) (string, error) {
	if c.Flags().Changed("repo-slug") {
		return c.Flag("repo-slug").Value.String(), nil
	}

	cmd := command.New("bash", "-c", "basename `git rev-parse --show-toplevel`")
	return cmd.RunAndReturnTrimmedCombinedOutput()
}

func appPrivacy(cmd *cobra.Command) (Privacy, error) {
	if cmd.Flags().Changed("privacy") {
		privacy := cmd.Flag("privacy").Value.String()
		if privacy == string(Private) || privacy == string(Public) {
			return Privacy(privacy), nil
		}
	}
	fmt.Println()
	fmt.Println()

	opts := []string{string(Private), string(Public)}
	privacy, err := goinp.SelectFromStringsWithDefault("Set privacy of the app", 1, opts)
	return Privacy(privacy), err
}

func create(params map[string]interface{}) error {
	response, err := services.RegisterApp(params)
	if err != nil {
		return errors.WithStack(err)
	}

	if response.Error != "" {
		return NewRequestFailedError(response)
	}

	return errors.WithStack(printOutputWithPrettyFormatter(response.Data, formatFlag != "json", &AppsCreateResponseModel{}))
}

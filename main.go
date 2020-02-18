package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var cmd = &cobra.Command{
		RunE:    run,
		Short:   "Sending message to Rocket.Chat using Rest API",
		Use:     "drone-rocketchat",
		Version: "0.0.1",
	}

	viper.AutomaticEnv()
	flags := cmd.Flags()
	flags.String("url", "", "Rocket.Chat url")
	viper.BindPFlag("url", flags.Lookup("url"))
	flags.Bool("insecure", false, "Skip ssl CA verification")
	viper.BindPFlag("insecure", flags.Lookup("insecure"))
	flags.String("trusted-ca", "", "Your own trusted CA")
	viper.BindPFlag("trusted-ca", flags.Lookup("trusted-ca"))
	flags.String("userId", "", "Rocket.chat user id")
	viper.BindPFlag("userId", flags.Lookup("userId"))
	flags.String("userToken", "", "Rocket.chat user API token")
	viper.BindPFlag("userToken", flags.Lookup("userToken"))
	flags.String("channel", "", "Rocket.chat channel name")
	viper.BindPFlag("channel", flags.Lookup("channel"))
	flags.String("message", "", "The message contents (up to 2000 characters)")
	viper.BindPFlag("message", flags.Lookup("message"))
	flags.String("avatar-url", "", "Override the default avatar of user")
	viper.BindPFlag("avatar-url", flags.Lookup("avatar-url"))
	flags.Bool("drone", false, "Environment is drone")
	viper.BindPFlag("drone", flags.Lookup("drone"))
	flags.String("repo", "", "Repository owner and repository name")
	viper.BindPFlag("repo", flags.Lookup("repo"))
	flags.String("repo.namespace", "", "Repository namespace")
	viper.BindPFlag("repo.namespace", flags.Lookup("repo.namespace"))
	flags.String("repo.name", "", "Repository name")
	viper.BindPFlag("repo.name", flags.Lookup("repo.name"))
	flags.String("commit.sha", "", "Git commit hash")
	viper.BindPFlag("commit.sha", flags.Lookup("commit.sha"))
	flags.String("commit.ref", "", "Git commit ref")
	viper.BindPFlag("commit.ref", flags.Lookup("commit.ref"))
	flags.String("commit.branch", "", "Git commit branch")
	viper.BindPFlag("commit.branch", flags.Lookup("commit.branch"))
	flags.String("commit.author", "", "Git commit author")
	viper.BindPFlag("commit.author", flags.Lookup("commit.author"))
	flags.String("commit.author.email", "", "Git commit author email")
	viper.BindPFlag("commit.author.email", flags.Lookup("commit.author.email"))
	flags.String("commit.author.avatar", "", "Git commit author avatar")
	viper.BindPFlag("commit.author.avatar", flags.Lookup("commit.author.avatar"))
	flags.String("commit.message", "", "Git commit message")
	viper.BindPFlag("commit.message", flags.Lookup("commit.message"))
	flags.String("build.event", "", "Build event")
	viper.BindPFlag("build.event", flags.Lookup("build.event"))
	flags.Int("build.number", 0, "Build number")
	viper.BindPFlag("build.number", flags.Lookup("build.number"))
	flags.String("build.status", "", "Build status")
	viper.BindPFlag("build.status", flags.Lookup("build.status"))
	flags.String("build.link", "", "Build link")
	viper.BindPFlag("build.link", flags.Lookup("build.link"))
	flags.String("build.tag", "", "Build tag")
	viper.BindPFlag("build.tag", flags.Lookup("build.tag"))
	flags.Float64("job.started", 0, "Job started")
	viper.BindPFlag("job.started", flags.Lookup("job.started"))
	flags.Float64("job.finished", 0, "Job finished")
	viper.BindPFlag("job.finished", flags.Lookup("job.finished"))
	flags.String("env-file", "", "Source environment file: file.env, file.yml, ...")
	viper.BindPFlag("env-file", flags.Lookup("env-file"))
	flags.Bool("github", false, "Environment is GitHub Action")
	viper.BindPFlag("drone", flags.Lookup("drone"))
	flags.String("github.workflow", "", "The name of the Github workflow")
	viper.BindPFlag("github.workflow", flags.Lookup("github.workflow"))
	flags.String("github.action", "", "The name of the Github action")
	viper.BindPFlag("github.action", flags.Lookup("github.action"))
	flags.String("github.event.name", "", "The name of the Github event.name")
	viper.BindPFlag("github.event.name", flags.Lookup("github.event.name"))
	flags.String("github.event.path", "", "The name of the Github event.path")
	viper.BindPFlag("github.event.path", flags.Lookup("github.event.path"))
	flags.String("github.event.workspace", "", "The name of the Github event.workspace")
	viper.BindPFlag("github.event.workspace", flags.Lookup("github.event.workspace"))

	cmd.Execute()
}

func viperGetStrings(strList []string) string {
	for _, s := range strList {
		if viper.IsSet(s) {
			return viper.GetString(s)
		}
	}
	return ""
}

func viperGetBool(strList []string) bool {
	for _, s := range strList {
		if viper.IsSet(s) {
			b, err := strconv.ParseBool(viper.GetString(s))
			if err != nil {
				log.Fatal(fmt.Sprintf("%s is not parsable as boolean", s))
			}
			return b
		}
	}
	return false
}

func viperGetFloat64(strList []string) float64 {
	for _, s := range strList {
		if viper.IsSet(s) {
			f, err := strconv.ParseFloat(viper.GetString(s), 64)
			if err != nil {
				log.Fatal(fmt.Sprintf("%s is not parsable as float", s))
			}
			return f
		}
	}
	return 0
}

func viperGetInt64(strList []string) int64 {
	for _, s := range strList {
		if viper.IsSet(s) {
			i, err := strconv.ParseInt(viper.GetString(s), 10, 64)
			if err != nil {
				log.Fatal(fmt.Sprintf("%s is not parsable an int", s))
			}
			return i
		}
	}
	return 0
}

func run(cmd *cobra.Command, args []string) error {

	// did we get an environment file
	if viper.IsSet("env-file") && viper.GetString("env-file") != "" {
		viper.SetConfigName(viper.GetString("env-file"))
		viper.AddConfigPath(".")
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}

	// populate plugin structs
	plugin := Plugin{
		GitHub: GitHub{
			Workflow:  viperGetStrings([]string{"github_workflow"}),
			Workspace: viperGetStrings([]string{"github_workspace"}),
			Action:    viperGetStrings([]string{"github_action"}),
			EventName: viperGetStrings([]string{"github_event_name"}),
			EventPath: viperGetStrings([]string{"github_event_path"}),
		},
		Repo: Repo{
			FullName:  viperGetStrings([]string{"drone_repo", "github_repository"}),
			Namespace: viperGetStrings([]string{"drone_repo_owner", "drone_repo_namespace", "github_actor"}),
			Name:      viperGetStrings([]string{"drone_repo_name"}),
		},
		Build: Build{
			Tag:      viperGetStrings([]string{"drone_build_tag"}),
			Number:   viperGetInt64([]string{"drone_build_number"}),
			Event:    viperGetStrings([]string{"drone_build_event"}),
			Status:   viperGetStrings([]string{"drone_build_status"}),
			Commit:   viperGetStrings([]string{"drone_commit_sha", "github_sha"}),
			RefSpec:  viperGetStrings([]string{"drone_commit_ref", "github_ref"}),
			Branch:   viperGetStrings([]string{"drone_commit_branch"}),
			Author:   viperGetStrings([]string{"drone_commit_author"}),
			Email:    viperGetStrings([]string{"drone_commit_author_email"}),
			Avatar:   viperGetStrings([]string{"drone_commit_author_avatar"}),
			Message:  viperGetStrings([]string{"drone_commit_message"}),
			Link:     viperGetStrings([]string{"drone_build_link"}),
			Started:  viperGetFloat64([]string{"drone_job_started"}),
			Finished: viperGetFloat64([]string{"drone_job_finished"}),
		},
		Config: Config{
			URL:       viperGetStrings([]string{"plugin_url", "url"}),
			Insecure:  viperGetBool([]string{"plugin_insecure", "insecure"}),
			TrustedCA: viperGetStrings([]string{"plugin_trustedca", "trustedca"}),
			UserID:    viperGetStrings([]string{"plugin_userId", "userId"}),
			Token:     viperGetStrings([]string{"plugin_userToken", "userToken"}),
			Message:   viperGetStrings([]string{"plugin_message", "message"}),
			Drone:     viperGetBool([]string{"drone"}),
			GitHub:    viperGetBool([]string{"plugin_github", "github"}),
		},
		Payload: Payload{
			Avatar:  viperGetStrings([]string{"plugin_avatar_url", "avatar_url"}),
			Channel: viperGetStrings([]string{"plugin_channel", "channel"}),
		},
	}

	//
	// fmt.Printf("url: %q\n", viperGetStrings([]string{"plugin_url", "url"}))
	// fmt.Printf("insecure: %q\n", viperGetStrings([]string{"insecure"}))
	// fmt.Printf("trusted-ca: %q\n", viperGetStrings([]string{"trusted-ca"}))
	// fmt.Printf("userId: %q\n", viperGetStrings([]string{"plugin_userId", "userId"}))
	// fmt.Printf("userToken: %q\n", viperGetStrings([]string{"plugin_userToken", "userToken"}))
	// fmt.Printf("channel: %q\n", viperGetStrings([]string{"plugin_channel", "channel"}))
	fmt.Printf("message: %q\n", viperGetStrings([]string{"plugin_message", "message"}))
	// fmt.Printf("avatar: %q\n", viperGetStrings([]string{"plugin_avatar_url", "avatar_url"}))
	// fmt.Printf("drone: %q\n", viperGetStrings([]string{"drone"}))
	// fmt.Printf("repo: %q\n", viperGetStrings([]string{"drone_repo", "github_repository"}))
	// fmt.Printf("repo.namespace: %q\n", viperGetStrings([]string{"drone_repo_owner", "drone_repo_namespace", "github_actor"}))
	// fmt.Printf("repo.name: %q\n", viperGetStrings([]string{"drone_repo_name"}))
	// fmt.Printf("commit.sha: %q\n", viperGetStrings([]string{"drone_commit_sha", "github_sha"}))
	// fmt.Printf("commit.ref: %q\n", viperGetStrings([]string{"drone_commit_ref", "github_ref"}))
	// fmt.Printf("commit.branch: %q\n", viperGetStrings([]string{"drone_commit_branch"}))
	// fmt.Printf("commit.author: %q\n", viperGetStrings([]string{"drone_commit_author"}))
	// fmt.Printf("commit.author.email: %q\n", viperGetStrings([]string{"drone_commit_author_email"}))
	// fmt.Printf("commit.author.avatar: %q\n", viperGetStrings([]string{"drone_commit_author_avatar"}))
	// fmt.Printf("commit.message: %q\n", viperGetStrings([]string{"drone_commit_message"}))
	// fmt.Printf("build.event: %q\n", viperGetStrings([]string{"drone_build_event"}))
	// fmt.Printf("build.number: %q\n", viperGetStrings([]string{"drone_build_number"}))
	// fmt.Printf("build.status: %q\n", viperGetStrings([]string{"drone_build_status"}))
	// fmt.Printf("build.link: %q\n", viperGetStrings([]string{"drone_build_link"}))
	// fmt.Printf("build.tag: %q\n", viperGetStrings([]string{"drone_build_tag"}))
	// fmt.Printf("job.started: %q\n", viperGetStrings([]string{"drone_job_started"}))
	// fmt.Printf("job.finished: %q\n", viperGetStrings([]string{"drone_job_finished"}))
	// fmt.Printf("github: %q\n", viperGetStrings([]string{"plugin_github", "github"}))
	// fmt.Printf("github.workflow: %q\n", viperGetStrings([]string{"github_workflow"}))
	// fmt.Printf("github.action: %q\n", viperGetStrings([]string{"github_action"}))
	// fmt.Printf("github.event.name: %q\n", viperGetStrings([]string{"github_event_name"}))
	// fmt.Printf("github.event.path: %q\n", viperGetStrings([]string{"github_event_path"}))
	// fmt.Printf("github.workspace: %q\n", viperGetStrings([]string{"github_workspace"}))

	return plugin.Exec()
	// return nil
}

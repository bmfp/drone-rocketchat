package main

import (
	"encoding/json"
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
	flags.String("custom-msg-fileds", "", "Custom fields, json dictionnary")
	viper.BindPFlag("custom-msg-fileds", flags.Lookup("custom-msg-fileds"))
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
	viper.BindPFlag("github", flags.Lookup("github"))
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
	flags.Bool("debug", false, "Debug")
	viper.BindPFlag("debug", flags.Lookup("debug"))

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

func viperGetJSON(strList []string) map[string]interface{} {
	ret := make(map[string]interface{})
	var v interface{}
	for _, s := range strList {
		if viper.IsSet(s) {
			json.Unmarshal([]byte(viper.GetString(s)), &v)
			ret = v.(map[string]interface{})

			return ret
		}
	}
	return ret
}

func run(cmd *cobra.Command, args []string) error {

	// did we get an environment file
	var gotEnvFile bool = false
	if viper.IsSet("env-file") && viper.GetString("env-file") != "" {
		viper.SetConfigName(viper.GetString("env-file"))
		gotEnvFile = true
	} else if viperGetStrings([]string{"plugin_env_file", "env_file"}) != "" {
		viper.SetConfigName(viperGetStrings([]string{"plugin_env_file", "env_file"}))
		gotEnvFile = true
	}
	if gotEnvFile {
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
			EnvFile:   viperGetStrings([]string{"plugin_env_file", "env_file"}),
			Debug:     viperGetBool([]string{"plugin_debug", "debug"}),
		},
		Payload: Payload{
			Avatar:          viperGetStrings([]string{"plugin_avatar_url", "avatar_url"}),
			Channel:         viperGetStrings([]string{"plugin_channel", "channel"}),
			CustomMSgFields: viperGetJSON([]string{"custom_msg_fields"}),
		},
	}

	if plugin.Config.Debug {
		fmt.Printf("url: %q\n", plugin.Config.URL)
		fmt.Printf("insecure: %q\n", plugin.Config.Insecure)
		fmt.Printf("trusted-ca: %q\n", plugin.Config.TrustedCA)
		fmt.Printf("userId: %q\n", plugin.Config.UserID)
		fmt.Printf("userToken: %q\n", plugin.Config.Token)
		fmt.Printf("channel: %q\n", plugin.Payload.Channel)
		fmt.Printf("message: %q\n", plugin.Config.Message)
		fmt.Printf("customfields: %q\n", plugin.Payload.CustomMSgFields)
		fmt.Printf("avatar: %q\n", plugin.Payload.Avatar)
		fmt.Printf("drone: %q\n", plugin.Config.Drone)
		fmt.Printf("repo: %q\n", plugin.Repo.FullName)
		fmt.Printf("repo.namespace: %q\n", plugin.Repo.Namespace)
		fmt.Printf("repo.name: %q\n", plugin.Repo.Name)
		fmt.Printf("commit.sha: %q\n", plugin.Build.Commit)
		fmt.Printf("commit.ref: %q\n", plugin.Build.RefSpec)
		fmt.Printf("commit.branch: %q\n", plugin.Build.Branch)
		fmt.Printf("commit.author: %q\n", plugin.Build.Author)
		fmt.Printf("commit.author.email: %q\n", plugin.Build.Email)
		fmt.Printf("commit.author.avatar: %q\n", plugin.Build.Avatar)
		fmt.Printf("commit.message: %q\n", plugin.Build.Message)
		fmt.Printf("build.event: %q\n", plugin.Build.Event)
		fmt.Printf("build.number: %q\n", plugin.Build.Number)
		fmt.Printf("build.status: %q\n", plugin.Build.Status)
		fmt.Printf("build.link: %q\n", plugin.Build.Link)
		fmt.Printf("build.tag: %q\n", plugin.Build.Tag)
		fmt.Printf("job.started: %q\n", plugin.Build.Started)
		fmt.Printf("job.finished: %q\n", plugin.Build.Finished)
		fmt.Printf("github: %q\n", plugin.Config.GitHub)
		fmt.Printf("github.workflow: %q\n", plugin.GitHub.Workflow)
		fmt.Printf("github.action: %q\n", plugin.GitHub.Action)
		fmt.Printf("github.event.name: %q\n", plugin.GitHub.EventName)
		fmt.Printf("github.event.path: %q\n", plugin.GitHub.EventPath)
		fmt.Printf("github.workspace: %q\n", plugin.GitHub.Workspace)
	}

	err := plugin.Exec()
	if err != nil {
		log.Fatal(err)
	}
	return err
}

package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pennywisdom/commitizen-go/commit"
	"github.com/pennywisdom/commitizen-go/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute() error {
	return rootCmd.Execute()
}

var (
	all     bool
	sign    bool
	debug   bool
	rootCmd = &cobra.Command{
		Use:  "commitizen-go",
		Long: `Command line utility to standardize git commit messages, golang version.`,
		Run:  RootCmd,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "tell the command to automatically stage files that have been modified and deleted, but new files you have not told Git about are not affected")
	rootCmd.Flags().BoolVarP(&sign, "sign", "S", false, "sign the commit and add sign-off line")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode, output debug info to debug.log")

	// viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.AddCommand(VersionCmd)
	rootCmd.AddCommand(InstallCmd)
}

func initConfig() {
	if !debug {
		log.SetOutput(io.Discard)
	} else {
		f, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		// defer f.Close()
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.SetOutput(f)
	}

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Printf("Get home dir failed, err=%v\n", err)
		os.Exit(1)
	}
	workingTreeRoot, err := git.WorkingTreeRoot()
	if err != nil || workingTreeRoot == "" {
		log.Printf("current directory is not a git working tree\n")
		os.Exit(1)
	}

	// Search config in repository working tree root directory and home directory with name ".git-czrc" or ".git-czrc.json".
	viper.AddConfigPath(workingTreeRoot)
	viper.AddConfigPath(home)
	viper.SetConfigName(".git-czrc")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Printf("can not find config file\n")
		} else {
			// Config file was found but another error was produced
			log.Printf("read config failed, err=%v\n", err)
		}
	} else {
		log.Printf("read config success\n")
	}

}

func RootCmd(command *cobra.Command, args []string) {
	// exit if not git repo
	if _, err := git.IsCurrentDirectoryGitRepo(); err != nil {
		fmt.Print(err)
		return
	}

	if message, err := commit.FillOutForm(); err == nil {
		// do git commit
		result, err := git.CommitMessage(message, all, sign)
		if err != nil {
			log.Printf("run git commit failed, err=%v\n", err)
			log.Printf("commit message is: \n\n%s\n\n", string(message))
		}
		fmt.Print(string(result))
	}
}

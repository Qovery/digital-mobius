package run

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(GetCurrentVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func GetCurrentVersion() string {
	return "0.1.1" // ci-version-check
}

func GetAscii() string {
	return "\n     _ _       _ _        _                      _     _           \n  __| (_) __ _(_) |_ __ _| |     _ __ ___   ___ | |__ (_)_   _ ___ \n / _` | |/ _` | | __/ _` | |    | '_ ` _ \\ / _ \\| '_ \\| | | | / __|\n| (_| | | (_| | | || (_| | |    | | | | | | (_) | |_) | | |_| \\__ \\\n \\__,_|_|\\__, |_|\\__\\__,_|_|    |_| |_| |_|\\___/|_.__/|_|\\__,_|___/\n         |___/                                                     "
}
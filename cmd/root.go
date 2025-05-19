package cmd

import (
	"fmt"
	"goctl-gin/generator"
	"goctl-gin/prepare"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "gengin",
		Short:   "生成基于 GIN 框架的 WEB 服务的极简项目结构",
		Example: "gengin --dir=. some.api",
		Args:    cobra.ExactValidArgs(1),
		RunE:    GenGinCode,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&prepare.OutputDir, "dir", ".", "生成项目目录")
}

func GenGinCode(cmd *cobra.Command, args []string) error {
	prepare.ApiFile = args[0]
	prepare.Setup()
	Must(generator.GenRoutes())
	return nil
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

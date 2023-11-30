package main

import (
	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/config"
	"github.com/starudream/go-lib/core/v2/utils/osutil"
)

var rootCmd = cobra.NewRootCommand(func(c *cobra.Command) {
	c.Use = "aliyunxiao-generator"

	c.PersistentFlags().String("ticket", "", "aliyun ticket")
	osutil.Must0(c.MarkPersistentFlagRequired("ticket"))

	c.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		config.LoadFlags(c.PersistentFlags())
	}
})

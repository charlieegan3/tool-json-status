package main

import (
	"context"
	"log"
	"os"

	"github.com/charlieegan3/toolbelt/pkg/database"
	"github.com/charlieegan3/toolbelt/pkg/tool"
	"github.com/spf13/viper"

	jsonStatusTool "github.com/charlieegan3/tool-json-status/pkg/tool"
)

func main() {
	var err error
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	params := viper.GetStringMapString("database.params")
	connectionString := viper.GetString("database.connectionString")
	db, err := database.Init(connectionString, params, params["dbname"], false)
	if err != nil {
		log.Fatalf("failed to init DB: %s", err)
	}

	tb := tool.NewBelt()
	tb.SetConfig(viper.GetStringMap("tools"))
	tb.SetDatabase(db)

	jst := &jsonStatusTool.JSONStatus{}
	err = tb.AddTool(jst)
	if err != nil {
		log.Fatalf("failed to add tool: %v", err)
	}

	if len(os.Args) == 2 {
		jobs, err := jst.Jobs()
		if err != nil {
			log.Fatalf("failed to get jobs: %v", err)
		}

		err = jobs[1].Run(context.Background())
		if err != nil {
			log.Fatalf("failed to run job: %v", err)
		}

		return
	}

	// go tb.RunJobs(context.Background())

	tb.RunServer(context.Background(), "0.0.0.0", "3000")
}

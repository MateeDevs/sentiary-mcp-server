package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MateeDevs/sentiary-tools/internal/buildinfo"
	"github.com/MateeDevs/sentiary-tools/internal/frontend/cli"
	"github.com/MateeDevs/sentiary-tools/internal/sentiary"
)

func main() {
	if buildinfo.IsVersionCommand(os.Args[1:]) {
		fmt.Println(buildinfo.String())
		return
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	api := sentiary.NewClient(sentiary.ConfigFromEnv(), httpClient)
	app := cli.New(api)
	if err := app.Run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if cli.IsUsageError(err) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}

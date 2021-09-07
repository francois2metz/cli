package crontasks

import (
	"os"

	"github.com/Scalingo/cli/config"
	httpclient "github.com/Scalingo/go-scalingo/v4/http"
	"github.com/Scalingo/go-utils/errors"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List(app string) error {
	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	cronTasks, err := client.CronTasksGet(app)
	if err != nil {
		rootError := errors.ErrgoRoot(err)
		if !httpclient.IsRequestFailedError(rootError) || rootError.(*httpclient.RequestFailedError).Code != 404 {
			return errgo.Notef(err, "fail to get cron tasks")
		}

		// A 404 only means there is no cron task configured on the application. In this case, we want to display an empty table.
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"Command", "Size"})

	for _, job := range cronTasks.Jobs {
		t.Append([]string{job.Command, job.Size})
	}
	t.Render()

	return nil
}
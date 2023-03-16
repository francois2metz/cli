package apps

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	scalingo "github.com/Scalingo/go-scalingo/v6"
	errors "github.com/Scalingo/go-utils/errors/v2"
)

func handleOperation(ctx context.Context, app string, res *http.Response) error {
	operationURL := res.Header.Get("Location")
	return handleOperationWithURL(ctx, app, operationURL)
}

func handleOperationWithURL(ctx context.Context, app string, operationURL string, containerLabel ...string) error {
	opURL, err := url.Parse(operationURL)
	if err != nil {
		return errors.Notef(ctx, err, "parse url of operation")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Notef(ctx, err, "get Scalingo client")
	}

	var op *scalingo.Operation

	opID := filepath.Base(opURL.Path)
	done := make(chan struct{})
	errs := make(chan error)
	defer close(done)
	defer close(errs)

	op, err = c.OperationsShow(ctx, app, opID)
	if err != nil {
		return errors.Notef(ctx, err, "get operation %v", opID)
	}

	go func() {
		for {
			op, err = c.OperationsShow(ctx, app, opID)
			if err != nil {
				errs <- err
				break
			}

			if op.Status == "done" || op.Status == "error" {
				done <- struct{}{}
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()

	if op.Type == scalingo.OperationTypeStartOneOff {
		fmt.Printf("-----> Starting container %v   ", containerLabel)
	} else {
		fmt.Print("Status:  ")
	}
	spinner := io.NewSpinner(os.Stderr)
	go spinner.Start()
	defer spinner.Stop()

	for {
		select {
		case err := <-errs:
			return errors.Notef(ctx, err, "get operation %v", op.ID)
		case <-done:
			if op.Status == "done" {
				fmt.Printf("\bDone in %.3f seconds\n", op.ElapsedDuration())
				return nil
			} else if op.Status == "error" {
				fmt.Printf("\bOperation '%s' failed, an error occurred: %v\n", op.Type, op.Error)
				return errors.Newf(ctx, "operation %v failed", op.ID)
			}
		}
	}
}

func GetAttachURLFromOperationWithURL(ctx context.Context, app string, operationURL string) (string, error) {
	opURL, err := url.Parse(operationURL)
	if err != nil {
		return "", errors.Notef(ctx, err, "parse url of operation")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return "", errors.Notef(ctx, err, "get Scalingo client")
	}

	var operation *scalingo.Operation
	opID := filepath.Base(opURL.Path)
	operation, err = c.OperationsShow(ctx, app, opID)
	if err != nil {
		return "", errors.Notef(ctx, err, "get operation %v", opID)
	}

	return operation.StartOneOffData.AttachURL, nil
}

package cmd

import (
	"context"
	"io"
)

func HandleHttp(ctx context.Context, w io.Writer, args []string) error {
	hc, err := flagConfig(w, args)
	if err != nil {
		return err
	}

	client := createHTTPClientWithTimeout()
	if hc.disableRedirect {
		client.CheckRedirect = redirectPolicyFunc
	}
	body, err := fetchRemoteResource(ctx, client, hc)
	if err != nil {
		return err
	}

	var output Output
	if hc.filePath != "" {
		output = &FileOutput{filePath: hc.filePath}
	} else {
		output = &ConsoleOutput{wirter: w}
	}

	err = output.Write(body)
	if err != nil {
		return err
	}

	return nil
}

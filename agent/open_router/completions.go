package open_router

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/242617/other/agent"
)

func (p *OpenRouter) completions(ctx context.Context, req request, fn responseFunc, tools agent.Tools) error {
	client := resty.New().
		SetBaseURL("https://openrouter.ai").
		SetAuthToken(p.token).
		SetHeader("User-Agent", p.name)

	// Prepare request
	body := request{
		Model: req.Model,
		Tools: tools.Info(),
	}
	body.Messages = req.Messages[:]

	var response response
	res, err := client.NewRequest().
		SetDebug(p.debug).
		SetContext(ctx).
		SetBody(body).
		SetResult(&response).
		Post("api/v1/chat/completions")
	if err != nil {
		slog.Error("client new request", "err", err, "body", body)
		return errors.Wrap(err, "client new request")
	}
	if res.StatusCode() != http.StatusOK {
		b, err := httputil.DumpResponse(res.RawResponse, true)
		slog.Error("unexpected status code", "err", err, "status code", res.StatusCode(), "body", string(b))
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return fn(response)
}

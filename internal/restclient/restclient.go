package restclient

import (
	"context"
	"net/http"
	"strings"
)

type RESTClient interface {
	Do(ctx context.Context, cli *http.Client, streams IOStreams) error
}

type RequestOptions struct {
	BaseURL string
	Method  string
	Path    string
	Params  []string
	Headers []string
	Streams IOStreams
}

func New(opts RequestOptions) RESTClient {
	rc := &restClientImpl{
		baseURL: opts.BaseURL,
		urlPath: opts.Path,
		verb:    opts.Method,
	}

	if tot := len(opts.Headers); tot > 0 {
		rc.requestHeaders = make([]string, tot)
		copy(rc.requestHeaders, opts.Headers)
	}

	if tot := len(opts.Params); tot > 0 {
		rc.requestParams = make([]string, tot)
		copy(rc.requestParams, opts.Params)
	}

	if rc.verb == "" {
		rc.verb = http.MethodGet
	}

	return rc
}

var _ RESTClient = (*restClientImpl)(nil)

type restClientImpl struct {
	baseURL        string
	urlPath        string
	verb           string
	requestParams  []string
	requestHeaders []string
}

func (hc *restClientImpl) Do(ctx context.Context, cli *http.Client, streams IOStreams) error {
	uri, err := composeURL(hc.baseURL, hc.urlPath, hc.requestParams...)
	if err != nil {
		return err
	}

	method := strings.ToUpper(strings.TrimSpace(hc.verb))
	if method == "" {
		method = http.MethodGet
	}

	call, err := http.NewRequestWithContext(ctx, method, uri, streams.In)
	if err != nil {
		return err
	}

	setHeaders(call, hc.requestHeaders...)

	respo, err := cli.Do(call)
	if err != nil {
		return err
	}
	defer respo.Body.Close()

	return dumpResponse(respo, streams.Out, streams.Err)
}

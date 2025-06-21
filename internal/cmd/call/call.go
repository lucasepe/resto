package call

import (
	"context"
	"fmt"

	"io"
	"log"
	"os"

	"github.com/lucasepe/resto/internal/restclient"
	getoptutil "github.com/lucasepe/resto/internal/util/getopt"
	ioutil "github.com/lucasepe/resto/internal/util/io"
	"github.com/lucasepe/resto/internal/util/retry"
	"github.com/lucasepe/x/getopt"
	"github.com/lucasepe/x/text/conv"
)

func Do(args []string) error {
	extras, opts, err := getopt.GetOpt(args,
		"X:H:f:v",
		[]string{
			"proxy-url",
			"max-attempts=",
			"max-delay=",
			"ca-cert=",
			"cert=",
			"cert-key=",
			"file=",
			"header=",
			"insecure=",
			"initial-delay=",
			"max-jitter=",
			"password=",
			"request=",
			"token=",
			"until=",
			"username=",
			"verbose",
		},
	)
	if err != nil {
		return err
	}

	if len(extras) < 1 {
		return fmt.Errorf("missing request uri")
	}

	reqOpts, err := requestOptions(extras, opts)
	if err != nil {
		return err
	}

	streams := ioStreams(opts)
	filename := getoptutil.OptVal(opts, []string{"-f", "--file"})
	in, close, err := ioutil.FileOrStdin(filename)
	if err == nil {
		streams.In = in
	}
	defer close()

	expr := getoptutil.EnvOrOptVal("UNTIL", opts, []string{"--until"})

	cfg := restClientConfig(opts)
	if cfg.ServerURL == "" {
		cfg.ServerURL = reqOpts.BaseURL
	}

	if cfg.Verbose && expr != "" {
		log.Printf("jq expression: %q\n", expr)
	}

	retryOpts := retryOptions(opts)

	cli, err := restclient.HTTPClientForConfig(cfg)
	if err != nil {
		return err
	}
	cli.Transport = retry.NewRoundTripperWithEval(cli.Transport, expr,
		retry.Jittered(retryOpts.MaxJitter),
		retry.NewRetrier(retryOpts),
	)

	return restclient.New(reqOpts).Do(context.Background(), cli, streams)
}

func restClientConfig(opts []getopt.OptArg) restclient.Config {
	cfg := restclient.ConfigFromEnv()

	proxyUrl := getoptutil.OptVal(opts, []string{"--proxy-url"})
	if proxyUrl != "" {
		cfg.ProxyURL = proxyUrl
	}

	username := getoptutil.OptVal(opts, []string{"--username"})
	password := getoptutil.OptVal(opts, []string{"--password"})
	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}

	token := getoptutil.OptVal(opts, []string{"--token"})
	if token != "" {
		cfg.Token = token
	}

	clientCert := getoptutil.OptVal(opts, []string{"--cert"})
	if clientCert != "" {
		cfg.ClientCertificateData = clientCert
	}

	clientKey := getoptutil.OptVal(opts, []string{"--cert-key"})
	if clientKey != "" {
		cfg.ClientKeyData = clientKey
	}

	caCert := getoptutil.OptVal(opts, []string{"--ca-cert"})
	if caCert != "" {
		cfg.CertificateAuthorityData = caCert
	}

	cfg.Insecure = getoptutil.HasOpt(opts, []string{"--insecure"})
	cfg.Verbose = getoptutil.HasOpt(opts, []string{"-v", "--verbose"})

	return cfg
}

func retryOptions(opts []getopt.OptArg) retry.RetryOptions {
	res := retry.OptionsFromEnv()

	val := getoptutil.OptVal(opts, []string{"--initial-delay"})
	if val != "" {
		res.InitialDelay = conv.Duration(val, res.InitialDelay)
	}

	val = getoptutil.OptVal(opts, []string{"--max-delay"})
	if val != "" {
		res.MaxDelay = conv.Duration(val, res.MaxDelay)
	}

	val = getoptutil.OptVal(opts, []string{"--max-attempts"})
	if val != "" {
		res.MaxAttempts = conv.Int(val, res.MaxAttempts)
	}

	val = getoptutil.OptVal(opts, []string{"--max-jitter"})
	if val != "" {
		res.MaxJitter = conv.Duration(val, res.MaxJitter)
	}

	return res
}

func requestOptions(extras []string, opts []getopt.OptArg) (restclient.RequestOptions, error) {
	baseURL, path, params, err := reverseURL(extras[0])
	if err != nil {
		return restclient.RequestOptions{}, err
	}

	return restclient.RequestOptions{
		BaseURL: baseURL,
		Method:  getoptutil.OptVal(opts, []string{"-X", "--request"}),
		Path:    path,
		Headers: getoptutil.AllOptArgs(opts, []string{"-H", "--header"}),
		Params:  params,
	}, nil
}

func ioStreams(opts []getopt.OptArg) restclient.IOStreams {
	ios := restclient.IOStreams{
		Out: os.Stdout,
		Err: os.Stderr,
	}

	verbose := getoptutil.HasOpt(opts, []string{"-v", "--verbose"})
	if verbose {
		ios.Out = io.Discard
		ios.Err = io.Discard
	}

	return ios
}

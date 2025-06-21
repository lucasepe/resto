package cmd

import (
	"fmt"
	"io"
	"strings"

	xtext "github.com/lucasepe/x/text"
)

const (
	appName = "resto"
)

func usage(wri io.Writer) {
	var (
		desc = []string{
			"A minimalist CLI REST client that calls APIs, waits for conditions, and retries intelligently.\n",
			"Supports waiting and polling until your target resource reaches the desired state, using simple JQ expressions.",
			"Makes scripting and automation of REST API calls simpler and more reliable in CI/CD pipelines and development workflows.",
		}

		donateInfo = []string{
			"If you find this tool helpful consider supporting with a donation.",
			"Every bit helps cover development time and fuels future improvements.\n",
			"Your support truly makes a difference — thank you!\n",
			"  * https://www.paypal.com/donate/?hosted_button_id=FV575PVWGXZBY\n",
		}
	)

	fmt.Fprintln(wri)
	fmt.Fprintln(wri, "┬─┐┌─┐┌─┐┌┬┐┌─┐")
	fmt.Fprintln(wri, "├┬┘├┤ └─┐ │ │ │")
	fmt.Fprintln(wri, "┴└─└─┘└─┘ ┴ └─┘")
	fmt.Fprintln(wri)

	for _, el := range desc {
		if el[0] == 194 {
			fmt.Fprintf(wri, "%s\n\n", xtext.Indent(xtext.Wrap(el, 60), "  "))
			continue
		}
		fmt.Fprintf(wri, "%s\n\n", xtext.Wrap(el, 76))
	}
	fmt.Fprintln(wri)

	fmt.Fprint(wri, "USAGE:\n\n")
	fmt.Fprintf(wri, "  %s [FLAGS] URL\n\n", appName)

	fmt.Fprint(wri, "FLAGS:\n\n")
	fmt.Fprint(wri, "  -X, --request          Specify request method to use (default: GET).\n\n")
	fmt.Fprint(wri, "  -H, --header           Add a custom request header (can be specified multiple times).\n")
	fmt.Fprint(wri, "                         Format: 'Key: Value'.\n\n")
	fmt.Fprint(wri, "      --proxy-url        HTTP proxy URL to use for the request.\n\n")
	fmt.Fprint(wri, "  -u, --until            JQ expression to evaluate on JSON response.\n")
	fmt.Fprint(wri, "                         Retries until it evaluates to true.\n\n")
	fmt.Fprint(wri, "      --max-attempts     The maximum number of retry attempts. The operation will be\n")
	fmt.Fprint(wri, "                         retried up to this many times before giving up.\n\n")
	fmt.Fprint(wri, "      --initial-delay    The starting delay duration before the first retry attempt.\n")
	fmt.Fprint(wri, "                         Specify as a time duration (e.g., 100ms, 1s).\n")
	fmt.Fprint(wri, "                         Determines how long to wait before retrying initially.\n\n")
	fmt.Fprint(wri, "      --max-delay        The maximum delay duration allowed between retry attempts.\n")
	fmt.Fprint(wri, "                         Subsequent retries will not exceed this delay.\n\n")
	fmt.Fprint(wri, "      --max-jitter       The maximum random jitter added to the retry delay.\n")
	fmt.Fprint(wri, "                         Specified as a time duration to spread out retry timing.\n\n")
	fmt.Fprint(wri, "      --ca-cert          Base64-encoded CA certificate for verifying the server's TLS cert.\n\n")
	fmt.Fprint(wri, "      --cert             Base64-encoded client certificate (PEM format) for TLS authentication.\n\n")
	fmt.Fprint(wri, "      --cert-key         Base64-encoded private key (PEM format) for the client certificate.\n\n")
	fmt.Fprint(wri, "      --insecure         Skip TLS certificate verification (insecure, use with caution).\n\n")

	fmt.Fprint(wri, "      --username         Username for Basic Auth. Used with --password.\n\n")
	fmt.Fprint(wri, "      --password         Password for Basic Auth. Used with --username.\n\n")

	fmt.Fprint(wri, "      --token            Bearer token for Authorization header.\n\n")

	fmt.Fprint(wri, "  -v, --verbose          Enable verbose output (prints headers, body, debug info).\n\n")
	fmt.Fprint(wri, "      --version          Show version and exit.\n")
	fmt.Fprint(wri, "      --help             Show help and exit.\n")
	fmt.Fprint(wri, "\n\n")

	fmt.Fprint(wri, "ENVIRONMENT:\n\n")
	fmt.Fprint(wri, "  Many long-form flags can alternatively be set using environment variables.\n\n")
	fmt.Fprint(wri, "  You can define them in a `.env` file or export them in your shell.\n\n")
	fmt.Fprint(wri, "  +---------------------+-----------------------+\n")
	fmt.Fprint(wri, "  |  flag               |  environment variable |\n")
	fmt.Fprint(wri, "  |---------------------+-----------------------|\n")
	fmt.Fprint(wri, "  |                     |  SERVER_URL           |\n")
	fmt.Fprint(wri, "  |     --proxy-url     |  PROXY_URL            |\n")
	fmt.Fprint(wri, "  |     --max-attempts  |  MAX_ATTEMPTS         |\n")
	fmt.Fprint(wri, "  |     --initial-delay |  INITIAL_DELAY        |\n")
	fmt.Fprint(wri, "  |     --max-delay     |  MAX_DELAY            |\n")
	fmt.Fprint(wri, "  |     --max-jitter    |  MAX_JITTER           |\n")
	fmt.Fprint(wri, "  |     --ca-cert       |  CA_CERT              |\n")
	fmt.Fprint(wri, "  |     --cert          |  CERT                 |\n")
	fmt.Fprint(wri, "  |     --cert-key      |  CERT_KEY             |\n")
	fmt.Fprint(wri, "  |     --insecure      |  INSECURE             |\n")
	fmt.Fprint(wri, "  |     --token         |  TOKEN                |\n")
	fmt.Fprint(wri, "  |     --username      |  USERNAME             |\n")
	fmt.Fprint(wri, "  |     --password      |  PASSWORD             |\n")
	fmt.Fprint(wri, "  | -v, --verbose       |  VERBOSE              |\n")
	fmt.Fprint(wri, "  +---------------------+-----------------------+\n\n")

	fmt.Fprint(wri, "  Example `.env` file:\n")
	fmt.Fprint(wri, "    TOKEN=your-token-here\n")
	fmt.Fprint(wri, "    SERVER_URL=https://k8s.example.com\n\n")
	fmt.Fprint(wri, "  Shell usage:\n")
	fmt.Fprint(wri, "    export TOKEN=your-token-here\n")
	fmt.Fprintf(wri, "    %s \"$SERVER_URL/api/v1/pods\"\n\n", appName)

	fmt.Fprint(wri, "  » Security Tip: Avoid committing `.env` files containing\n")
	fmt.Fprint(wri, "                  sensitive data to version control.\n")
	fmt.Fprint(wri, "\n\n")

	fmt.Fprint(wri, "EXAMPLES:\n\n")

	fmt.Fprint(wri, " » Perform a simple GET request:\n\n")
	fmt.Fprintf(wri, "     %s https://httpbin.org/get\n\n", appName)

	fmt.Fprint(wri, " » Add custom headers to the request:\n\n")
	fmt.Fprintf(wri, "     %s -H \"X-Token: abc123\" -H \"Accept: application/json\" https://httpbin.org/headers\n\n", appName)

	fmt.Fprint(wri, " » POST a JSON body from stdin:\n\n")
	fmt.Fprintf(wri, "     echo '{\"hello\": \"world\"}' | %s -X POST -H \"Content-Type: application/json\" https://httpbin.org/post\n\n", appName)

	fmt.Fprint(wri, " » Retry until a JQ expression is true:\n\n")
	fmt.Fprintf(wri, "     %s --until '.status == \"ok\"' https://example.com/api/status\n\n", appName)

	fmt.Fprint(wri, " » Use Basic Auth credentials:\n\n")
	fmt.Fprintf(wri, "     %s --username user --password pass https://httpbin.org/basic-auth/user/pass\n\n", appName)

	fmt.Fprint(wri, " » Send request via HTTP proxy:\n\n")
	fmt.Fprintf(wri, "     %s --proxy-url http://localhost:8080 https://httpbin.org/ip\n\n", appName)

	fmt.Fprint(wri, "SUPPORT:\n\n")
	fmt.Fprint(wri, xtext.Indent(strings.Join(donateInfo, "\n"), "  "))
	fmt.Fprint(wri, "\n\n")

	fmt.Fprintln(wri, "Copyright (c) 2025 Luca Sepe")
}

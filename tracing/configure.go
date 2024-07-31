package tracing

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"go.opentelemetry.io/otel"
	otlpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otlphttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type TraceConfiguration struct {
	Endpoint string
	Headers  map[string]string
	Exporter string
}

func Configure(ctx context.Context, cfg *TraceConfiguration, appName string, version string) (func(ctx context.Context) error, error) {
	exporter, err := createExporter(ctx, cfg)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
			semconv.ServiceVersionKey.String(version),
		)),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signals
		fmt.Printf("Received %s, stopping\n", s)

		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}()

	return tp.Shutdown, nil
}

func parseOtelEnvironmentHeaders(fromEnv string) (map[string]string, error) {

	headers := map[string]string{}

	if fromEnv == "" {
		return headers, nil
	}

	for _, pair := range strings.Split(fromEnv, ",") {
		index := strings.Index(pair, "=")
		if index == -1 {
			return nil, fmt.Errorf("unable to parse '%s' as a key=value pair, missing a '='", pair)
		}

		key := strings.TrimSpace(pair[0:index])
		val := strings.TrimSpace(pair[index+1:])

		headers[key] = val
	}

	return headers, nil
}

func createExporter(ctx context.Context, cfg *TraceConfiguration) (sdktrace.SpanExporter, error) {

	switch cfg.Exporter {
	case "none":
		return nil, nil

	case "stdout":
		return stdouttrace.New(stdouttrace.WithPrettyPrint())

	case "stderr":
		return stdouttrace.New(stdouttrace.WithPrettyPrint(), stdouttrace.WithWriter(os.Stderr))
	}

	endpoint := "localhost:4317"
	if cfg.Endpoint != "" {
		endpoint = cfg.Endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(endpoint, "https://") || strings.HasPrefix(endpoint, "http://") {

		opts := []otlphttp.Option{}

		hostAndPort := u.Host
		if u.Port() == "" {
			if u.Scheme == "https" {
				hostAndPort += ":443"
			} else {
				hostAndPort += ":80"
			}
		}
		opts = append(opts, otlphttp.WithEndpoint(hostAndPort))

		if u.Path == "" {
			u.Path = "/v1/traces"
		}
		opts = append(opts, otlphttp.WithURLPath(u.Path))

		if u.Scheme == "http" {
			opts = append(opts, otlphttp.WithInsecure())
		}

		opts = append(opts, otlphttp.WithHeaders(cfg.Headers))

		return otlphttp.New(ctx, opts...)
	} else {
		opts := []otlpgrpc.Option{}

		opts = append(opts, otlpgrpc.WithEndpoint(endpoint))

		isLocal, err := isLoopbackAddress(endpoint)
		if err != nil {
			return nil, err
		}

		if isLocal {
			opts = append(opts, otlpgrpc.WithInsecure())
		}

		opts = append(opts, otlpgrpc.WithHeaders(cfg.Headers))

		return otlpgrpc.New(ctx, opts...)
	}

}

func isLoopbackAddress(endpoint string) (bool, error) {
	hpRe := regexp.MustCompile(`^[\w.-]+:\d+$`)
	uriRe := regexp.MustCompile(`^(http|https)`)

	endpoint = strings.TrimSpace(endpoint)

	var hostname string
	if hpRe.MatchString(endpoint) {
		parts := strings.SplitN(endpoint, ":", 2)
		hostname = parts[0]
	} else if uriRe.MatchString(endpoint) {
		u, err := url.Parse(endpoint)
		if err != nil {
			return false, err
		}
		hostname = u.Hostname()
	}

	ips, err := net.LookupIP(hostname)
	if err != nil {
		return false, err
	}

	allAreLoopback := true
	for _, ip := range ips {
		if !ip.IsLoopback() {
			allAreLoopback = false
		}
	}

	return allAreLoopback, nil
}

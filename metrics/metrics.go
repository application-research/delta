package metrics

import (
	"io"
	"net/http"
	"runtime"

	"github.com/labstack/gommon/log"

	//#nosec G108 - exposing the profiling endpoint is expected
	_ "net/http/pprof"

	"contrib.go.opencensus.io/exporter/prometheus"
	promclient "github.com/prometheus/client_golang/prometheus"
)

// It creates a Prometheus exporter that exports the metrics from the default Prometheus registry
func Exporter() http.Handler {
	// Prometheus globals are exposed as interfaces, but the prometheus
	// OpenCensus exporter expects a concrete *Registry. The concrete type of
	// the globals are actually *Registry, so we downcast them, staying
	// defensive in case things change under the hood.
	registry, ok := promclient.DefaultRegisterer.(*promclient.Registry)
	if !ok {
		log.Warnf("failed to export default prometheus registry; some metrics will be unavailable; unexpected type: %T", promclient.DefaultRegisterer)
	}
	exporter, err := prometheus.NewExporter(prometheus.Options{
		Registry:  registry,
		Namespace: "delta",
	})
	if err != nil {
		log.Errorf("could not create the prometheus stats exporter: %v", err)
	}

	return exporter
}

// It writes the stack traces of all goroutines to the given writer
func WriteAllGoroutineStacks(w io.Writer) error {
	buf := make([]byte, 64<<20)
	for i := 0; ; i++ {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		if len(buf) >= 1<<30 {
			// Filled 1 GB - stop there.
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	_, err := w.Write(buf)
	return err
}

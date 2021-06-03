package promflush_test

import (
	"github.com/StevenLeRoux/promflush"
	"github.com/prometheus/client_golang/prometheus"

	//"fmt"
	"io/ioutil"

	"os"
	"testing"
)

func TestPromflush(t *testing.T) {

	expectedOut := `1622755603191947// test_gauge{foo=bar} 42
`
	ts := "1622755603191947"

	registry := prometheus.NewRegistry()

	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "test_gauge",
			Help: "test gauge",
		},
		[]string{"foo"},
	)

	registry.MustRegister(gauge)

	gauge.With(prometheus.Labels{"foo": "bar"}).Set(42)

	tmpfile, err := ioutil.TempFile("", "prom_registry_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if err := promflush.WriteToTextfile(ts, tmpfile.Name(), registry); err != nil {
		t.Fatal(err)
	}

	fileBytes, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	fileContents := string(fileBytes)

	//t.Errorf(fileContents)
	//fmt.Print(fileContents)

	if fileContents != expectedOut {
		t.Errorf(
			"files don't match, got:\n%s\nwant:\n%s",
			fileContents, expectedOut,
		)
	}

}

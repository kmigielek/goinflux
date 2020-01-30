package goinflux

import (
	"testing"
	"time"

	"github.com/influxdata/influxdb1-client/models"
	client "github.com/influxdata/influxdb1-client/v2"
	"gotest.tools/assert"
)

var currentTime = time.Now()

func TestAs(t *testing.T) {
	var slice []TestStructTime
	um := models.Row{
		Name:    "",
		Columns: []string{"time", "bool_field", "int_field", "float_field"},
		Tags:    map[string]string{"tag": "tag"},
		Values:  [][]interface{}{{currentTime.Unix(), false, 1, 1.5}},
	}
	cases := []struct {
		name     string
		rows     []models.Row
		as       []TestStructTime
		expected []TestStructTime
	}{
		{"UnixTime", []models.Row{um}, slice, unixStructs},
		{"EmptySlice", []models.Row(nil), slice, []TestStructTime(nil)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rr := RawResponse{Response: &client.Response{Results: []client.Result{{Series: c.rows}}}}
			rr.As(&c.as)
			assert.DeepEqual(t, c.expected, c.as)
		})
	}
}

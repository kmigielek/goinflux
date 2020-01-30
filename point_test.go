package goinflux

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	client "github.com/influxdata/influxdb1-client/v2"
	"gotest.tools/assert"
)

type TestStructTime struct {
	Time       interface{} `influx:"time"`
	Tag        string      `influx:"tag"`
	FloatField float64     `influx:"float_field"`
	IntField   int         `influx:"int_field"`
	BoolField  bool        `influx:"bool_field"`
}

type TestStruct struct {
	Tag        string  `influx:"tag"`
	FloatField float64 `influx:"float_field"`
	IntField   int     `influx:"int_field"`
	BoolField  bool    `influx:"bool_field"`
}

var rfcStructs = []TestStructTime{{currentTime, "tag", 1.5, 1, false}}
var unixStructs = []TestStructTime{{currentTime.Unix(), "tag", 1.5, 1, false}}
var structs = []TestStruct{{"tag", 1.5, 1, false}}

func TestSliceToInfluxPoitns(t *testing.T) {
	tp, _ := client.NewPoint("", map[string]string{"tag": "tag"}, map[string]interface{}{"float_field": 1.5, "int_field": 1, "bool_field": false}, currentTime)
	p, _ := client.NewPoint("", map[string]string{"tag": "tag"}, map[string]interface{}{"float_field": 1.5, "int_field": 1, "bool_field": false})
	cases := []struct {
		name     string
		slice    interface{}
		expected []*client.Point
	}{
		{"UnixTime", unixStructs, []*client.Point{tp}},
		{"RFC3339Time", rfcStructs, []*client.Point{tp}},
		{"NoTime", structs, []*client.Point{p}},
		{"EmptySlice", []TestStructTime{}, []*client.Point(nil)},
		{"NilSlice", nil, []*client.Point(nil)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			points, _ := sliceToInfluxPoints(c.slice, "")
			assert.DeepEqual(t, c.expected, points, cmpopts.IgnoreUnexported(client.Point{}))
		})
	}
}

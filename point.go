package goinflux

import (
	"fmt"
	"reflect"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type point struct {
	measurement string
	tags        map[string]string
	fields      map[string]interface{}
	time        *time.Time
}

func sliceToInfluxPoints(slice interface{}, measurement string) ([]*client.Point, error) {
	var influxPoints []*client.Point
	if slice == nil {
		return influxPoints, nil
	}
	s := reflect.ValueOf(slice)
	if s.Len() < 1 {
		return influxPoints, nil
	}
	for i := 0; i < s.Len(); i++ {
		ip, err := createInfluxPoint(measurement, s.Index(i))
		if err != nil {
			return influxPoints, err
		}
		influxPoints = append(influxPoints, ip)
	}
	return influxPoints, nil
}

func createInfluxPoint(measurement string, s reflect.Value) (*client.Point, error) {
	p := point{measurement: measurement, tags: make(map[string]string), fields: make(map[string]interface{})}
	p.parseFieldsFromSliceElement(s)
	if p.time == nil {
		return client.NewPoint(p.measurement, p.tags, p.fields)
	}
	return client.NewPoint(p.measurement, p.tags, p.fields, *p.time)
}

func (p *point) parseFieldsFromSliceElement(v reflect.Value) {
	for j := 0; j < reflect.TypeOf(v).NumField(); j++ {
		p.parseField(v, j)
	}
}

func (p *point) parseField(value reflect.Value, i int) {
	field := value.Type().Field(i)
	tag := field.Tag.Get("influx")
	if tag != "" {
		if tag == "time" {
			p.getTimeFromValue(value, i)
		} else if field.Type.Kind() == reflect.String {
			p.tags[tag] = value.Field(i).String()
		} else {
			p.fields[tag] = value.Field(i).Interface()
		}
	}
}

func (p *point) getTimeFromValue(value reflect.Value, i int) {
	var t time.Time
	fv := value.Field(i).Interface()
	if unix, ok := fv.(int64); ok {
		t, _ = time.Parse(fmt.Sprintf("%d", time.Now().Unix()), fmt.Sprintf("%d", unix))
	} else {
		t = fv.(time.Time)
	}
	p.time = &t
}

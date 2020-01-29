package goinflux

import (
	"reflect"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/mitchellh/mapstructure"
)

//InfluxHTTPClient is a http client to influxdb
type InfluxHTTPClient struct {
	addr     string
	username string
	password string
	Client   client.Client
}

//NewHTTPClient creates new InfluxHTTPClient
func NewHTTPClient(address, username, password string) InfluxHTTPClient {
	return InfluxHTTPClient{
		addr:     address,
		username: username,
		password: password,
	}
}

//Open sets up influxdb http client
func (ic *InfluxHTTPClient) Open() error {
	var err error
	ic.Client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:               ic.addr,
		Username:           ic.password,
		Password:           ic.username,
		InsecureSkipVerify: true,
	})

	return err
}

//RawResponse contains raw influxdb data returned by http client
type RawResponse struct {
	*client.Response
	Err error
}

//ReadQuery performs read operation on influxdb
func (ic InfluxHTTPClient) ReadQuery(query, db, precission string, parameters map[string]interface{}) *RawResponse {
	defer ic.Client.Close()
	q := client.NewQueryWithParameters(query, db, precission, parameters)
	resp, err := ic.Client.Query(q)
	return &RawResponse{resp, err}
}

//As parses raw influx response to value pointed by v
func (rr RawResponse) As(v interface{}) error {
	if rr.Err != nil {
		return rr.Err
	}
	if rr.Response == nil || len(rr.Results) == 0 {
		return nil
	}
	influxData := make([]map[string]interface{}, 0)

	for _, series := range rr.Results[0].Series {
		for _, v := range series.Values {
			r := make(map[string]interface{})
			for i, c := range series.Columns {
				if len(v) >= i+1 {
					r[c] = v[i]
				}
			}
			for tag, val := range series.Tags {
				r[tag] = val
			}

			influxData = append(influxData, r)
		}
	}

	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           v,
		TagName:          "influx",
		WeaklyTypedInput: false,
		ZeroFields:       false,
		DecodeHook: func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if t.ConvertibleTo(reflect.TypeOf(time.Time{})) && f.Kind() == reflect.String {
				return time.Parse(time.RFC3339, data.(string))
			}

			return data, nil
		},
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(influxData)
}

//WriteQuery inserts points to influxdb
func (ic InfluxHTTPClient) WriteQuery(db, measurement, precision, retentionPolicy string, points interface{}) error {
	influxPoints, err := sliceToInfluxPoints(points, measurement)
	if err != nil {
		return err
	}
	return ic.InsertRaw(db, precision, retentionPolicy, influxPoints)
}

//InsertRaw insert raw influx client points to database
func (ic InfluxHTTPClient) InsertRaw(db, precision, retentionPolicy string, points []*client.Point) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:        db,
		Precision:       precision,
		RetentionPolicy: retentionPolicy,
	})
	if err != nil {
		return err
	}

	bp.AddPoints(points)

	err = ic.Client.Write(bp)
	if err != nil {
		return err
	}

	if err := ic.Client.Close(); err != nil {
		return err
	}
	return nil
}

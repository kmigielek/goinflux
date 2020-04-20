# goinflux

[![Build Status](https://travis-ci.org/kmigielek/goinflux.svg?branch=master)](https://travis-ci.org/kmigielek/goinflux) [![GoDoc](https://godoc.org/github.com/kmigielek/goinflux?status.svg)](https://godoc.org/github.com/kmigielek/goinflux) [![Go Report Card](https://goreportcard.com/badge/github.com/kmigielek/goinflux)](https://goreportcard.com/report/github.com/kmigielek/goinflux)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkmigielek%2Fgoinflux.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkmigielek%2Fgoinflux?ref=badge_shield)

Golang ORM helper for influx database.

## Examples

Open HTTP influx client:
```
client := goinflux.NewHTTPClient("http://localhost:8086", "", "")
err := client.Open()
if err != nil {
	log.Error(err)
}
```

Create compatible struct:
```
type Foo struct {
	Time       int64   `influx:"time"`
	Name       string  `influx:"name"`
	FloatValue float64 `influx:"float_value"`
	IntValue   int64   `influx:"int_value"`
	BoolValue  bool    `influx:"bool_value"`
}
```

Insert data to database:
```
var data []Foo
.
.
.
err := client.WriteQuery("EXAMPLE", "examples", "ms", "", data)
```
Read data from database:
```
var data []Foo
err := client.ReadQuery("SELECT * FROM examples", "EXAMPLE", "ms", nil).As(&data)
```

Execute query with parameters:
```
var data []Foo
params := make(map[string]interface{})
params["name"] = "foo"
err := client.ReadQuery("SELECT * FROM examples WHERE name=$name", "EXAMPLE", "ms", params).As(&data)
```

Helper exposes influxdb-client as well so it is possible to use features from standard influx library. For example:
```
client.Client.Query("SELECT * FROM examples")
```

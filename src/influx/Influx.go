package influx

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"time"
)

func Conn() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "root",
		Password: "root",
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

// QueryDB query
func QueryDB(cli client.Client, database string, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: database,
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

// Insert insert
func Insert(cli client.Client, database string, measurement string, tags map[string]string, fields map[string]interface{}) (err error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: database,
	})
	if err != nil {
		return err
	}

	pt, err := client.NewPoint(measurement, tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		return err
	}
	log.Println("insert success")
	return nil
}

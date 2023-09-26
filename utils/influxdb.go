package utils

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"time"
)

var (
	client influxdb2.Client
)

func init() {
	token := "ghEvZJ0xnUsEQTymkQz9Ybnjl8cTfOLJJDxi74M7AmS8shz7ra-8eO5M-mDwzRec3kymkKm0-nnjWGSYCCekqw=="
	url := "http://192.168.40.20:8086"
	client = influxdb2.NewClient(url, token)
}

func Write(tags map[string]string, fields map[string]interface{}, org, bucket string) error {
	writeAPI := client.WriteAPIBlocking(org, bucket)
	point := write.NewPoint("measurement1", tags, fields, time.Now())

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		fmt.Println("write into influxdb failure, ", err)
		return err
	}

	return nil
}

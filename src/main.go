package main

import (
	"database/sql"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/mattn/go-sqlite3"
	"github.com/savsgio/atreugo/v11"
	"iot-demo-golang/src/influx"
	"iot-demo-golang/src/mqtt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	//连接influx
	influxConn := influx.Conn()
	//连接MQTT
	mqConn := mqtt.Conn("golang-demo")
	mqConn.Subscribe("test", 1, func(client MQTT.Client, message MQTT.Message) {
		log.Println(message)
	})
	//连接Sqlite
	sqlConn, sqlConnErr := sql.Open("sqlite3", "test.db")

	if sqlConnErr != nil {
		log.Fatal(sqlConnErr)
	}
	server := atreugo.New(atreugo.Config{Addr: "0.0.0.0:8080"})

	server.GET("/save", func(ctx *atreugo.RequestCtx) error {
		//随机温度
		temperature := rand.Intn(21) + 16
		data := Data{
			SensorName:  "testSensor",
			Temperature: strconv.Itoa(temperature),
		}

		tagMap := make(map[string]string)
		tagMap["id"] = "1"
		filedMap := make(map[string]interface{})
		filedMap["temperature"] = temperature

		//保存到influx
		err := influx.Insert(influxConn, "test", "temperature", tagMap, filedMap)
		if err != nil {
			return ctx.ErrorResponse(err, http.StatusInternalServerError)
		}

		//发送到MQ
		mqConn.Publish("test", 1, false, temperature)

		_, err = sqlConn.Exec("update temperature_data set temperature=" + strconv.Itoa(temperature) + " where id =1")
		if err != nil {
			return ctx.ErrorResponse(err, http.StatusInternalServerError)
		}
		return ctx.JSONResponse(data, http.StatusOK)

	})

	server.GET("/queryResult", func(ctx *atreugo.RequestCtx) error {
		res, err := influx.QueryDB(influxConn, "test", "SELECT MEAN(temperature) FROM temperature WHERE time > now() - 20m")
		if err != nil {
			return ctx.ErrorResponse(err, http.StatusInternalServerError)
		}
		return ctx.JSONResponse(res, http.StatusOK)
	})

	server.GET("/testBlock", func(ctx *atreugo.RequestCtx) error {
		log.Println("testBlock in")
		time.Sleep(time.Second * 3)
		log.Println("testBlock out")
		return ctx.TextResponse("ok", http.StatusOK)
	})

	//启动路由
	routeErr := server.ListenAndServe()
	if routeErr != nil {
		log.Fatal(routeErr)
	}

}

type Data struct {
	SensorName  string
	Temperature string
}

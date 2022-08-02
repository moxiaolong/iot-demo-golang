package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"iot-demo-golang/src/influx"
	"iot-demo-golang/src/mqtt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	//连接influx
	influxConn := influx.Conn()
	//连接MQTT
	mqConn := mqtt.Conn("golang-demo")
	//连接Sqlite
	sqlConn, sqlConnErr := sql.Open("sqlite3", "test.db")
	if sqlConnErr != nil {
		log.Fatal(sqlConnErr)
	}

	//创建路由
	router := gin.Default()
	// 绑定路由规则，执行的函数
	router.GET("/save", func(ctx *gin.Context) {
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
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		//发送到MQ
		mqConn.Publish("test", 1, false, temperature)

		//保存到Sqlite
		_, err = sqlConn.Exec("update temperature_data set temperature=" + strconv.Itoa(temperature) + " where id =1")
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, data)
	})

	router.GET("/queryResult", func(ctx *gin.Context) {
		res, err := influx.QueryDB(influxConn, "test", "SELECT MEAN(temperature) FROM temperature WHERE time > now() - 20m")
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, res)
	})

	router.GET("/testBlock", func(ctx *gin.Context) {
		log.Println("testBlock in")
		//协程阻塞
		go func() {
			time.Sleep(time.Second * 3)
			ctx.String(http.StatusOK, "ok")
			log.Println("testBlock ok")
		}()
		log.Println("testBlock out")
	})

	//启动路由
	routeErr := router.Run(":8080")
	if routeErr != nil {
		log.Fatal(routeErr)
	}

}

type Data struct {
	SensorName  string
	Temperature string
}

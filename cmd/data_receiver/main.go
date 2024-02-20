package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wvalencia19/tolling/types"
)

var kafkaTopic = "obudata"

func main() {

	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgCh chan types.OBUData
	conn  *websocket.Conn
	prod  DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p   DataProducer
		err error
	)

	p, err = NewKafkaProducer()
	if err != nil {
		return nil, err
	}

	p = NewLogMiddleWare(p)

	return &DataReceiver{
		msgCh: make(chan types.OBUData, 10),
		prod:  p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}

	conn, err := u.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsServe()
}

func (dr *DataReceiver) wsServe() {
	fmt.Println("new OBU client connected")
	for {
		var data types.OBUData

		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
		}
		if err := dr.produceData(data); err != nil {
			fmt.Println("kafka produce error", err)
		}
	}
}

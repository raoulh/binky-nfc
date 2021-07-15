package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/raoulh/acr122u"
	"github.com/raoulh/binky-nfc/mdns"
	"github.com/raoulh/binky-nfc/gpio"
)

var (
	quitCh  chan interface{}
	conn    *websocket.Conn
	macAddr string
)

//Init main streamr app
func Init(mac string) error {
	if mac == "" {
		return fmt.Errorf("no mac address defined")
	}

	macAddr = mac
	quitCh = make(chan interface{})

	return nil
}

//Run all background tasks
func Run() {
	go func() {
		//discover websocket server and connect

		for {
			err := runWs()
			if err != nil {
				select {
				case <-quitCh:
					return
				default:
					log.Printf("failure: %v", err)
				}
			}

			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		//setup NFC and listen

		for {
			acrCtx, err := acr122u.EstablishContext()
			if err != nil {
				time.Sleep(time.Second * 1)
				break
			}

			h := &handler{log.New(os.Stdout, "", 0)}

			log.Println("Waiting for NFC card...")
			acrCtx.Serve(h)
		}
	}()

	go func() {
		//Switch LED on GPIO11 high
		gpio.Pin("11").Output().High()
	}
}

//Shutdown all background jobs
func Shutdown(ctx context.Context) error {
	var e error

	close(quitCh)

	return e
}

func runWs() error {
	// Discover all services on the network
	host, err := mdns.DiscoverBinkyServer()
	if err != nil {
		return err
	}

	log.Println("Found host:", host)

	conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", host), nil)
	if err != nil {
		return err
	}
	defer conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	defer conn.Close()

	log.Println("Websocket connected")

	readWebsocket(conn)

	return fmt.Errorf("end")
}

func readWebsocket(c *websocket.Conn) {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			_ = c.Close()
			return
		}

		log.Printf("recv: %s", msg)
	}
}

type handler struct {
	acr122u.Logger
}

type MsgData struct {
	ClientId string           `json:"client_id,omitempty"`
	Msg      string           `json:"msg,omitempty"`
	Payload  *json.RawMessage `json:"payload,omitempty"`
}

type NFCMsgData struct {
	ClientId string     `json:"client_id,omitempty"`
	Msg      string     `json:"msg,omitempty"`
	Payload  NFCMessage `json:"payload,omitempty"`
}

type NFCMessage struct {
	Msg        string `json:"msg,omitempty"`
	NFCID      string `json:"nfc_id,omitempty"`
	MACAddress string `json:"mac_address,omitempty"`
}

func (h *handler) ServeCard(c acr122u.Card) {
	m := &NFCMsgData{
		Msg: "nfc",
		Payload: NFCMessage{
			Msg:        "present",
			NFCID:      fmt.Sprintf("%x", c.UID()),
			MACAddress: macAddr},
	}

	data, err := json.Marshal(m)
	if err != nil {
		return
	}

	log.Printf("Card present %s", m.Payload.NFCID)

	if conn != nil {
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("write:", err)
			_ = conn.Close()
			return
		}
	}
}

func (h *handler) CardRemoved() {
	m := &NFCMsgData{
		Msg: "nfc",
		Payload: NFCMessage{
			Msg:        "removed",
			MACAddress: macAddr},
	}

	data, err := json.Marshal(m)
	if err != nil {
		return
	}

	log.Printf("Card removed")

	if conn != nil {
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("write:", err)
			_ = conn.Close()
			return
		}
	}
}

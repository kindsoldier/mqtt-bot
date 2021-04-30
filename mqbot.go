/*
 * Copyright:  Pixel Networks <support@pixel-networks.com> 
 * Author: Oleg Borodin <oleg.borodin@pixel-networks.com>
 */

package main

import (
    "flag"
    "fmt"
    "os"
    "encoding/json"
    "path/filepath"
    "log"
    "strconv"
    "time"
    "strings"
    "sync"

    "app/pmconfig"
    "app/pmtools"
    "app/pmdaemon"
    "app/transport"
)

type Application struct {
    config  *pmconfig.Config
    trans   *transport.Transport
    room    *Room
}

type Room struct {
    Temp            int         `json:"temp"`
    TempMtx         sync.Mutex  `json:"-"`
    Light           int         `json:"light"`
    LightMtx        sync.Mutex  `json:"-"`
    Current         int         `json:"current"`
    CurrentMtx      sync.Mutex  `json:"-"`

    LightOn         bool        `json:"lightOn"`
    LightOnMtx      sync.Mutex  `json:"-"`
    CurrentOn       bool        `json:"currentOn"`
    CurrentOnMtx    sync.Mutex  `json:"-"`
}

func NewRoom() *Room {
    return &Room{
        LightOn:        true,
        CurrentOn:      true,
    }
}

func (this *Room) Json() []byte {
    this.LightOnMtx.Lock() 
    this.CurrentOnMtx.Lock()
    res, _ := json.Marshal(this)
    this.LightOnMtx.Unlock()
    this.CurrentOnMtx.Unlock()
    return res
}

func (this *Room) GetLightOn() bool {
    this.LightOnMtx.Lock()
    defer this.LightOnMtx.Unlock()
    return this.LightOn
}

func (this *Room) SetLightOn(value bool) {
    this.LightOnMtx.Lock()
    this.LightOn = value
    this.LightOnMtx.Unlock()
}

func (this *Room) GetCurrentOn() bool {
    this.CurrentOnMtx.Lock()
    defer this.CurrentOnMtx.Unlock()
    return this.CurrentOn
}
func (this *Room) SetCurrentOn(value bool) {
    this.CurrentOnMtx.Lock()
    this.CurrentOn = value
    this.CurrentOnMtx.Unlock()
}

func (this *Room) GetLight() int {
    this.LightMtx.Lock()
    defer this.LightMtx.Unlock()
    return this.Light
}

func (this *Room) SetLight(value int) {
    this.LightMtx.Lock()
    this.Light = value
    this.LightMtx.Unlock()
}

func (this *Room) GetCurrent() int {
    this.CurrentMtx.Lock()
    defer this.CurrentMtx.Unlock()
    return this.Current
}
func (this *Room) SetCurrent(value int) {
    this.CurrentMtx.Lock()
    this.Current = value
    this.CurrentMtx.Unlock()
}

func (this *Room) GetTemp() int {
    this.TempMtx.Lock()
    defer this.TempMtx.Unlock()
    return this.Temp
}
func (this *Room) SetTemp(value int) {
    this.TempMtx.Lock()
    this.Temp = value
    this.TempMtx.Unlock()
}

func NewApplication() *Application {
    config := pmconfig.New()
    room := NewRoom()
    return &Application{
        config:     config,
        room:       room,
    }
}

func (this *Application) Configure() {
    /* Parse cli options */
    flag.BoolVar(&this.config.Debug, "debug", this.config.Debug, "debug mode")
    flag.BoolVar(&this.config.Devel, "devel", this.config.Debug, "devel mode")
    flag.BoolVar(&this.config.Foreground, "foreground", this.config.Foreground, "foreground mode")

    flag.StringVar(&this.config.Broker.Hostname, "host", this.config.Broker.Hostname, "broker hostname")
    flag.IntVar(&this.config.Broker.Port, "port", this.config.Broker.Port, "broker port")
    flag.StringVar(&this.config.Broker.Username, "user", this.config.Broker.Username, "broker username")
    flag.StringVar(&this.config.Broker.Password, "pass", this.config.Broker.Password, "broker password")

	flag.BoolVar(&this.config.Operation.ShowConfig, "config", this.config.Operation.ShowConfig, "custom operation: show current config")
	flag.BoolVar(&this.config.Operation.ShowVersion, "version", this.config.Operation.ShowVersion, "custom operation: show version")

    exeName := filepath.Base(os.Args[0])

    flag.Usage = func() {
        fmt.Println(exeName + " version " + this.config.Version)
        fmt.Println("")
        fmt.Printf("usage: %s command [option]\n", exeName)
        fmt.Println("")
        flag.PrintDefaults()
        fmt.Println("")
    }
    flag.Parse()

    //if len(os.Getenv("POG_DEBUG")) > 0 {
    //    this.config.Debug = true
    //}
}

func (this *Application) Start() error {
    var err error
    this.Configure()

	switch {
		case this.config.Operation.ShowVersion == true:
			this.ShowVersion()
			return err
		case this.config.Operation.ShowConfig == true:
			this.ShowConfig()
			return err
	}

    this.config.Foreground = true
    daemon := pmdaemon.NewDaemon(
        this.config.MessageLogPath,
        this.config.PidPath,
        this.config.Debug,
        this.config.Foreground)

    err = daemon.Daemonize()
    if err != nil {
        return err
    }
    pmdaemon.SetSignalHandler()
    
	// Default operation
    err = this.Run()
    if err != nil {
        return err
    }
    return err
}

func (this *Application) Run() error {
    var err error

    this.trans = transport.New()
    err = this.trans.Bind(this.config.Broker.Hostname,
            this.config.Broker.Port,
            this.config.Broker.Username,
            this.config.Broker.Password)

    lightOnFunc := func(topic string, payload []byte) {
        log.Println("topic", topic, "payload", string(payload))
        switch strings.ToLower(string(payload)) {
            case "on":
                this.room.SetLightOn(true)
            case "off":
                this.room.SetLightOn(false)
        }
    } 

    currentOnFunc := func(topic string, payload []byte) {
        log.Println("topic", topic, "payload", string(payload))
        switch strings.ToLower(string(payload)) {
            case "push":
                if this.room.GetCurrentOn() {
                    this.room.SetCurrentOn(false)
                } else {
                    this.room.SetCurrentOn(true)
                }
        }
    } 

    this.trans.Subscribe("/room1/currentOn", currentOnFunc)
    this.trans.Subscribe("/room1/lightOn", lightOnFunc)
    
	timer := time.NewTicker(2000 * time.Millisecond)

	for timeX := range timer.C {
        
            if this.room.GetCurrentOn() { 
                this.room.SetCurrent(pmtools.GetRandomInt(10, 15))
            } else {
                this.room.SetCurrent(0)
            }

            if this.room.GetLightOn() {
                this.room.SetLight(pmtools.GetRandomInt(1000, 1500))
            } else {
                this.room.SetLight(0)
            }
            this.room.SetTemp(pmtools.GetRandomInt(150, 250))

            this.trans.Publish("/room1/current", strconv.Itoa(this.room.GetCurrent()))
            this.trans.Publish("/room1/light", strconv.Itoa(this.room.GetLight()))
            this.trans.Publish("/room1/temp", strconv.Itoa(this.room.GetTemp()))

            this.trans.Publish("/room1/state", string(this.room.Json()))
            this.trans.Publish("/time", timeX.String())
	}
    return err
}

func (this *Application) ShowVersion() {
	fmt.Println("version: ", this.config.Version)
	return
}

func (this *Application) ShowConfig() {
	fmt.Println(this.config.ToYaml())
	return
}

func main() {
    app := NewApplication()
    err := app.Start()
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
}

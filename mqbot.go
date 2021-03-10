/*
 * Copyright:  Pixel Networks <support@pixel-networks.com> 
 * Author: Oleg Borodin <oleg.borodin@pixel-networks.com>
 */

package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    //"time"

    //"encoding/json"
    //"errors"
    //"log"
    //"strconv"
    //"strings"
    "time"
    //"net/url"

    "app/config"
    //"app/tools"
    "app/transport"
    "app/daemon"
    
    //mqtt "github.com/eclipse/paho.mqtt.golang"
    //"log"
)

type Application struct {
    config  *config.Config
}

func NewApplication() *Application {
    config := config.New()
    return &Application{
        config:     config,
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

func (this *Application) Run() error {
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

    daemon := daemon.New(this.config)
    //err = daemon.Daemonize()
    //if err != nil {
    //    return err
    //}
    //}
    daemon.SetSignalHandler()
    
	// Default operation
    err = this.Loop()
    if err != nil {
        return err
    }
    return err
}

func (this *Application) Loop() error {
    var err error

    trans := transport.New()
    err = trans.Bind(this.config.Broker.Hostname,
            this.config.Broker.Port,
            this.config.Broker.Username,
            this.config.Broker.Password)

	timer := time.NewTicker(1 * time.Second)
	for time := range timer.C {
      trans.Publish("time", time.String())
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
    err := app.Run()
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
}

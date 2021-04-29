/*
 * Copyright: Pixel Networks <support@pixel-networks.com> 
 * Author: Oleg Borodin <oleg.borodin@pixel-networks.com>
 */

package pmconfig

import (
    //"path/filepath"

    "io/ioutil"
    "encoding/json"
    "os"

    "github.com/go-yaml/yaml"
)

type Config struct {
    Verbose             bool    `yaml:"-"`
    Version             string  `yaml:"-"`

    Foreground          bool    `yaml:"-"`

    Debug               bool    `yaml:"debug"`
    Devel               bool    `yaml:"-"`

    ConfigPath          string  `yaml:"-"`
    LibDir              string  `yaml:"-"`
    DataDir             string  `yaml:"datadir"`
    PidPath             string  `yaml:"pidfile"`
    MessageLogPath      string  `yaml:"messagelog"`
    AccessLogPath       string  `yaml:"accesslog"`

    User                string  `yaml:"user"`
    Group               string  `yaml:"group"`
    //CertPath            string  `yaml:"cert"`
    //KeyPath             string  `yaml:"key"`

    Operation       Operation   `yaml:"-"`
    Broker          Broker      `yaml:"broker"`
}

type Operation struct {
    ShowVersion     bool
    ShowConfig      bool
}

type Broker struct {
    Hostname        string  `yaml:"hostname"`
    Port            int     `yaml:"port"`
    Username        string  `yaml:"username"`
    Password        string  `yaml:"password"`
}

func (this *Config) ToJson() string {
    json, _ := json.MarshalIndent(this, "", "    ")
    return string(json)
}
func (this *Config) ToYaml() string {
    json, _ := yaml.Marshal(this)
    return string(json)
}

func (this *Config) Write(fileName string) error {
    var data []byte
    var err error

    //fileName, _ := filepath.Abs(this.ConfigPath)
    os.Rename(fileName, fileName + "~")

    data, err = yaml.Marshal(this)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(fileName, data, 0640)
}

func (this *Config) Read(fileName string) error {
    var data []byte
    var err error

    //fileName, _ := filepath.Abs(this.ConfigPath)
    data, err = ioutil.ReadFile(fileName)
    if  err != nil {
        return err
    }
    return yaml.Unmarshal(data, &this)
}

func New() *Config {
    broker := Broker{
        Hostname:       "v7.unix7.org",
        Port:           1883,
        Username:       "device",
        Password:       "qwerty",
    }
    return &Config{
        Debug:              false,
        Verbose:            false,

        Devel:          false,
        Foreground:     false,

        ConfigPath:     "/home/ziggi/projects/pix-mqtt-bot//mqbot.yml",
        LibDir:         "/home/ziggi/projects/pix-mqtt-bot/",
        DataDir:        "/home/ziggi/projects/pix-mqtt-bot",

        PidPath:        "/home/ziggi/projects/pix-mqtt-bot/run/mqbot.pid",
        MessageLogPath: "/home/ziggi/projects/pix-mqtt-bot/log/message.log",
        AccessLogPath:  "/home/ziggi/projects/pix-mqtt-bot/log/access.log",

        User:           "root",
        Group:          "root",

        //CertPath:       "/home/ziggi/projects/pix-mqtt-bot//mqbot.crt",
        //KeyPath:        "/home/ziggi/projects/pix-mqtt-bot//mqbot.key",

        Version:            "0.1.1",
        Broker:             broker,
    }
}


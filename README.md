# MQTT bot

Simple publisher service for emulating different classes of MQTT devices 
for development and verification purposes

### Individual topics

* /room1/light - integer as string, 10..15
* /room1/current - integer as string, 1000..1500
* /room1/temp - integer as string, 150..250
* /time - string, real time

### JSON topic


* /room1/state 

`{"temp":197,"light":1037,"current":15,"lightOn":true,"currentOn":true}`

### Command topics

$ mosquitto_sub -h v7.unix7.org -u gateway -P qwerty -t /room1/state  

`{"temp":197,"light":1037,"current":15,"lightOn":true,"currentOn":true}`

#### Pseudo button, topic "/room1/currentOn", command: "push" 

$ mosquitto_pub -h v7.unix7.org -u device -P qwerty -t /room1/currentOn -m push

`{"temp":161,"light":1469,"current":0,"lightOn":true,"currentOn":false}`

$ mosquitto_pub -h v7.unix7.org -u device -P qwerty -t /room1/currentOn -m push

#### Pseudo switch, topic "/room1/lightOn", commands: "on", "off"

$ mosquitto_pub -h v7.unix7.org -u device -P qwerty -t /room1/lightOn -m off

`{"temp":180,"light":0,"current":14,"lightOn":false,"currentOn":true}`

$ mosquitto_pub -h v7.unix7.org -u device -P qwerty -t /room1/lightOn -m on

`{"temp":158,"light":1324,"current":13,"lightOn":true,"currentOn":true}`


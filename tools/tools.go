/*
 * Copyright: Pixel Networks <support@pixel-networks.com> 
 * Author: Oleg Borodin <oleg.borodin@pixel-networks.com>
 */


package tools

import (
    "math/rand"
    "time"

    //"github.com/tidwall/pretty"
    "github.com/satori/go.uuid"
)

type UUID = string

func GetNewUUID() string {
    id := uuid.NewV4()
    return id.String()
}

func GetRandomBool() bool {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(2) == 1
}

func GetRandomPercent() int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(100)
}

func GetRandomInt(min int, max int) int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(max - min + 1) + min
}

//EOF


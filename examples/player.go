package main

import (
	"fmt"
	. "github.com/trejjam/spotify/accessToken"
	. "github.com/trejjam/spotify/player"
	"os"
)

func printDevices(accessToken *AccessToken) {
	devices, err := GetDevices(accessToken)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(devices.Devices); i++ {
		device := devices.Devices[i]

		fmt.Println(device.String())
	}
}

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Errorf("Use username and password as arguments"))
		return
	}

	username := os.Args[1]
	password := os.Args[2]

	accessToken, err := GetAccessToken(username, password)
	if err != nil {
		panic(err)
	}

	printDevices(accessToken)
}

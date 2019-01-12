package main

import (
	"fmt"
	. "github.com/trejjam/spotify/accessToken"
	"github.com/trejjam/spotify/object"
	. "github.com/trejjam/spotify/player"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Errorf("Use username and password as arguments"))
		return
	}

	username := os.Args[1]
	password := os.Args[2]
	var deviceId1 string
	if len(os.Args) > 2 {
		deviceId1 = os.Args[3]
	}
	var deviceId2 string
	if len(os.Args) > 3 {
		deviceId2 = os.Args[4]
	}

	accessToken, err := GetAccessToken(username, password)
	if err != nil {
		panic(err)
	}

	devices, err := GetDevices(accessToken)
	if err != nil {
		panic(err)
	}

	fmt.Println(devices.String())

	var device1 *object.Device
	var device2 *object.Device
	for i := 0; i < len(devices.Devices); i++ {
		iDevice := devices.Devices[i]

		if iDevice.Id == deviceId1 {
			device1 = &iDevice
		}
		if iDevice.Id == deviceId2 {
			device2 = &iDevice
		}
	}

	if device1 != nil {
		err = TransferPlayback(accessToken, true, device1)
		if err != nil {
			panic(err)
		}

		err = SetRepeatMode(accessToken, RepeatModeTrack, device1)
		if err != nil {
			panic(err)
		}

		err = SetRepeatMode(accessToken, RepeatModeContext, device1)
		if err != nil {
			panic(err)
		}

		err = SetRepeatMode(accessToken, RepeatModeOff, device1)
		if err != nil {
			panic(err)
		}

		err = NextTrack(accessToken, device1)
		if err != nil {
			panic(err)
		}

		time.Sleep(2 * time.Second)

		err = Pause(accessToken, device1)
		if err != nil {
			panic(err)
		}

		time.Sleep(2 * time.Second)

		err = Play(accessToken, device1)
		if err != nil {
			panic(err)
		}

		time.Sleep(2 * time.Second)

		err = PreviousTrack(accessToken, device1)
		if err != nil {
			panic(err)
		}

		time.Sleep(2 * time.Second)

		if device2 != nil {
			err = Pause(accessToken, device1)
			if err != nil {
				panic(err)
			}

			time.Sleep(1 * time.Second)

			err = TransferPlayback(accessToken, false, device2)
			if err != nil {
				panic(err)
			}

			time.Sleep(1 * time.Second)

			err = Play(accessToken, device1)
			if err != nil {
				panic(err)
			}

			time.Sleep(3 * time.Second)

			err = TransferPlayback(accessToken, true, device1)
			if err != nil {
				panic(err)
			}
		}
	}
}

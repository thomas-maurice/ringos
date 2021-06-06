# RingOS
A firmware for ESP8266 microcontroller to drive LED strips and rings

![An example LED ring in action](https://github.com/thomas-maurice/ringos/blob/master/_assets/ring.gif)

## What is this ?
This project is a firmware that aims at providing an easy way to drive LED strips (WS2812B at the moment).
It offers an API as well as a web interface that should cover the basic modes of operations you want.
For now, it supports a handfuls of modes of operation:

* Static lighting
* Breathing mode
* Chase mode (the ring or light having a bright LED then a trail of fading ones moving along the strip)
  * Speed is adjustable
  * Direction is adjustable

In all the colour modes you can adjust the following:
* Colour
* Brightness

## Setting it up
### Hardware
You need 2 components to get started, an ESP8266 microcontroller and a WS2812B LED ring or strip, of any length (the firmware hardcodes 30, but this can be changed)

You need to connect both to power, then wire the DIN (or DI, or IN depending on how the packaging is marked) pin of the strip to the D3 pin of the ESP, and you are done.

### Tuning the firmwares
You can edit a few behaviours in the [src/main.cpp](https://github.com/thomas-maurice/ringos/blob/master/src/main.cpp) file

### Firmware
Simply flash the firmware to the device.
If you are using the esp12e module:
```
pio run -e esp12e -t upload
pio run -e esp12e -t uploadfs
```

If you are using a Wemos D1 module:
```
pio run -e d1_mini -t upload
pio run -e d1_mini -t uploadfs
```

:warning: You might need to edit the `upload_port` and `upload_protocol` variables from your target's [platformio.ini](https://github.com/thomas-maurice/ringos/blob/master/platformio.ini) file to upload through serial, depending on what the name of the serial port of your module is.

:information_source: After the initial flash is done, and you managed to connect the module to the wifi, subsequent updates can happen using the Arduino Over-The-Air method that allows you to upload the code/UI over WiFi.

### Setting it up
When you first power up the device, it will not know what network to connect to. It will then broadcast its own SSID, it will look like `ringos-<some digits>`, it will be passwordless. Connect to it, navigate the network tab and add your own network, then you can reboot the device and it will connect to your network when it starts.

Note that you can add as many networks as you want (with respect of the internal flash size of course). At startup if any of the defined networks is available, the strongest one should be joined.

## Screenshots
![Main screen](https://github.com/thomas-maurice/ringos/blob/master/_assets/main.png)

## API documentation
TODO

## Limitations
Eventhough the project works pretty well as is, it has a few limitations

* This is an 80Mhz microcontroller, it won't handle 100 requests per second
* Serving the UI can be quite slow, use the API whenever possible

## Further plans

- [ ] API documentation
- [ ] Revamp the API to make it more consistant
- [ ] Make the UI on par with the API
- [ ] Optional password authentication for the API calls
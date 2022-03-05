package ringo

const (
	DIRECTION_CLOCKWISE         = 1
	DIRECTION_COUNTER_CLOCKWISE = -1
	MODE_BREATHING              = "breathing"
	MODE_STATIC                 = "static"
	MODE_CHASE                  = "chase"
)

// ColourResponse is the result of a call to the API
// to get the status of the device colour/mode wise
type ColourResponse struct {
	Colour     string `json:"colour"`
	R          int64  `json:"R"`
	G          int64  `json:"G"`
	B          int64  `json:"B"`
	Brightness int64  `json:"brightness"`
	Mode       string `json:"mode"`
}

func (c *ColourResponse) ToColourRequest() *ColourRequest {
	return &ColourRequest{
		Colour:     c.Colour,
		Brightness: c.Brightness,
		Mode:       c.Mode,
	}
}

type ColourRequest struct {
	Colour     string `json:"colour,omitempty"`
	Mode       string `json:"mode,omitempty"`
	Brightness int64  `json:"brightness,omitempty"`
	Persist    bool   `json:"persist"`
}

// BreathingResponse gets you how the device
// is configured for the breathing mode
type BreathingResponse struct {
	Speed int64 `json:"speed,omitempty"`
}

type BreathingRequest struct {
	BreathingResponse
	Persist bool `json:"persist"`
}

func (b *BreathingResponse) ToBreathingRequest() *BreathingRequest {
	return &BreathingRequest{
		BreathingResponse: *b,
	}
}

// ChaseResponse gets you how the device
// is configured for the chase mode
type ChaseResponse struct {
	Speed int64 `json:"speed,omitempty"`
	// 1 is clockwise, -1 is counter-clockwise
	Direction int64 `json:"direction,omitempty"`
	Length    int64 `json:"length,omitempty"`
}

type ChaseRequest struct {
	ChaseResponse
	Persist bool `json:"persist"`
}

func (c *ChaseResponse) ToChaseRequest() *ChaseRequest {
	return &ChaseRequest{
		ChaseResponse: *c,
		Persist:       false,
	}
}

// Device status gives info about the device itself
type DeviceStatus struct {
	FreeHeap   int64  `json:"free_heap"`
	SSID       string `json:"ssid"`
	MAC        string `json:"mac"`
	Address    string `json:"address"`
	LEDOn      bool   `json:"led_on"`
	Hostname   string `json:"hostname"`
	LEDs       int64  `json:"leds"`
	WiFiStatus string `json:"wifi_status"`
	ChipID     string `json:"chip_id"`
}

// ScannedNetwork is a network as scanned by the device
type ScannedNetwork struct {
	RSSI    int64  `json:"rssi"`
	SSID    string `json:"ssid"`
	BSSID   string `json:"bssid"`
	Channel int64  `json:"channel"`
	Secure  int64  `json:"secure"`
	Hidden  bool   `json:"hidden"`
}

// ScannedNetworks is a list of ScannedNetwork
type ScannedNetworks []ScannedNetwork

type GenericAnswer struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Error   string `json:"error"`
}

type KnownNetworkList []string

type DeleteNetworkRequest string

type AddNetworkRequest struct {
	SSID string `json:"ssid"`
	PSK  string `json:"psk"`
}

type Config struct {
	Hostname *string `json:"hostname"`
}

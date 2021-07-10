#include <Arduino.h>
#include <ESP8266WiFi.h>
#include <ESP8266WiFiMulti.h>
#include <ESP8266mDNS.h>
#include <ESPAsyncTCP.h>
#include <ESPAsyncWebServer.h>
#include <AsyncJson.h>
#include <ArduinoJson.h>
#include <ArduinoOTA.h>
#include <LittleFS.h>
#include <FastLED.h>

#include "fsUtil.h"
#include "mdnsUtil.h"
#include "otaUtil.h"
#include "misc.h"
#include "hash.h"
#include "captivePortal.h"
#include "webUtils.h"

/* Chip specific settings */
// LED pin for blinking purposes
#define LED_PIN LED_BUILTIN

// Colour order for FastLED
#define COLOUR_ORDER GRB

// Chipset for FastLED
#define CHIPSET WS2812B

// Maximum number of LEDs this firmware will
// manage, the actual number of LEDs can be
// adjusted changing the NUM_LEDS variable.
// To save on ram lower this since they
// will be allocated even if not used
#define MAX_NUM_LEDS 120
#define DEFAULT_NUM_LEDS 24

// Pin onto which the ring/strip is plugged
#define SIG_LED D3

/* Some default parameters */
// Default LED brightness
#define DEFAULT_BRIGHTNESS 100

// Default operation mode
// Can be any of:
// * chase
// * static
// * breathing
#define DEFAULT_OPERATION_MODE "chase"

// Number of times per second FastLED is updated
// Higher numers tend to result in poorer performance
#define FRAMES_PER_SECOND 30

// Default chase speed, lower is wuicker
#define DEFAULT_CHASE_SPEED 3

// Default chase direction 1 is clockwise
#define DEFAULT_CHASE_DIRECTION 1

/* Resets the chip (effectively a soft reboot) */
void (*softReset)(void) = 0;

/* Variables to control reboot */
// Has a reboot order been issued ?
bool ASKED_FOR_REBOOT = false;
// In how many seconds the reboot will happen ?
// this is mostly so the api client gets the
// full replies before the connexion resets
long REBOOT_TICKS = 3 * FRAMES_PER_SECOND;

AsyncWebServer server(80);
ESP8266WiFiMulti wifiMulti;
DNSServer dnsServer;

/* Global variables to manage the brightness and colour modes */
int R = 0;
int G = 0;
int B = 0;
int NUM_LEDS = MAX_NUM_LEDS;
int BRIGHTNESS = DEFAULT_BRIGHTNESS;
String OPERATION_MODE = DEFAULT_OPERATION_MODE;

/* Variables for the LED chase mode */
int chasePauseCounter = 0;
int chaseSpeed = 3;
int chaseLeadPosition = 0;
int chaseDirection = -1;
int chaseTrailLength = NUM_LEDS - 1;

/* Representation of the LEDs stip/ring */
CRGB leds[MAX_NUM_LEDS];

// Returns the WiFi status in a more human readable form
String getWiFiStatus()
{
  switch (WiFi.status())
  {
  case WL_CONNECTED:
    return String("connected");
    break;
  case WL_NO_SSID_AVAIL:
    return "unknown SSID";
    break;
  case WL_IDLE_STATUS:
    return "idle";
    break;
  case WL_CONNECT_FAILED:
    return "connection failed";
    break;
  case WL_DISCONNECTED:
    return "disconnected";
    break;
  default:
    return "unknown";
    break;
  }
}

// Reads the colour from the configuration
int readColour(String colour)
{
  String c = readFile("/config/" + colour, 16);
  int intColour;
  sscanf(c.c_str(), "%d", &intColour);
  return intColour;
}

String readMode()
{
  return readFile("/config/mode", 16);
}

void saveBrightness(int value)
{
  char stringBrightness[16];
  sprintf(stringBrightness, "%d", value);
  writeFile("/config/brightness", stringBrightness);
}

int readBrightness()
{
  String c = readFile("/config/brightness", 16);
  int intColour;
  sscanf(c.c_str(), "%d", &intColour);
  if (intColour < 0 || intColour > 255)
  {
    saveBrightness(DEFAULT_BRIGHTNESS);
    Serial.println(PSTR("saving default brightness because of an invalid value"));
    return DEFAULT_BRIGHTNESS;
  }
  return intColour;
}

void saveColour(String colour, int value)
{
  char stringColour[16];
  sprintf(stringColour, "%d", value);
  writeFile("/config/" + colour, stringColour);
}

void setup()
{
  // Start the filesystem
  LittleFS.begin();
  // Start the serial port with a decent baudrate
  Serial.begin(115200);
  delay(10);
  Serial.println("\n\n");
  // Setup the builtin LED
  pinMode(SIG_LED, OUTPUT);
  pinMode(LED_PIN, OUTPUT);

  digitalWrite(LED_PIN, HIGH);

  NUM_LEDS = readIntFile("/config/leds", 16);
  if (NUM_LEDS == 0)
  {
    NUM_LEDS = DEFAULT_NUM_LEDS;
  }
  FastLED.addLeds<CHIPSET, SIG_LED, COLOUR_ORDER>(leds, NUM_LEDS).setCorrection(UncorrectedColor);

  BRIGHTNESS = readBrightness();
  chaseSpeed = readIntFile("/config/chase/speed", 16);
  if (chaseSpeed == 0)
  {
    chaseSpeed = DEFAULT_CHASE_SPEED;
  }
  chaseTrailLength = readIntFile("/config/chase/trail", 16);
  if (chaseTrailLength < 0 || chaseTrailLength >= NUM_LEDS)
  {
    chaseTrailLength = NUM_LEDS - 1;
  }
  chaseDirection = readIntFile("/config/chase/direction", 16);
  if (chaseDirection != -1 && chaseDirection != 1)
  {
    chaseDirection = DEFAULT_CHASE_DIRECTION;
  }
  OPERATION_MODE = readFile("/config/mode", 16);
  if (OPERATION_MODE == "")
  {
    OPERATION_MODE = DEFAULT_OPERATION_MODE;
  }

  FastLED.setBrightness(BRIGHTNESS);

  for (int i = 0; i < NUM_LEDS; i++)
  {
    leds[i] = CRGB(100, 100, 100);
  }
  FastLED.show();

  pinMode(SIG_LED, OUTPUT);
  digitalWrite(SIG_LED, LOW);

  // Setup the WiFi
  WiFi.setAutoConnect(false);
  WiFi.setAutoReconnect(true);
  WiFi.disconnect();
  WiFi.mode(WIFI_STA);

  DynamicJsonDocument wifiConfig(1024);
  DeserializationError err = deserializeJson(wifiConfig, readFile("/config/networks.json", 1024));
  if (err)
  {
    Serial.printf("could not deserialize wifi json config: %s\n", err.c_str());
    writeFile("/config/networks.json", "[]");
    Serial.println("reset the wifi configuration");
    deserializeJson(wifiConfig, "[]");
  }

  Serial.println("scanning networks...");
  WiFi.scanNetworks(false);

  Serial.println("adding known wifi networks to wifiMulti");
  for (JsonVariant network : wifiConfig.as<JsonArray>())
  {
    wifiMulti.addAP(network["ssid"].as<const char *>(), network["psk"].as<const char *>());
    Serial.printf("  * %s : %s\n", network["ssid"].as<const char *>(), network["psk"].as<const char *>());
  }

  WiFi.hostname(getHostname());
  Serial.println("");

  Serial.print("connecting...");

  wifiMulti.run(20000);

  Serial.println(getWiFiStatus());

  if (WiFi.status() != WL_CONNECTED)
  {
    Serial.println("could not connect to any known networks :'(");
    Serial.println("setting up captive portal...");
    setupCaptivePortal();
    dnsServer.setErrorReplyCode(DNSReplyCode::NoError);
    if (!dnsServer.start(DNS_PORT, "*", getCaptivePortalIP()))
    {
      Serial.println("failed to start DNS server on port 53");
    }
    else
    {
      Serial.println("started DNS server on port 53");
    }
  }
  else
  {
    Serial.print("connected to:\t");
    Serial.println(WiFi.SSID());
    Serial.print("IP address:\t");
    Serial.println(WiFi.localIP());
    Serial.print("MAC address:\t");
    Serial.println(WiFi.macAddress());
  }

  R = readColour("R");
  G = readColour("G");
  B = readColour("B");

  Serial.printf("LEDs:\t%d\n", NUM_LEDS);
  Serial.printf("MAX LEDs:\t%d\n", MAX_NUM_LEDS);
  Serial.printf("colours:\tR:%d G:%d B:%d\n", R, G, B);
  Serial.printf("brightness:\t%d\n", readBrightness());
  Serial.printf("chase trail:\t%d\n", chaseTrailLength);
  Serial.printf("chase direction:\t%d\n", chaseDirection);

  Serial.println();

  initOTA();
  initMDNS();

  digitalWrite(SIG_LED, HIGH);

  server.on(
      "/version", HTTP_GET, [](AsyncWebServerRequest *request)
      { request->send(200, "text/plain", readFile("/version", 128)); });

  server.on(
      "/api/status", HTTP_GET, [](AsyncWebServerRequest *request)
      {
        AsyncResponseStream *response = request->beginResponseStream("application/json");
        DynamicJsonDocument jsonBuffer(2048);
        JsonVariant root = jsonBuffer.as<JsonVariant>();
        root["free_heap"] = ESP.getFreeHeap();
        root["ssid"] = WiFi.SSID();
        root["mac"] = WiFi.macAddress();
        root["address"] = WiFi.localIP().toString();
        root["led_on"] = (digitalRead(LED_PIN) == HIGH) ? false : true;
        root["hostname"] = getHostname();
        root["leds"] = NUM_LEDS;
        root["wifi_status"] = getWiFiStatus();
        root["chip_id"] = getChipID();
        serializeJson(jsonBuffer, *response);
        request->send(response);
      });

  server.on(
      "/api/config", HTTP_GET, [](AsyncWebServerRequest *request)
      {
        AsyncResponseStream *response = request->beginResponseStream("application/json");
        DynamicJsonDocument jsonBuffer(2048);
        JsonVariant root = jsonBuffer.as<JsonVariant>();
        root["hostname"] = readFile("/config/hostname", 256);
        serializeJson(jsonBuffer, *response);
        request->send(response);
      });

  server.on(
      "/api/toggle-led", HTTP_GET, [](AsyncWebServerRequest *request)
      {
        AsyncResponseStream *response = request->beginResponseStream("application/json");
        DynamicJsonDocument jsonBuffer(2048);
        digitalWrite(LED_PIN, (digitalRead(LED_PIN) == 1) ? LOW : HIGH);
        JsonVariant root = jsonBuffer.as<JsonVariant>();
        root["led_on"] = (digitalRead(LED_PIN) == 1) ? false : true;
        serializeJson(jsonBuffer, *response);
        request->send(response);
      });

  server.on(
      "/api/colour", HTTP_GET, [](AsyncWebServerRequest *request)
      {
        char colour[7];
        sprintf(colour, "%02X%02X%02X", R, G, B);

        AsyncResponseStream *response = request->beginResponseStream("application/json");
        response->setCode(200);

        DynamicJsonDocument jsonBuffer(256);
        JsonVariant root = jsonBuffer.as<JsonVariant>();

        root["colour"] = colour;
        root["R"] = R;
        root["G"] = G;
        root["B"] = B;
        root["brightness"] = BRIGHTNESS;
        root["mode"] = OPERATION_MODE;
        serializeJson(jsonBuffer, *response);

        request->send(response);
      });

  server.on(
      "/api/chase", HTTP_GET, [](AsyncWebServerRequest *request)
      {
        AsyncResponseStream *response = request->beginResponseStream("application/json");
        response->setCode(200);

        DynamicJsonDocument jsonBuffer(256);
        JsonVariant root = jsonBuffer.as<JsonVariant>();

        root["speed"] = chaseSpeed;
        root["direction"] = chaseDirection;
        root["length"] = chaseTrailLength;
        serializeJson(jsonBuffer, *response);

        request->send(response);
      });

  server.on(
      "/api/ping", HTTP_GET, [](AsyncWebServerRequest *request)
      {
        AsyncResponseStream *response = request->beginResponseStream("application/json");
        response->setCode(200);

        DynamicJsonDocument jsonBuffer(64);
        JsonVariant root = jsonBuffer.as<JsonVariant>();

        root["ping"] = "pong";
        serializeJson(jsonBuffer, *response);

        request->send(response);
      });

  server.on(
      "/api/net/scan", HTTP_GET, [](AsyncWebServerRequest *request)
      {
        String json = "[";
        int n = WiFi.scanComplete();
        if (n == -2)
        {
          WiFi.scanNetworks(true);
        }
        else if (n)
        {
          for (int i = 0; i < n; ++i)
          {
            if (i)
              json += ",";
            json += "{";
            json += "\"rssi\":" + String(WiFi.RSSI(i));
            json += ",\"ssid\":\"" + WiFi.SSID(i) + "\"";
            json += ",\"bssid\":\"" + WiFi.BSSIDstr(i) + "\"";
            json += ",\"channel\":" + String(WiFi.channel(i));
            json += ",\"secure\":" + String(WiFi.encryptionType(i));
            json += ",\"hidden\":" + String(WiFi.isHidden(i) ? "true" : "false");
            json += "}";
          }
          WiFi.scanDelete();
          if (WiFi.scanComplete() == -2)
          {
            WiFi.scanNetworks(true);
          }
        }
        json += "]";
        request->send(200, "application/json", json);
        json = String();
      });

  server.on(
      "/api/net/list", HTTP_GET, [](AsyncWebServerRequest *request)
      {
        StaticJsonDocument<1024> jsonBuffer;
        DeserializationError err = deserializeJson(jsonBuffer, readFile("/config/networks.json", 1024));

        if (err)
        {
          Serial.print("failed to deserialize stored network list: ");
          Serial.println(err.c_str());
          return jsonError(request, 500, err.c_str());
        }

        DynamicJsonDocument jsonResponse(1024);
        JsonArray storedConfig = jsonBuffer.as<JsonArray>();

        for (JsonVariant value : storedConfig)
        {
          jsonResponse.add(value["ssid"]);
        }

        AsyncResponseStream *response = request->beginResponseStream("application/json");
        serializeJson(jsonResponse, *response);
        request->send(response);
      });

  server.on("/", HTTP_GET, [](AsyncWebServerRequest *request)
            { request->redirect("/index.html"); });

  server.on("/generate_204", HTTP_GET, [](AsyncWebServerRequest *request)
            { request->send(204, "text/plain", "hello"); });

  server.on("/fwlink", HTTP_GET, [](AsyncWebServerRequest *request)
            { request->send(204, "text/plain", "hello"); });

  AsyncCallbackJsonWebHandler *networkAddHandler = new AsyncCallbackJsonWebHandler(
      "/api/net/add", [](AsyncWebServerRequest *request, JsonVariant &json)
      {
        if (json["ssid"] == nullptr || json["psk"] == nullptr)
        {
          return jsonError(request, 400, "missing 'psk' or 'ssid' key");
        }

        DynamicJsonDocument config(1024);
        DeserializationError err = deserializeJson(config, readFile("/config/networks.json", 1024));

        if (err)
        {
          Serial.print("failed to deserialize stored network list: ");
          Serial.println(err.c_str());
          return jsonError(request, 500, err.c_str());
        }

        bool isNewNetwork = true;

        for (JsonVariant network : config.as<JsonArray>())
        {
          if (network["ssid"].as<String>().equals(json["ssid"].as<String>()))
          {
            network["psk"] = json["psk"];
            isNewNetwork = false;
            break;
          }
        }

        if (config.as<JsonArray>().size() >= 5 && isNewNetwork)
        {
          return jsonError(request, 400, "you cannot save more than 5 networks");
        }

        if (isNewNetwork)
        {
          StaticJsonDocument<256> newSSID;

          newSSID["ssid"] = json["ssid"];
          newSSID["psk"] = json["psk"];

          config.add(newSSID);
          Serial.printf("adding network to the config %s\n", newSSID["ssid"].as<String>().c_str());
        }

        File configFile = LittleFS.open("/config/networks.json", "w");
        if (serializeJson(config, configFile) == 0)
        {
          Serial.println("failed to save the network configuration");
          configFile.close();
          return jsonError(request, 500, "failed to save the configuration to file");
        }
        configFile.close();

        return jsonSuccess(request, 200, "successfully added network");
      });

  AsyncCallbackJsonWebHandler *networkDelHandler = new AsyncCallbackJsonWebHandler(
      "/api/net/del", [](AsyncWebServerRequest *request, JsonVariant &json)
      {
        if (json["ssid"] == nullptr)
        {
          return jsonError(request, 400, "missing 'ssid' key");
        }

        DynamicJsonDocument config(1024);
        DeserializationError err = deserializeJson(config, readFile("/config/networks.json", 1024));

        if (err)
        {
          Serial.print("failed to deserialize stored network list: ");
          Serial.println(err.c_str());
          return jsonError(request, 500, err.c_str());
        }

        for (JsonArray::iterator network = config.as<JsonArray>().begin(); network != config.as<JsonArray>().end(); ++network)
        {
          if ((*network)["ssid"].as<String>().equals(json["ssid"].as<String>()))
          {
            Serial.printf("removed network from config %s\n", (*network)["ssid"].as<String>().c_str());
            config.as<JsonArray>().remove(network);
          }
        }

        File configFile = LittleFS.open("/config/networks.json", "w");
        if (serializeJson(config, configFile) == 0)
        {
          Serial.println("failed to save the network configuration");
          configFile.close();
          return jsonError(request, 500, "failed to save the configuration to file");
        }
        configFile.close();

        return jsonSuccess(request, 200, "successfully removed network");
      });

  AsyncCallbackJsonWebHandler *pingHandler = new AsyncCallbackJsonWebHandler(
      "/api/ping", [](AsyncWebServerRequest *request, JsonVariant &json)
      {
        String response;
        serializeJson(json["data"], response);
        request->send(200, "application/json", response);
      });

  AsyncCallbackJsonWebHandler *networkResetHandler = new AsyncCallbackJsonWebHandler(
      "/api/net/reset", [](AsyncWebServerRequest *request, JsonVariant &json)
      {
        if (json["reset"] != nullptr && json["reset"].as<bool>())
        {
          writeFile("/config/networks.json", "[]");
          return jsonSuccess(request, 200, "successfully wiped network configuration");
        }

        return jsonError(request, 400, "you must set 'reset to true'");
      });

  AsyncCallbackJsonWebHandler *chaseSetHandler = new AsyncCallbackJsonWebHandler(
      "/api/chase", [](AsyncWebServerRequest *request, JsonVariant &json)
      {
        bool persist = isPersist(json);
        if (json["direction"] != nullptr)
        {
          int d = json["direction"].as<int>();
          if (d != 1 && d != -1)
          {
            return jsonError(request, 400, "speed should be 1 or -1");
          }
          chaseDirection = d;
          if (persist)
          {
            writeIntFile("/config/chase/direction", d);
          }
          Serial.printf("chase direction set to %d\n", d);
        }

        if (json["length"] != nullptr)
        {
          int length = json["length"].as<int>();
          if (length < 0 || length > NUM_LEDS - 1)
          {
            return jsonError(request, 400, "length should be between 0 and NUM_LEDS");
          }
          chaseTrailLength = length;
          if (persist)
          {
            writeIntFile("/config/chase/trail", length);
          }
          Serial.printf("chase trail set to %d\n", length);
        }

        if (json["speed"] != nullptr)
        {
          int speed = json["speed"].as<int>();
          if (speed < 1 || speed > 10)
          {
            return jsonError(request, 400, "speed should be between 1 and 10");
          }
          chaseSpeed = speed;
          if (persist)
          {
            writeIntFile("/config/chase/speed", speed);
          }
          Serial.printf("chase speed set to %d\n", speed);
        }

        return jsonSuccess(request, 200, "successfully changed chase configuration");
      });

  AsyncCallbackJsonWebHandler *colourSetHandler = new AsyncCallbackJsonWebHandler(
      "/api/colour", [](AsyncWebServerRequest *request, JsonVariant &json)
      {
        bool persist = isPersist(json);

        if (json["colour"] != nullptr)
        {
          String colour = json["colour"].as<String>();

          Serial.printf("new colour: %s\n", colour.c_str());

          R = (int)strtol(colour.substring(0, 2).c_str(), NULL, 16);
          G = (int)strtol(colour.substring(2, 4).c_str(), NULL, 16);
          B = (int)strtol(colour.substring(4, 6).c_str(), NULL, 16);

          if (persist)
          {
            saveColour("R", R);
            saveColour("G", G);
            saveColour("B", B);
          }

          Serial.printf("R:%d G:%d B:%d\n", R, G, B);
        }

        if (json["brightness"] != nullptr)
        {
          int brightness = json["brightness"].as<int>();
          if (brightness < 0 || brightness > 256)
          {
            return jsonError(request, 400, "you must set 'brightness' to a correct value");
          }
          if (persist)
          {
            saveBrightness(brightness);
          }
          BRIGHTNESS = brightness;
        }

        if (json["mode"] != nullptr)
        {
          String mode = json["mode"].as<String>();
          if (mode == "chase" || mode == "static" || mode == "breathing")
          {
            if (persist)
            {
              writeFile("/config/mode", mode);
            }
            OPERATION_MODE = mode;
          }
          else
          {
            return jsonError(request, 400, "invalid mode " + mode);
          }
        }

        return jsonSuccess(request, 200, "successfully changed colour mode");
      });

  AsyncCallbackJsonWebHandler *systemSetHandler = new AsyncCallbackJsonWebHandler(
      "/api/system", [](AsyncWebServerRequest *request, JsonVariant &json)
      {
        if (json["num_leds"] != nullptr)
        {
          int num_leds = json["num_leds"].as<int>();
          if (num_leds < 0 || num_leds > MAX_NUM_LEDS)
          {
            return jsonError(request, 400, "you must set 'num_leds' to a correct value");
          }
          writeIntFile("/config/leds", num_leds);
        }

        if (json["hostname"] != nullptr)
        {
          String hostname = json["hostname"].as<String>();
          writeFile("/config/hostname", hostname);
        }

        return jsonSuccess(request, 200, "successfully changed system infos, you'll need a reboot");
      });
  
  AsyncCallbackJsonWebHandler *rebootHandler = new AsyncCallbackJsonWebHandler(
      "/api/reboot", [](AsyncWebServerRequest *request, JsonVariant &json)
      {
        if (json["reboot"] != nullptr)
        {
          bool reboot = json["reboot"].as<bool>();
          if(reboot) {
            ASKED_FOR_REBOOT=true;
          } else {
            return jsonSuccess(request, 200, "no 'reboot' field, no reboot");
          }
        }

        return jsonSuccess(request, 200, "asked for a reboot");
      });

  server.addHandler(pingHandler);
  server.addHandler(networkAddHandler);
  server.addHandler(networkDelHandler);
  server.addHandler(networkResetHandler);

  server.addHandler(colourSetHandler);
  server.addHandler(systemSetHandler);
  server.addHandler(chaseSetHandler);
  server.addHandler(rebootHandler);
  server.addHandler(pingHandler);

  server.serveStatic("/index.html", LittleFS, "/static/index.html");
  server.serveStatic("/bootstrap.js", LittleFS, "/static/bootstrap.js").setCacheControl("max-age=3600");
  server.serveStatic("/jquery.js", LittleFS, "/static/jquery.js").setCacheControl("max-age=3600");
  server.serveStatic("/notify.js", LittleFS, "/static/notify.js").setCacheControl("max-age=3600");
  server.serveStatic("/bootstrap.css", LittleFS, "/static/bootstrap.css").setCacheControl("max-age=3600");
  server.serveStatic("/favicon.ico", LittleFS, "/static/favicon.ico").setCacheControl("max-age=3600");
  server.serveStatic("/static/", LittleFS, "/static/").setCacheControl("max-age=600");

  server.begin();

  WiFi.scanNetworks(true);
}

// TODO: make the variables names more sensible
int increment = 1;
int base = 0;
int maxVal = 255;
int val = base;

void loop()
{
  MDNS.update();
  ArduinoOTA.handle();
  dnsServer.processNextRequest();

  if (FastLED.getBrightness() != BRIGHTNESS)
  {
    FastLED.setBrightness(BRIGHTNESS);
  }

  if (OPERATION_MODE == "chase")
  {
    if (chasePauseCounter == 0)
    {
      FastLED.clear();
      leds[chaseLeadPosition] = CRGB(R, G, B);
      for (int i = 1; i < chaseTrailLength; i++)
      {
        double factor = double(chaseTrailLength - i) / double(chaseTrailLength);
        leds[mod(chaseLeadPosition - chaseDirection * i, NUM_LEDS)] = CRGB(R * factor, G * factor, B * factor);
      }

      chaseLeadPosition = (chaseLeadPosition + chaseDirection) % NUM_LEDS;
      if (chaseLeadPosition < 0)
      {
        chaseLeadPosition = NUM_LEDS - 1;
      }
    }
    chasePauseCounter = (chasePauseCounter + 1) % chaseSpeed;
  }
  else if (OPERATION_MODE == "static")
  {
    for (int i = 0; i < NUM_LEDS; i++)
    {
      leds[i] = CRGB(R, G, B);
    }
  }
  else
  {
    double factor = double(val) / double(maxVal);

    int r = (int)map(R * factor, 0, 256, 0, 230);
    int g = (int)map(G * factor, 0, 256, 0, 230);
    int b = (int)map(B * factor, 0, 256, 0, 230);

    for (int i = 0; i < NUM_LEDS; i++)
    {
      leds[i] = CRGB(r, g, b);
    }

    val += increment;
    if (val == base || val == (maxVal))
    {
      increment *= -1;
    }
  }

  FastLED.show();
  FastLED.delay(1000 / FRAMES_PER_SECOND);

  if (ASKED_FOR_REBOOT)
  {
    REBOOT_TICKS--;
    if (REBOOT_TICKS <= 0)
    {
      softReset();
    }
  }
}

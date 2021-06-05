#include <Arduino.h>
#include <ESP8266WiFi.h>
#include <DNSServer.h>
#include "captivePortal.h"
#include "misc.h"

IPAddress captivePortalIP(192, 168, 69, 1);
IPAddress captivePortalNetmask(255, 255, 255, 0);

void setupCaptivePortal()
{
    WiFi.setAutoConnect(false);
    WiFi.setAutoReconnect(true);
    WiFi.disconnect();
    WiFi.mode(WIFI_AP);
    wifi_set_sleep_type(NONE_SLEEP_T);

    WiFi.softAPConfig(captivePortalIP, captivePortalIP, captivePortalNetmask);
    WiFi.softAP("ringos-" + getChipID());
    Serial.println("starting network ringos-" + getChipID());
    delay(500);
    Serial.print("Software access point IP address:\t");
    Serial.println(WiFi.softAPIP());
    Serial.print("MAC address:\t");
    Serial.println(WiFi.macAddress());
}

IPAddress getCaptivePortalIP()
{
    return captivePortalIP;
}
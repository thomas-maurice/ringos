#include <Arduino.h>
#include <ArduinoOTA.h>

#include "otaUtil.h"
#include "misc.h"

#define OTA_PASSWORD "esp8266"

void initOTA()
{
    ArduinoOTA.setHostname(getHostname().c_str());
    ArduinoOTA.setPassword(OTA_PASSWORD);

    ArduinoOTA.onStart(
        []()
        {
            Serial.println("OTA update start");
            delay(100);
        });

    ArduinoOTA.onEnd(
        []()
        {
            Serial.println("\nOTA update end");
            digitalWrite(LED_BUILTIN, HIGH);
            delay(100);
        });

    ArduinoOTA.onProgress(
        [](unsigned int progress, unsigned int total)
        {
            Serial.printf("OTA progress: %u%%\r", (progress / (total / 100)));
            digitalWrite(LED_BUILTIN, (progress / (total / 100)) % 2);
        });

    ArduinoOTA.onError(
        [](ota_error_t error)
        {
            Serial.printf("Error[%u]: ", error);
            if (error == OTA_AUTH_ERROR)
                Serial.println("Auth Failed");
            else if (error == OTA_BEGIN_ERROR)
                Serial.println("Begin Failed");
            else if (error == OTA_CONNECT_ERROR)
                Serial.println("Connect Failed");
            else if (error == OTA_RECEIVE_ERROR)
                Serial.println("Receive Failed");
            else if (error == OTA_END_ERROR)
                Serial.println("End Failed");
            digitalWrite(BUILTIN_LED, 1);
            delay(100);
        });

    ArduinoOTA.begin();
    Serial.println("OTA updater ready");
}
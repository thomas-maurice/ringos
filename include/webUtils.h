#ifndef _H_WEBUTILS
#define _H_WEBUTILS

#include <Arduino.h>
#include <ESPAsyncTCP.h>
#include <ESPAsyncWebServer.h>
#include <AsyncJson.h>
#include <ArduinoJson.h>

void jsonError(AsyncWebServerRequest *request, int httpCode, String errorMessage, String status = "failure");
void jsonSuccess(AsyncWebServerRequest *request, int httpCode, String message, String status = "success");

#endif
#include <Arduino.h>
#include <ESPAsyncTCP.h>
#include <ESPAsyncWebServer.h>
#include <AsyncJson.h>
#include <ArduinoJson.h>

#include "webUtils.h"

void jsonError(AsyncWebServerRequest *request, int httpCode, String errorMessage, String status) {
    AsyncResponseStream *response = request->beginResponseStream("application/json");
    response->setCode(httpCode);

    DynamicJsonDocument jsonError(512);
    jsonError["error"] = errorMessage;
    jsonError["status"] = status;
    serializeJson(jsonError, *response);
    request->send(response);
}

void jsonSuccess(AsyncWebServerRequest *request, int httpCode, String message, String status) {
    AsyncResponseStream *response = request->beginResponseStream("application/json");
    response->setCode(httpCode);

    DynamicJsonDocument jsonError(512);
    jsonError["message"] = message;
    jsonError["status"] = status;
    serializeJson(jsonError, *response);
    request->send(response);
}
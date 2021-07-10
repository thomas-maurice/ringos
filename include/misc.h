#ifndef _MISC_H
#define _MISC_H

#include <Arduino.h>
#include <ArduinoJson.h>

String getHostname();
String getChipID();
int mod(int a, int b);
bool isPersist(JsonVariant &json);

#endif
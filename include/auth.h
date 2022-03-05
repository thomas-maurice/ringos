#ifndef _AUTH_H
#define _AUTH_H

#include <Arduino.h>
#include <ESPAsyncWebServer.h>

bool needsAuth(String password);

#endif
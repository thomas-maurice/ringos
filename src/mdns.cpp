#include <Arduino.h>
#include <ESP8266mDNS.h>
#include "fsUtil.h"
#include "misc.h"

bool initMDNS()
{
    String hostname = getHostname();
    MDNS.setHostname(hostname);
    MDNS.addService("http", "tcp", 80);
    return MDNS.begin(hostname);
}
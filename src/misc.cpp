#include <Arduino.h>
#include "fsUtil.h"
#include "misc.h"

// returns the hostname
String getHostname()
{
    String hostname = readFile("/config/hostname", 128);
    if (hostname == "")
    {
        return "jar-" + getChipID();
    }
    return hostname;
}

// gets the chip ID
String getChipID()
{
    char espID[9];
    itoa(ESP.getChipId(), espID, 16);
    String espIDString = String(espID);
    espIDString.toUpperCase();
    return espIDString;
}

// modulus function
int mod(int a, int b)
{
    int r = a % b;
    return r < 0 ? r + b : r;
}
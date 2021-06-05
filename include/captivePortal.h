#ifndef _H_CAPTIVE
#define _H_CAPTIVE

#include <DNSServer.h>

#define DNS_PORT 53

void setupCaptivePortal();
IPAddress getCaptivePortalIP();

#endif
#ifndef _FS_UTIL_H
#define _FS_UTIL_H

#include <Arduino.h>
#include <LittleFS.h>

String readFile(String fname, int maxSize);
int readIntFile(String fname, int maxSize);
int writeFile(String fname, String data);
int writeIntFile(String fname, int data);

#endif
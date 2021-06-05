#include <Arduino.h>
#include <Crypto.h>
#include <SHA256.h>

#include "hash.h"

#define HASH_SIZE_BYTES 32

// Hash hashes a string with the default SHA256 algorithm
String Hash(String input)
{
    SHA256 sha256;
    sha256.update(input.c_str(), input.length());
    char buf[HASH_SIZE_BYTES];
    sha256.finalize(&buf, 32);
    String output;
    for (int i = 0; i < HASH_SIZE_BYTES; i++)
    {
        char hex[2];
        sprintf(hex, "%x", buf[i]);
        output += hex;
    }

    return output;
}
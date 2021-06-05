#include <Arduino.h>
#include <LittleFS.h>

// Reads a file from the filesystem
String readFile(String fname, int maxSize)
{
  File file = LittleFS.open(fname, "r");
  if (!file)
  {
    Serial.println("Failed to open file " + fname + " for reading");
    return "";
  }
  String result("");
  while (file.available() && maxSize > 0)
  {
    result += char(file.read());
    maxSize--;
  }
  file.close();
  return result;
}

// Reads an integer from the filesystem
int readIntFile(String fname, int maxSize)
{
  String sResult = readFile(fname, maxSize);
  int result;
  sscanf(sResult.c_str(), "%d", &result);
  return result;
}

// Writes data to the filesystem
int writeFile(String fname, String data)
{
  File file = LittleFS.open(fname, "w");

  if (!file)
  {
    Serial.println("There was an error opening " + fname + " for writing");
    file.close();
    return -1;
  }
  if (!file.print(data))
  {
    Serial.println("File write failed on " + fname);
    file.close();
    return -1;
  }
  file.close();
  return 0;
}

// Writes integer data to the filesystem
int writeIntFile(String fname, int data)
{
  File file = LittleFS.open(fname, "w");

  if (!file)
  {
    Serial.println("There was an error opening " + fname + " for writing");
    file.close();
    return -1;
  }
  char stringData[16];
  sprintf(stringData, "%d", data);
  if (!file.print(stringData))
  {
    Serial.println("File write failed on " + fname);
    file.close();
    return -1;
  }
  file.close();
  return 0;
}
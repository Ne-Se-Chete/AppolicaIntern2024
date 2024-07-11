#include <GxEPD.h>
#include <GxGDEH0154D67/GxGDEH0154D67.h>  // 1.54" b/w 200x200, SSD1681
#include GxEPD_BitmapExamples
#include <Fonts/FreeSansBold18pt7b.h>
#include <GxIO/GxIO_SPI/GxIO_SPI.h>
#include <GxIO/GxIO.h>

GxIO_Class io(SPI, /*CS=5*/ SS, /*DC=*/ 17, /*RST=*/ 16);
GxEPD_Class display(io, /*RST=*/ 16, /*BUSY=*/ 4);

void setup() {
  Serial.begin(115200);
  Serial.println();
  Serial.println("setup");

  display.init(115200); 
  drawHelloWorld();
  display.update();

  Serial.println("setup done");
  
  // display.powerDown(); // Ne se chete s tova
}

void loop() {};

const char HelloWorld[] = "Hello World!";

void drawHelloWorld() {
  display.setRotation(1);
  display.setFont(&FreeSansBold18pt7b);
  display.setTextColor(GxEPD_BLACK);
  int16_t tbx, tby; uint16_t tbw, tbh;
  display.getTextBounds(HelloWorld, 0, 0, &tbx, &tby, &tbw, &tbh);
  // center bounding box by transposition of origin:
  uint16_t x = ((display.width() - tbw) / 2) - tbx;
  uint16_t y = ((display.height() - tbh) / 2) - tby;
  display.fillScreen(GxEPD_WHITE);
  display.setCursor(x, y);
  display.print(HelloWorld);
}

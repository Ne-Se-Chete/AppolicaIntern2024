#include <HX711_ADC.h>
#include <EEPROM.h>
#define EEPROM_SIZE 16
#include <GxEPD.h>
#include <GxGDEH0154D67/GxGDEH0154D67.h>  // 1.54" b/w 200x200, SSD1681
#include GxEPD_BitmapExamples
#include <Fonts/FreeSansBold18pt7b.h>
#include <GxIO/GxIO_SPI/GxIO_SPI.h>
#include <GxIO/GxIO.h>

GxIO_Class io(SPI, /*CS=5*/ SS, /*DC=*/ 17, /*RST=*/ 16);
GxEPD_Class display(io, /*RST=*/ 16, /*BUSY=*/ 4);


//pins:
const int HX711_dout = 22; //mcu > HX711 dout pin
const int HX711_sck = 15; //mcu > HX711 sck pin

//HX711 constructor:
HX711_ADC LoadCell(HX711_dout, HX711_sck);

const int calVal_eepromAdress = 0;
unsigned long t = 0;
float initialMeasurement = 0;
int readingsForAverage = 50;
const int serialPrintInterval = 200; //increase value to slow down serial print activity
bool cooking = false;


void setup() {
  Serial.begin(115200); 
  delay(10);
  Serial.println();
  Serial.println("Starting...");
  startLoadCell();
  initialMeasurement = LoadCell.getData();
}

void loop() {
  float avrgReading = averageReading(readingsForAverage);
  if(!cooking) {
    if(initialMeasurement - avrgReading > 20) {
    delay(5000);
    if(initialMeasurement - avrgReading > 20)
      cooking = true;
  }
  }
  
  int hasStoppedCooking = 0;
  int newWeight = avrgReading;
  while(cooking) {
    //5,6 gr/min
    delay(30000);
    int tmpWeight = averageReading(readingsForAverage);
    if(newWeight - tmpWeight < 3)
      hasStoppedCooking++;
    else
      newWeight = avrgReading;

    if(hasStoppedCooking > 1)
      cooking = false;
  }


}

void calibrate() {
  Serial.println("***");
  Serial.println("Start calibration:");
  Serial.println("Place the load cell an a level stable surface.");
  Serial.println("Remove any load applied to the load cell.");
  Serial.println("Send 't' from serial monitor to set the tare offset.");

  boolean _resume = false;
  while (_resume == false) {
    LoadCell.update();
    if (Serial.available() > 0) {
      if (Serial.available() > 0) {
        char inByte = Serial.read();
        Serial.println(LoadCell.getTareOffset());
        if (inByte == 't') LoadCell.tareNoDelay();
        delay(2000);
      }
    }
    if (LoadCell.getTareStatus() == true) {
      EEPROM.writeLong(9, LoadCell.getTareOffset());
      Serial.println("Tare complete");
      _resume = true;
    }
  }

  Serial.println("Now, place your known mass on the loadcell.");
  Serial.println("Then send the weight of this mass (i.e. 100.0) from serial monitor.");

  float known_mass = 0;
  _resume = false;
  while (_resume == false) {
    LoadCell.update();
    if (Serial.available() > 0) {
      known_mass = Serial.parseFloat();
      if (known_mass != 0) {
        Serial.print("Known mass is: ");
        Serial.println(known_mass);
        _resume = true;
      }
    }
  }

  LoadCell.refreshDataSet(); //refresh the dataset to be sure that the known mass is measured correct
  float newCalibrationValue = LoadCell.getNewCalibration(known_mass); //get the new calibration value

  Serial.print("New calibration value has been set to: ");
  Serial.print(newCalibrationValue);
  Serial.println(", use this as calibration value (calFactor) in your project sketch.");
  Serial.print("Save this value to EEPROM adress ");
  Serial.print(calVal_eepromAdress);
  Serial.println("? y/n");

  _resume = false;
  while (_resume == false) {
    if (Serial.available() > 0) {
      char inByte = Serial.read();
      if (inByte == 'y') {
        // EEPROM.begin(EEPROM_SIZE);
        EEPROM.writeFloat(calVal_eepromAdress, newCalibrationValue);
        EEPROM.commit();
        newCalibrationValue = EEPROM.readFloat(calVal_eepromAdress);
        Serial.print("Value ");
        Serial.print(newCalibrationValue);
        Serial.print(" saved to EEPROM address: ");
        Serial.println(calVal_eepromAdress);
        _resume = true;

      }
      else if (inByte == 'n') {
        Serial.println("Value not saved to EEPROM");
        _resume = true;
      }
    }
  }

  Serial.println("End calibration");
  Serial.println("***");
  Serial.println("To re-calibrate, send 'r' from serial monitor.");
  Serial.println("For manual edit of the calibration value, send 'c' from serial monitor.");
  Serial.println("***");
}

void changeFromSavedCalFactor() {
  float oldCalibrationValue = LoadCell.getCalFactor();
  long int tareValue = 0;

  Serial.println("***");
  Serial.print("Current value is: ");
  Serial.println(oldCalibrationValue);
  float newCalibrationValue;
  newCalibrationValue = EEPROM.readFloat(calVal_eepromAdress);
  tareValue = EEPROM.readLong(9);
  Serial.println(tareValue);
  Serial.print("New calibration factor is set to: ");
  Serial.println(newCalibrationValue);
  LoadCell.setCalFactor(newCalibrationValue);
  LoadCell.setTareOffset(tareValue);
  Serial.print(LoadCell.getTareOffset());
  delay(2000);
  Serial.println("End change calibration value");
  Serial.println("***");
}

void startLoadCell() {
  EEPROM.begin(EEPROM_SIZE);
  int cal_val_eeprom = 0;

  LoadCell.begin();
  //LoadCell.setReverseOutput(); //uncomment to turn a negative output value to positive
  unsigned long stabilizingtime = 2000; // preciscion right after power-up can be improved by adding a few seconds of stabilizing time
  boolean _tare = false; //set this to false if you don't want tare to be performed in the next step
  LoadCell.start(stabilizingtime, _tare);
  if (LoadCell.getTareTimeoutFlag() || LoadCell.getSignalTimeoutFlag()) {
    Serial.println("Timeout, check MCU>HX711 wiring and pin designations");
    while (1);
  }
  else {
    cal_val_eeprom = EEPROM.readFloat(calVal_eepromAdress);
    if(cal_val_eeprom != 0) {
        changeFromSavedCalFactor();
        Serial.println("Startup is complete");
    }
    else {
      Serial.println("The sensors should be calibrated!");
      while (!LoadCell.update());
        calibrate();
    }
  }
}

float averageReading(int timesToRead) {
  float sum = 0;
  static boolean newDataReady = 0;
  for(int i = 0; i < timesToRead; i++) {
    if (LoadCell.update()) newDataReady = true;

  // get smoothed value from the dataset:
    if (newDataReady) {
      if (millis() > t + serialPrintInterval) {
        float measurement = LoadCell.getData();
        sum += measurement;
        Serial.print("Load_cell output val: ");
        Serial.println(measurement);
        newDataReady = 0;
        t = millis();
      }
    }
  }

  return sum/timesToRead;
}

  // // receive command from serial terminal
  // if (Serial.available() > 0) {
  //   char inByte = Serial.read();
  //   if (inByte == 't') LoadCell.tareNoDelay(); //tare
  //   else if (inByte == 'r') calibrate(); //calibrate
  // }

  // // check if last tare operation is complete
  // if (LoadCell.getTareStatus() == true) {
  //   Serial.println("Tare complete");
  // }
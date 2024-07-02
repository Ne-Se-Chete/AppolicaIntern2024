#include <WiFi.h>
#include <HTTPClient.h>
#include <NTPClient.h> //install this library
#include <WiFiUdp.h>


// Replace these with your Wi-Fi credentials and API endpoint
const char* ssid = "##############";
const char* password = "##############";
const char* testServerName = "##############/test";
const char* usersServerName = "##############/Users";
const char* grillServerName = "##############/Grill";

const char* bearerToken = "##############";


WiFiUDP ntpUDP;
NTPClient timeClient(ntpUDP, "pool.ntp.org", 0, 60000);  // Update every 60 seconds


void setup() {
  Serial.begin(115200);

  // Connect to Wi-Fi
  WiFi.begin(ssid, password);
  Serial.print("Connecting to Wi-Fi");
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println();
  Serial.println("Connected to Wi-Fi");

  timeClient.begin();
  while(!timeClient.update()) {
    timeClient.forceUpdate();
  }


  // Call each function to send data to the respective endpoints
  // post_to_test_endpoint();
  // post_to_users_endpoint();
  // post_to_grill_endpoint();
  post_to_grill_endpoint(get_user_id("Pencho"));
}

String get_user_id(String userName) {
  String user_id = "";

  // Create the URL for the GET request
  String url = String(usersServerName) + "?name=" + userName;

  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;
    http.begin(url);
    http.addHeader("Authorization", String("Bearer ") + bearerToken);

    int httpResponseCode = http.GET();

    if (httpResponseCode > 0) {
      String response = http.getString();
      Serial.println(httpResponseCode);
      Serial.println(response);

      // Parse the JSON response to find the _id associated with the user name
      int idIndex = response.indexOf("\"_id\":");
      if (idIndex != -1) {
        // Move the index past "_id":" and the initial quote
        idIndex = response.indexOf("\"", idIndex + 6) + 1; 
        int idEndIndex = response.indexOf("\"", idIndex);
        user_id = response.substring(idIndex, idEndIndex);
      } else {
        Serial.println("Error: _id not found in response.");
      }
    } else {
      Serial.print("Error on sending GET: ");
      Serial.println(httpResponseCode);
    }

    http.end();
  } else {
    Serial.println("Error in WiFi connection");
  }

  return user_id;
}

void post_to_test_endpoint() {
  // Prepare the data to send
  String test_option = "testing with esp32";

  // Create the JSON object
  String jsonData = String("{\"test_option\":\"") + test_option + "\"}";

  // Send the data to the server
  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;
    http.begin(testServerName);
    http.addHeader("Content-Type", "application/json");
    http.addHeader("Authorization", String("Bearer ") + bearerToken);

    int httpResponseCode = http.POST(jsonData);

    if (httpResponseCode > 0) {
      String response = http.getString();
      Serial.println(httpResponseCode);
      Serial.println(response);
    } else {
      Serial.print("Error on sending POST: ");
      Serial.println(httpResponseCode);
    }

    http.end();
  } else {
    Serial.println("Error in WiFi connection");
  }
}

void post_to_users_endpoint() {
  // Prepare the data to send
  String user_name = "Pencho";
  int points_cooking = 2;

  // Create the JSON object
  String jsonData = String("{\"name\":\"") + user_name + "\",\"points cooking\":" + points_cooking + "}";

  // Send the data to the server
  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;
    http.begin(usersServerName);
    http.addHeader("Content-Type", "application/json");
    http.addHeader("Authorization", String("Bearer ") + bearerToken);

    int httpResponseCode = http.POST(jsonData);

    if (httpResponseCode > 0) {
      String response = http.getString();
      Serial.println(httpResponseCode);
      Serial.println(response);
    } else {
      Serial.print("Error on sending POST: ");
      Serial.println(httpResponseCode);
    }

    http.end();
  } else {
    Serial.println("Error in WiFi connection");
  }
}

void post_to_grill_endpoint(String user_griled) {
  // Prepare the data to send
  float grill_start_gas = 8.0;
  float grill_end_gas = 5.2;

  // Format the current time and end time
  time_t now = timeClient.getEpochTime();
  char end_time[20];
  char start_time[20];
  strftime(end_time, sizeof(end_time), "%Y-%m-%d %H:%M:%S", localtime(&now));
  now -= 20 * 60;  // Subtracting 20 minutes
  strftime(start_time, sizeof(start_time), "%Y-%m-%d %H:%M:%S", localtime(&now));

  // Create the JSON object
  String jsonData = String("{\"grill start gas\":") + grill_start_gas +
                    ",\"grill end gas\":" + grill_end_gas +
                    ",\"user grilled\":\"" + user_griled + "\"" +
                    ",\"start time\":\"" + start_time + "\"" +
                    ",\"end time\":\"" + end_time + "\"}";

  // Send the data to the server
  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;
    http.begin(grillServerName);
    http.addHeader("Content-Type", "application/json");
    http.addHeader("Authorization", String("Bearer ") + bearerToken);

    int httpResponseCode = http.POST(jsonData);

    if (httpResponseCode > 0) {
      String response = http.getString();
      Serial.println(httpResponseCode);
      Serial.println(response);
    } else {
      Serial.print("Error on sending POST: ");
      Serial.println(httpResponseCode);
    }

    http.end();
  } else {
    Serial.println("Error in WiFi connection");
  }
}



void loop() {
  // Put your main code here, to run repeatedly
}

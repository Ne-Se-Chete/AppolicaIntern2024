# AppolicaIntern2024

### Project Description

This repository is for a project aimed at enhancing the management and tracking of BBQ proccess. The application will feature comprehensive notifications, historical data tracking, and estimation tools to optimize the cooking and consumption experience. It will include integrations with Slack, a webpage and Slack bot for order management, and additional predictive and recommendation features. The goal is to provide the user how many kgs are left in the gas tank and when it must be refiled.

### Features of the project

1. [ ] Measure the weight of the bottle
2. [ ] Create an estimation of how long the bottle will last based on:
   - [ ] Which ones have been rotated
   - [ ] How many rotations
3. [ ] Track the history of the cooking
4. [ ] Set up notifications (when to refill, if you've forgotten)

### Bonus Tasks

- **Slack Bot:**
  - Dependent on:
    - [ ] Hard-coded cooking times
  - Features:
    - [ ] Predict beer consumption
    - [ ] Manage orders

- **Scoreboard:**
  - Dependent on:
    - [ ] Log in functionality
  - Features:
    - [ ] Master Chef leaderboard

- **Beer Prediction:**
  - Implement a system for predicting beer needs and preferences based on historical data and user behavior


## Day 1
### To Do
1. <del> To meet our mentors <del>
2. <del> To be assigned a project <del>
3. <del> To divide into groups and to see who will do what <del>
4. <del> To describe the technologies we will be using <del>
5. <del> To make a roadmap with all of our plans <del>
6. <del> To make a plan that we will follow daily <del>

### What we did
1. We met our mentors - they are awesome. We met many people from the company and all of them were friendly

2. We will be working on a project called "Smart Scale". The idea is described below/higher

3. For the software Bobi and Moni will be resposible and for the hardware - Yaskata and Valeri

4. We will be using:
- Bubble io for the frontend, backend and the DB
- Go(lang) for the Bot
- ESP32 for communication with the smart scale and the backend
- Load cell three wire sensors for the measurment
5. Our roadmap is ![here](https://github.com/Ne-Se-Chete/AppolicaIntern2024/blob/main/images/roadmap.png)

6. We have made a daily plan which includes:
- daily meetings at 10:00
- following the plan
- helping eachother when we have problems
- every night to write a short description of what we did


## Day 2
### To Do
1. <del> To have a daily meeting<del>
2. <del>To setup the scale<del>
3. <del>To do testing with the Bubble io<del>
4. <del>To fix the DB schematic<del>
5. <del>To discuss further the idea for the software<del>

### What we did
1. We had daily meeting on which we discussed the plan for the day and we made some changes
2. We find a bug with the sensors which we fixed
3. We made the DB and we tested the connection to it
4. We fixed the DB schematic and we have ideas for future expansion


## Day 3
### To Do
1. <del>To have a daily meeting <del>
2. <del> To fix the 3D printed modules for the sensors because they are too big (with sandpaper)<del>
3. <del>To test the scale with the 3D printed modules<del>
4. <del>To cut a wooden board for the scale body<del>
5. <del>To desing the fronend it bubble io<del>
6. <del>To make basic endpoints in bubble io <del>
7. <del>To connect the esp32 to the backend <del>

### What we did
1. We had a daily meet on which discussed the plan for the day.
2. The 3D components were a little bigger so we made them smaller with a sandpaper.
3. The scale had a little off-set because it wasn't levered properly but when we made it with a wooden board and the sensors with the 3D printed parts it was much better.
4. We made some basic desinges on how the frontend would look like
5. We made the backend endpoints so we can add values to the DB
6. We made a basic code that can access the bubble io's DB. 
7. Made some features with the bot - it can now recieve orders and add them to the db


## Day 4
### To Do
1. <del> To test if everything is correct with the sensors <del>
2. <del> To test and finished the communication with the ESP32 and the backend <del>
3. <del> To display the values to the display and the frontend <del> 
4. <del>To fix the frontend - to display orders history and orders for the home page <del>

### What we did
1. We had a daily meet on which discussed the plan for the day.
2. We finished the code for the ESP32 - now it can display the code to the display, it can measure with precision to 1% what is the change of the kilograms and to send this to the backend.
3. We can now display the on the frontend the data.
4. We made some pages including - homepage (which will hold the data for the current cooking session), hisotry gas (which will display the usage of the gas and time time when the grill was working), history order (which will display the last 4 orders and who ordered them) and single grill history.
5. We start developing the bot functionalities for starting a new order and then finishing it.
6. We changed ![the struct of the DB](https://github.com/Ne-Se-Chete/AppolicaIntern2024/blob/main/images/Bubble_io_DB.png) again because we needed more things to add.


## Day 5
### To Do
1. <del> To discuss further what are the plans for today <del>
2. <del> To finish the functionalities of the bot <del>
3. <del> To make cool desing on the frontend<del> 
4. <del> To add some features to the bot (to be more useful)<del> 
5. <del> To test the whole project<del> 

### What we did
1. We had a daily meet on which discussed the plan for the day.
2. We added the functionalities we started developing for the bot
3. We started making the desing
4. We made the desing to be cool
5. We added many commands for the bot that we will describe in the slack_bot branch
6. We tested the whole thing and it works!!


## Day 6
### To Do
1. <del> To discuss further what are the plans for today <del>
2. <del> To make cool fonts on the frontend. <del>
3. <del> To make the frontend responsive and to be good in mobile. <del>
4. <del> To test the deep sleep mode of the esp32.<del>
5. <del> If it is possible and practicle to add the deep sleep. <del>

### What we did
1. We had a daily meet on which discussed the plan for the day.
2. We made the frontend responsive and to look good.
3. We made some features in the frontend that weren't added.
4. We tested the deep sleep mode and we made it so it works with deep sleep.
5. We added receipt recongition with the Slack Bot.


## To be cheked
- ngrok.com, https://www.qovery.com/
- да се направят папки optionally за дните
- ДА има в main page-а multiple choice dropdown Който да ни казва кой е активния и така да направим в базата данни да се запазват и после да се пращат
- как да се направи bubble io-то за изчисляване на данни - дали да се записва в база данни или някак по друг начин
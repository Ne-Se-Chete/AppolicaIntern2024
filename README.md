# AppolicaIntern2024

### Project Description

This repository is for a project aimed at enhancing the management and tracking of BBQ process. The application will feature comprehensive notifications, historical data tracking, and estimation tools to optimize the cooking and consumption experience. It will include integrations with Slack, a webpage and Slack bot for order management, and additional predictive and recommendation features. The goal is to provide the user with how many kgs are left in the gas tank and when it must be refiled.

### Features of the project

1. [ ] Measure the weight of the bottle
2. [ ] Create an estimation of how long the bottle will last based on:
   - [ ] Which ones have been rotated
   - [ ] How many rotations
3. [ ] Track the history of the cooking
4. [ ] Set up notifications

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

3. For the software Bobi and Moni will be responsible and for the hardware - Yaskata and Valeri

4. We will be using:
- Bubble io for the frontend, backend and the DB
- Go(lang) for the Bot
- ESP32 for communication with the smart scale and the backend
- Load cell three wire sensors for the measurment
5. Our roadmap is ![here](https://github.com/Ne-Se-Chete/AppolicaIntern2024/blob/main/images/roadmap.png)

6. We have made a daily plan which includes:
- daily meetings at 10:00
- following the plan
- helping each other when we have problems
- every night write a short description of what we did


## Day 2
### To Do
1. <del> To have a daily meeting<del>
2. <del>To setup the scale<del>
3. <del>To do testing with the Bubble io<del>
4. <del>To fix the DB schematic<del>
5. <del>To discuss further the idea for the software<del>

### What we did
1. We had daily meeting in which we discussed the plan for the day and we made some changes
2. We found a bug with the sensors which we fixed
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
1. We had a daily meeting to discuss the plan for the day.
2. The 3D components were a little bigger so we made them smaller with sandpaper.
3. The scale was slightly off-set because it wasn't leveled properly but when we made it with a wooden board and the sensors with the 3D printed parts it was much better.
4. We made some basic designs on how the frontend would look like
5. We made the backend endpoints so we can add values to the DB
6. We made a basic code that can access the bubble io's DB. 
7. Made some features with the bot - it can now receive orders and add them to the db


## Day 4
### To Do
1. <del> To test if everything is correct with the sensors <del>
2. <del> To test and finished the communication with the ESP32 and the backend <del>
3. <del> To display the values to the display and the frontend <del> 
4. <del>To fix the frontend - to display orders history and orders for the home page <del>

### What we did
1. We had a daily meeting on which we discussed the plan for the day.
2. We finished the code for the ESP32 - now it can display the code to the display, it can measure with precision to 1% what is the change of the kilograms and send this to the backend.
3. We can now display the on the frontend the data.
4. We made some pages including - homepage (which will hold the data for the current cooking session), history gas (which will display the usage of the gas and time time when the grill was working), history order (which will display the last 4 orders and who ordered them) and single grill history.
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
1. We had a daily meeting to discuss the plan for the day.
2. We added the functionalities we started developing for the bot
3. We started making the design
4. We made the design to be cool
5. We added many commands for the bot that we will describe in the slack_bot branch
6. We tested the whole thing and it works!!


## Day 6
### To Do
1. <del> To discuss further what are the plans for today <del>
2. <del> To make cool fonts on the frontend. <del>
3. <del> To make the frontend responsive and to be good in mobile. <del>
4. <del> To test the deep sleep mode of the esp32.<del>
5. <del> If it is possible and practical to add the deep sleep. <del>

### What we did
1. We had a daily meeting on which we discussed the plan for the day.
2. We made the frontend responsive and look good.
3. We made some features in the frontend that weren't added.
4. We tested the deep sleep mode and we made it so it works with deep sleep.
5. We added receipt recognition with the Slack Bot.


## Day 7
### To Do
1. <del> To discuss further what are the plans for today. <del>
2. <del> To continue the work on the frontend. <del> 
3. <del> To test the esp32 with batteries <del>
4. <del> To find better batteries and to buy a voltage regulator.<del> 

### What we did
1. We had a daily meeting to discuss the plan for the day.
2. We fixed some bugs on the frontend.
3. We tested the batteries and unfortunately they are too old and don't function as intended.


## Day 8
### To Do
1. <del>To discuss further what are the plans for today.<del>
2. <del>To continue the work on the frontend. <del>
3. <del>To find replacement of the batteries. <del>
4. <del>To refactor the bot code because we changed the DB diagram. <del>
5. <del> To update the diagram. <del>

### What we did
1. We had a daily meeting to discuss the plan for the day.
2. We continued the tasks for the frontend.
3. We discussed some ideas with our mentors and we came to the conclusion that we will use the batteries and nothing else. The options were:
- Power bank - it "goes to sleep" when there is really low consumption and the esp32 can't be powered up when it is in deep sleep
- Charger - there is no possibility to put a power strip where the charger will be plugged
4. We changed the DB diagram and the new diagram is in the `images` folder. Now it is cleaner and has some added tables for our convenience.
5. We updated the bot so it works with the new DB.
6. We deprecated the idea of having receipt recognition because there is no way to predict how much each individual order will cost.

## Day 9
### To Do
1. <del>To discuss further what are the plans for today.<del>
2. <del>To continue the work on the frontend. <del>
3. <del>To upload the updated ESP32 code so it gets the gas weight from the db.<del>
4. <del> To set up the batteries.<del>
5. <del> To merge the branches and resolve merge conflicts.<del>

### What we did
1. We had a daily meeting to discuss the plan for the day.
2. We continued the task for the frontend and we finished the important tasks.
3. We merged the branches.
4. We fixed some issues with the hardware.
5. We made a login system and tested it for bugs.
6. We tested the bot and the esp32 code.
7. We are making the presentation for tommorrow.

## Day 10
### To do
1. <del>To test the whole project and to make a demo <del>

### What we did
1. We tested everything, it worked, we made demo and we ate :)

## To be cheked
- ngrok.com, https://www.qovery.com/

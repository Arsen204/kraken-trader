# Tinkoff fintech autumn 2021 "Golang developer" coursework.
# Description

----
This robot allows you to make automatic buying and selling decisions based on the Stochastic indicator on Kraken stock exchange.
The task description you can see in file /docs/task.md.

# Preparation

----
## Install

    git clone https://github.com/Arsen204/kraken-trader

## Fill config.env file

      API_KEY=
      API_SECRET=
      BOT_TOKEN=
      CHAT_ID=
      DB_PASSWORD=


## Up database

    docker-compose up

## Run the app

    make run

----
# Rest API

----
Below you can read the descriptions of the endpoints calls

----

**Run**
----
This endpoint allows you to run the robot. You should provide robot with necessary parameters using url. 

* **URL**

  /run

* **Method:**

  `GET`

*  **URL Params**

**Required:**
  ```
  productID
  period
  size
  limitCoef
  ```

* **Data Params**

  None

* **Success Response:**

  If successful, then you should receive only status code.

   * **Code:** `200 OK`

* **Error Response:**

  In case of failure, you should receive status code.

   * **Code:** `400 BAD REQUEST`

* **Sample Call:**

  ```
  curl "http://localhost:5000/run?productID=PI_XBTUSD&period=candles_trade_1m&size=1&limitCoef=0.1"
  ```
  or

   ```
   make req
   ```

  ----

**Stop**
----
This endpoint allows you to stop the robot.

* **URL**

  /stop

* **Method:**

  `GET`

*  **URL Params**

   None

* **Data Params**

  None

* **Success Response:**

  If successful, then you should receive only status code.

   * **Code:** `200 OK`

* **Error Response:**

  In case of failure, you should receive status code.

   * **Code:** `400 BAD REQUEST`

* **Sample Call:**

  ```
  curl "http://localhost:5000/stop"
  ```
----

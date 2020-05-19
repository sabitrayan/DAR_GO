# Code structure

We used RabbitMQ server to make message exchange

##### Command sending design

![picture](https://github.com/DairovOlzhas/dar-project/raw/master/media/commandsexchange.png)


##### Online design

![picture](https://github.com/DairovOlzhas/dar-project/raw/master/media/onlineexchange.png)

Every 2-3 seconds each player sends message to others that he is an online.
If player exit game then others detect that he didn't send message long time and delete him from game. 
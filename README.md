# Planning poker server application

## API

create a room  
`POST /v1/rooms`  
<- `{ "room_id": "" }`

create a player (the first player will be room owner)  
`POST /v1/rooms/<room_id>/players`  
-> `{ "name": "" }`  
<- `{ "access_token": "" }`

get the room state  
_authorized_  
`GET /v1/rooms/<room_id>`

delete the room  
_authorized (owner)_  
`DELETE /v1/rooms/<room_id>`

create a next game  
_authorized (owner)_  
`POST /v1/rooms/<room_id>/games`  
-> `{ "name": ""}`

update the current game (change name, complete, reset)  
_authorized (owner)_  
`PATCH /v1/rooms/<room_id>/currentgame`  
-> `{ "name": "", "complete": false, "reset": false }`

post the current game card  
_authorized_  
`PUT /v1/rooms/<room_id>/currentgame/cards`  
-> `{ "score": 0 }`

## Client flow

### Room owner flow
1. create a room
1. add a player to the room and save , the first player will be room owner
1. observe the room state
1. post cards
1. switch games
1. update the current game

### Room player flow
1. add player to the existing room
1. observe the room state
1. post cards

## Launch

The simpliest way to run and test application is using docker and sample scripts for API

Usage
``` bash
docker build -t planning-poker-app .
docker run -dp 3000:8080 planning-poker-app
cd ./scripts
export POKER_HOST=localhost:3000
create_room.sh
export POKER_ROOM="<room_id>"
create_player.sh
export POKER_SESSION="<access_token>"
get_room.sh
...
```

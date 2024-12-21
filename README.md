# Callgo

## Endpoints
POST /video/{sessionID}/{memberID} - submit videodata as json object with "video" key for a specific member \
GET /video/{sessionID}/{memberID} - get a specific member's video data

POST /session - Creates an initial session and returns a host id (which also acts as meeting id). Expects a host display name json object with the "name" key \
POST /session/{id} - Adds a member to the meeting and returns its id. Also expects display name \
GET /session/{id} - Gets all the members in a session 

## Links
The client is [here](https://github.com/HoriaBosoanca/callgo-client).
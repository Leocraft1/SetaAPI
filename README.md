# Seta-API Go Edition

This api interfaces with SETA's servers and gives back correct information to implement on a site.

It operates at:

- `/arrivals/{id}` where id is the id of the stop you want. [COMPLETE, missing news and aep support]
- `/busesinservice` list of all buses operating. [COMPLETE, missing news, aep support and periodic run to save stops (missing INSERTs)]
- `/vehicleinfo/:id` where id is the id of the vehicle you want the informations of, it needs to be operating. [COMPLETE]
- `/lineslist` list of all lines. [COMPLETE]
- `/modelslist` list of all bus models. [COMPLETE]
- `/stoplist` list of bus stops. [COMPLETE]
- `/routecodes` list of route codes.
- `/routestops/:id` where id is the route code you want to obtain the stops of.
- `/nextstops/:id` where id is the journey code of the shift you want to obtain the remaining stops of. [COMPLETE]
- `/allnews` to get all the news. [COMPLETE]
- `/news?link=[news link]` to get the content of the selected news. [COMPLETE]
- `/routeproblems` list of all route problems.
- `/routeproblems/:id` where id is the num of the route you want to know the problems of.
- `/timetable` **TODO** needs to be completely redone
- `/routemap/:id` where id is the id of the route (only the last part) you want the map of.

## Setup

To run it you will need to have the Go compiler installed in your system (installing from apt is fine).

- Clone the repo with `git clone <repo-url>`.
- Open a terminal in your folder containing the cloned repo.
- Type `go mod tidy && go build` to install the dependencies or update them and compile the executable.
- Run it with `./setaapi`.
- It might be necessary to do a `sudo chmod -X ./setaapi`.

As a development state only, it starts on port 5001 but you will be able to change this later.

## Credits

Scraping endpoints kindly done by @Curry141

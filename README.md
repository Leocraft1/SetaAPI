# Seta-API Go Edition
This api interfaces with SETA's servers and gives back correct information to implement on a site.

It operates at :

- /arrivals/{id} where id is the id of the stop you want. [COMPLETE, missing news and aep support]
- /busesinservice for the list of all buses operating. [COMPLETE, missing news, aep support and periodic run to save stops (missing INSERTs)]
- /vehicleinfo/:id where id is the id of the vehicle you want the informations of, it needs to be operating. [COMPLETE]
- /lineslist for the list of all route numbers (not static, will update when new routes are operating). [COMPLETE]
- /modelslist for the list of all bus models. [COMPLETE]
- /stoplist the result of a fetch to obtain the bus stops. [COMPLETE]
- /routestops/:id where id is the route code of the stop you want to obtain the stops of.
- /nextstops/:id where id is the journey code of the shift you want to obtain the remaining stops of. [COMPLETE]
- /allnews to get all the news. [COMPLETE]
- /news?link=[news link] to get the content of the selected news. [COMPLETE]
- /routeproblems to fetch route problems.
- /routeproblems/:id where id is the num of the route you want to know the news of.
- /shitcodes to get the damn horrible hidden codes seta uses to identify routes in his website.
- /routemap/:id where id is the id of the route (only the last part) you want the map of.
- /aepnums to get the bus numbers where it is installed the new AEP AVM

# Setup
To run it you will need to have the Go compiler installed in your system (installing from apt is fine).

- Clone the repo with `git clone <repo-url>`.
- Open a terminal in your folder containing the cloned repo.
- Type `go mod tidy && go build` to install the dependencies or update them and compile the executable.
- Run it with `./(TODO exec name here)`.
- It might be necessary to do a `sudo chmod -X ./(same thing as before)`.

As a development state only, it starts on port 5001 but you will be able to change this later.
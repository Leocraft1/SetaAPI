# Seta-API
This api interfaces with SETA's servers and gives back correct information to implement on a site.

It operates at :

- /stoplist for the list of bus stops.
- /arrivals/:id where id is the id of the stop you want.
- /busesinservice for the list of all buses operating.
- /vehicleinfo/:id where id is the id of the vehicle you want the informations of, it needs to be operating.
- /routenumberslist for the list of all route numbers (not static, will update when new routes are operating).
- /busmodels for the list of all bus models.
- /stopcodesarchive the result of the fetch to obtain the bus stops.
- /routestops/:id where id is the route code of the stop you want to obtain the stops of.
- /nextstops/:id where id is the journey code of the shift you want to obtain the remaining stops of.
- /routeproblems to fetch route problems.
- /routeproblems/:id where id is the num of the route you want to know the news of.
- /shitcodes to get the damn horrible hidden codes seta uses to identify routes in his website.

# Setup
To run it you will need Node.js (at least v22.17.1 LTS).

- Open a terminal in your folder containing the `server.js` and the `package.json` files,
- Type `npm install` to install the dependencies used in the project.

As default it starts on port 5001 but you can customize the port varying the "port" constant in the server.js file.
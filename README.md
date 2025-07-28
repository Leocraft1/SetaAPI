# Seta-API
This api interfaces with SETA's servers and gives back correct information to implement on a site.

It operates at /stoplist for the list of bus stops and at /arrivals/:id where id is the id of the stop you want.

# Setup
To run it you will need Node.js (at least v22.17.1 LTS).

- Open a terminal in your folder containing the `server.js` and the `package.json` files,
- Type `npm install` to install the dependencies used in the project.

As default it starts on port 5001 but you can customize the port varying the "port" constant in the server.js file.
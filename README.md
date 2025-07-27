# Seta-API
This api interfaces with SETA's servers and gives back correct information to implement on a site.

To run it you will need Node.js
As default it starts on port 5001 but you can customize the port varying the "port" constant in server.js

It operates at /stoplist for the list of bus stops and at /arrivals/:id where id is the id of the stop you want.

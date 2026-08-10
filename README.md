# Seta-API Go Edition

This api interfaces with SETA's servers and corrects information via database to implement on a site.

It operates at:

- `/arrivals/{id}` where id is the id of the stop you want. [COMPLETE, missing news and aep support]
- `/busesinservice` list of all buses operating. [COMPLETE, missing news, aep support and periodic run to save stops (missing INSERTs)]
- `/vehicleinfo/{id}` where id is the id of the vehicle you want the informations of, it needs to be operating. [COMPLETE]
- `/lineslist` list of all lines. [COMPLETE]
- `/modelslist` list of all bus models. [COMPLETE]
- `/stoplist` list of bus stops. [COMPLETE]
- `/routecodes` list of route codes, descriptions, display lines and destination grouped by line number. [COMPLETE]
- `/routestops/{id}` where id is the route code you want to obtain the stops of.
- `/nextstops/{id}` where id is the journey code of the shift you want to obtain the remaining stops of. [COMPLETE]
- `/allnews` to get all the news. [COMPLETE]
- `/news?link=[news link]` to get the content of the selected news. [COMPLETE]
- `/lineproblems` list of all route problems. [COMPLETE]
- `/lineproblems/{id}` where id is the num of the route you want to know the problems of. [COMPLETE]
- `/timetable` **TODO** needs to be completely redone
- `/routemap/{id}` where id is the id of the route (only the last part) you want the map of.

## Setup

### Program configuration

To run it you will need to have the Go compiler installed in your system (installing from apt is fine).

- Clone the repo with `git clone <repo-url>`.
- Open a terminal in your folder containing the cloned repo.
- Type `go mod tidy && go build` to install the dependencies or update them and compile the executable.
- Run it with `./setaapi`.
- It might be necessary to do a `sudo chmod -X ./setaapi`.

As a development state only, it starts on port 5001 but you will be able to change this later.

### Database configuration

The API needs to interface with a database to fix all the information and provide registered content. It's made to be paired with a MySQL compliant database such as MariaDB or MySQL. Simply spin it up and configure it in its `.env` file **(will be available in the future)**.

#### **Data structure**

Regards data structure, you'll need two databases so called:

- `ertpl_mezzi`
- `seta_api_content`

with various tables in them.

To obtain up to date data you can reach out to <info.ertpl@protonmail.com> asking for temporary credentials and dump our database, or just download the dumps that will be updating on this repository under the `dumps` folder every now and then.

#### **User configuration**

You'll need to create a user that can access the database and you need to grant it following privileges:

- `SELECT` on `ertpl_mezzi`
- `SELECT` on `seta_api_content`

If you don't know how to create a user or run into some access denied issue here's the syntax:

```sql
CREATE USER `username`@`%` IDENTIFIED BY `password`;
```

**REMEMBER TO CHANGE username AND password WITH YOUR WANTED CREDENTIALS!**

To generate password you can also use

```bash
openssl rand -hex 16
```

and just use its output as the password.

## Credits

Scraping endpoints kindly done by [@Daniongitub](https://github.com/Daniongithub)

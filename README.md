# Seta-API Go Edition

This api interfaces with SETA's servers and corrects information via database to implement on a site.

It operates at:

- `/arrivals/{id}` where id is the id of the stop you want.
- `/busesinservice` list of all buses operating.
- `/vehicleinfo/{id}` where id is the id of the vehicle you want the informations of, it needs to be operating.
- `/lineslist` list of all lines.
- `/modelslist` list of all bus models.
- `/stoplist` list of bus stops.
- `/routecodes` list of route codes, descriptions, display lines and destination grouped by line number.
- `/routestops/{id}` where id is the route code you want to obtain the stops of.
- `/nextstops/{id}` where id is the journey code of the shift you want to obtain the remaining stops of.
- `/allnews` to get all the news.
- `/news?link=[news link]` to get the content of the selected news.
- `/lineproblems` list of all route problems.
- `/lineproblems/{id}` where id is the num of the route you want to know the problems of.
- `/timetable?line=&verse=` list of all journeys for given line and verse (can be As for going or Di for return).
- `/routemap/{id}` where id is the id of the route (only the last part) of the map you want.

[Insert tasks are being developed]

## Setup

### Program configuration

To run it you will need to have the Go compiler installed in your system (installing from apt is fine).

- Clone the repo with `git clone <repo-url>`.
- Open a terminal in your folder containing the cloned repo.
- Type `go mod tidy && go build` to install the dependencies or update them and compile the executable.
- Run it with `./setaapi`.
- It might be necessary to do a `sudo chmod +X ./setaapi`.

### Database configuration

The API needs to interface with a database to fix all the information and provide registered content. It's made to be paired with a MySQL compliant database such as MariaDB or MySQL. Simply spin it up and configure it in its `.env` file.

#### **.env configuration**

Here there are all needed env variables:

```env
PORT=":5001"
DB_HOST="yourHost"
DB_PORT=3306
DB_USER="yourUser"
DB_PASS="yourPassword"
```

`yourUser` and `yourPassword` are your configured db credentials, `yourHost` is ip addres or DNS of your db.

#### **Data structure**

Regarding data structure, you'll need two databases so called:

- `ertpl_mezzi`
- `seta_api_content`

with various tables in them.

To obtain up to date data you can reach out to <info.ertpl@protonmail.com> asking for temporary credentials and dump our database, or just download the dumps that will be updating on this repository under the `dumps` folder every now and then.

#### **User configuration**

You'll need to create a user that can access the database and you need to grant it following privileges:

- `SELECT` on `ertpl_mezzi`
- `SELECT`, `INSERT` and `UPDATE` on `seta_api_content`

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

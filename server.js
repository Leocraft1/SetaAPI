const express = require('express');
const axios = require('axios');
const cheerio = require('cheerio');
const cors = require('cors');
const fs = require('fs');
const cron = require('node-cron');

const app = express();
const port = 5001;

app.use(cors());
app.use(express.json());

//URLs declaration section
const routeNumbersUrl = "https://wimb.setaweb.it/publicmapbe/routes/getroutesinfo/MO";
const stopCodesUrl = "https://wimb.setaweb.it/publicmapbe/vehicles/map/MO";

//Intervals for updating data
const stopcodesInterval = setInterval(updateStopCodes,20000);
const routecodesInterval = setInterval(updateRouteCodes,600000);
cron.schedule("0 */8 * * *",updateRouteNumbers);

app.get('/routenumberslist', async (req, res) => {
    try {
        res.json(await updateRouteNumbers());
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: 'Failed to sync routes' });
    }
});

app.get('/busmodels', async (req, res) => {
    // Read setabus.json
    const setabuspre = await axios.get(`https://ertpl.pages.dev/scripts/setabus.json`);
    setabus = setabuspre.data;

    // Extract only the "modello" values, removing duplicates and sort alphabetically
    const modelli = Array.from(new Set(setabus.map(bus => bus.modello)))
        .filter(Boolean)
        .sort((a, b) => a.localeCompare(b, 'it'));

    // Write to busmodels.json
    fs.writeFileSync('./busmodels.json', JSON.stringify(modelli, null, 2), 'utf8');

    res.json(modelli);
});

app.get('/routestops/:id', async (req, res) => {
    const routeId = req.params.id;
    try{
        const stopsPre = await axios.get(`https://wimb.setaweb.it/publicmapbe/waypoints/getroutewaypoints/`+routeId);
        stops = stopsPre.data;
        fixStopNames(stops);

        res.json(stops);
    }catch(error){
        res.json({"error" : "Percorso non trovato"});
    }
});

app.get('/nextstops/:id', async (req, res) => {
    const journeyId = req.params.id;
    try{
        const stopsPre = await axios.get(`https://wimb.setaweb.it/publicmapbe/vehicles/getwaypointarrivals/`+journeyId);
        stops = stopsPre.data;

        res.json(stops);
    }catch(error){
        res.json({"error" : "Corsa non trovata"});
    }
});

app.get('/stoplist', async (req, res) => {
    /* Vecchio sistema, lista sbagliata
    try {
        const url = 'https://www.setaweb.it/mo/quantomanca';
        const response = await axios.get(url);
        const $ = cheerio.load(response.data);
        const results = [];
		//era l'id del select principale, va a iterare su tutti gli option
        $('#qm_palina option').each((i, el) => {
            const element = $(el);
            const fermata = element.text(); //testo del select
            const valore = element.attr('value'); //avevo nel select un value="cose", aggiustalo per come ti serve
            results.push({
                fermata,
                valore
            });
        });

        res.json(results);
    } catch (error) {
        console.error('Errore:', error.message);
        res.status(500).send('Errore nel recupero dei dati');
    }
    */
    //Nuovo sistema, lista nuova, statica perchè Seta merda
    const data = JSON.parse(fs.readFileSync('./stoplist07-2025.json', 'utf8'));
    res.json(data);
});

app.get('/arrivals/:id', async (req, res) => {
    const stopId = req.params.id;
    if(stopId=="test"){
        //Serve per test
    }else{
        try {
            const response = await axios.get(`https://avm.setaweb.it/SETA_WS/services/arrival/${stopId}`);
            //const response = await axios.get(`http://localhost:5002/SETA_WS/services/arrival/${stopId}`);
            //Varianti
            response.data.arrival.services.forEach(service => {
                //Sant'Anna (Dislessia)
                if(service.destination=="SANT  ANNA"){
                    service.destination="SANT'ANNA";
                }
                //S.Caterina (Dislessia)
                if(service.destination=="S. CATERINA"){
                    service.destination="S.CATERINA";
                }
                //D'Avia (Dislessia)
                if(service.destination=="D AVIA"){
                    service.destination="D'AVIA";
                }
                //La torre (Dislessia)
                if(service.destination=="L A TORRE"){
                    service.destination="LA TORRE";
                }
                //1A Modena Est
                if(service.service=="1"&&service.destination=="MODENA EST"){
                    service.service="1A";
                }
                //1A Polo Leonardo
                if(service.service=="1"&&service.destination=="POLO LEONARDO"){
                    service.service="1A";
                }
                //1B Ariete
                if(service.service=="1"&&service.destination=="V. ZETA - ARIETE"){
                    service.service="1B";
                    service.destination="ARIETE";
                }
                //2A San Donnino
                if(service.service=="2"&&service.destination=="SAN DONNINO"){
                    service.service="2A";
                }
                //3A Vaciglio
                if(service.service=="3"&&service.destination=="VACIGLIO MATTARELLA"){
                    service.service="3A";
                }
                //3A Portorico (Domenica)
                if(service.service=="3"&&service.destination=="PORTORICO"){
                    service.service="3A";
                }
                //3A Vaciglio (Domenica)
                if(service.service=="3"&&service.destination=="VACIGLIO"){
                    service.service="3A";
                }
                //5 Dalla Chiesa -> La Torre              
                if(service.service=="5"&&service.destination=="DALLA CHIESA"){
                    service.destination="DALLA CHIESA (LA TORRE)";
                }
                //5A Tre Olmi
                if(service.service=="5"&&service.destination=="TRE OLMI"){
                    service.service="5A";
                }
                //7A STAZIONE FS -> GOTTARDI
                if(service.service=="7A"&&service.destination=="STAZIONE FS"){
                    service.destination="GOTTARDI";
                }
                //7/ Stazione FS
                if(service.service=="7"&&service.destination=="STAZIONE FS"){
                    service.service="7/";
                }
                //9A Marzaglia Nuova
                if(service.service=="9"&&service.destination=="MARZAGLIA"){
                    service.service="9A";
                    service.destination="MARZAGLIA NUOVA";
                }
                //9C Rubiera
                if(service.service=="9"&&service.destination=="RUBIERA"){
                    service.service="9C";
                }
                //9/ Stazione FS
                if(service.service=="9"&&service.destination=="STAZIONE FS"){
                    service.service="9/";
                }
                //10A La Rocca
                if(service.service=="10"&&service.destination=="LA ROCCA"){
                    service.service="10A";
                }
                //13A Carcere
                if(service.service=="13"&&service.destination=="CARCERI"){
                    service.service="13A";
                }
                //13F Variante di merda
                if(service.service=="13"&&service.codice_corsa.includes("133")){
                    service.service="13F";
                }
                //14A Nazioni
                if(service.service=="14"&&service.destination=="NAZIONI"){
                    service.service="14A";
                }
                //15/ Santi
                if(service.service=="15"&&service.destination=="SANTI"){
                    service.service="15/";
                }
            });
            // Step 1: Mappa i servizi per codice_corsa divisi per tipo
            const plannedMap = new Map();
            const realtimeMap = new Map();

            response.data.arrival.services.forEach(service => {
            if (service.type == "planned") {
                plannedMap.set(service.codice_corsa, service);
            } else if (service.type == "realtime") {
                realtimeMap.set(service.codice_corsa, service);
            }
            });

            // Step 2: Filtra i servizi
            const filteredServices = response.data.arrival.services.filter(service => {
            if (service.type == "realtime") return true;
            return !realtimeMap.has(service.codice_corsa);
            });

            // Step 3: Aggiungi "delay" dove possibile
            filteredServices.forEach(service => {
            if (service.type == "realtime") {
                const planned = plannedMap.get(service.codice_corsa);
                if (planned) {
                const delayMinutes = computeDelay(planned.arrival, service.arrival);
                service.delay = delayMinutes;
                }
            }
            });

            // Funzione di utilità per calcolare il ritardo in minuti
            function computeDelay(plannedTime, realtimeTime) {
                const [pHour, pMin] = plannedTime.split(":").map(Number);
                const [rHour, rMin] = realtimeTime.split(":").map(Number);
                return (rHour * 60 + rMin) - (pHour * 60 + pMin);
            }

            // Aggiorna i dati
            response.data.arrival.services = filteredServices;
            
            res.json(response.data);
        } catch (error) {
            console.error(error);
            res.json({
                "arrival" : {
                    "error" : "no arrivals scheduled in next 90 minutes",
                }
            });
        }
    }
});

//Passtrough + correggi businservizio
app.get('/busesinservice', async (req, res) => {
    const stopId = req.params.id;
    if(stopId=="test"){
        //Serve per test
    }else{
        try {
            const response = await axios.get(`https://wimb.setaweb.it/publicmapbe/vehicles/map/MO`);
            //const response = await axios.get(`https://wimb.setaweb.it/publicmapbe/vehicles/map/MO`);
            //Varianti
            fixBusRouteAndNameWimb(response);
            // Sort features by numeric part of linea
            response.data.features.sort((a, b) => {
                // Extract numeric part from linea (e.g., "13F" -> 13)
                const getNum = linea => parseInt((linea || '').match(/\d+/)?.[0] || 0, 10);
                const numA = getNum(a.properties.linea);
                const numB = getNum(b.properties.linea);
                if (numA !== numB) return numA - numB;
                // If numbers are equal, sort alphabetically
                return (a.properties.linea || '').localeCompare(b.properties.linea || '', 'it');
            });
            // Move features where linea starts with a letter to the bottom of the array
            const features = response.data.features;
            const numericLineas = [];
            const letterLineas = [];
            features.forEach(bus => {
                const linea = bus.properties.linea || '';
                if (/^[A-Za-z]/.test(linea)) {
                    letterLineas.push(bus);
                } else {
                    numericLineas.push(bus);
                }
            });
            response.data.features = numericLineas.concat(letterLineas);
            res.json(response.data);
        } catch (error) {
            console.error(error);
            res.json({
                "arrival" : {
                    "error" : "Error",
                }
            });
        }
    }
});

app.get("/vehicleinfo/:id", async (req, res) => {
    const id = req.params.id;
    try {
        const response = await axios.get(`https://wimb.setaweb.it/publicmapbe/vehicles/map/vehicle/tracking/${id}`);
        //const response = await axios.get(`https://wimb.setaweb.it/publicmapbe/vehicles/map/vehicle/tracking/${id}`);
        fixBusRouteAndNameWimb(response);
        fixPlate(response);
        fixServiceTag(response);
        await addPostiTotali(response,id);
        fixPedana(response);
        await addNextStop(response,id);
        res.json(response.data);
    } catch (error) {
        console.error(error);
        res.json({
            "arrival" : {
                "error" : "Error, or vehicle not operating",
            }
        });
    }
});

//Machine learning codici fermate
app.get('/stopcodesarchive', async (req, res) => {
    res.json(await updateStopCodes());
});

//Machine learning codici percorsi
app.get('/routecodesarchive', async (req, res) => {
    res.json(await updateRouteCodes());
});

async function updateRouteCodes() {
    const response = await axios.get(stopCodesUrl);
    const remoteRC = response.data.features;

    // Percorso file archivio locale
    const filePath = './rc_new.json';
    let localRC = [];
    if (fs.existsSync(filePath)) {
        localRC = JSON.parse(fs.readFileSync(filePath, 'utf8'));
    }

    // Trasforma l'archivio locale in una mappa per linea
    const localMap = new Map();
    localRC.forEach(item => {
        localMap.set(item.linea, new Set(item.codes));
    });

    // Costruisci una mappa temporanea per raccogliere i codici remoti per ogni linea
    const tempMap = new Map();
    remoteRC.forEach(route => {
        const linea = route.properties.linea;
        const code = route.properties.route_code;
        if (!tempMap.has(linea)) tempMap.set(linea, new Set());
        tempMap.get(linea).add(code);
    });

    // Aggiorna la mappa locale con i nuovi codici
    let changed = false;
    tempMap.forEach((codes, linea) => {
        if (!localMap.has(linea)) {
            localMap.set(linea, new Set(codes));
            changed = true;
        } else {
            codes.forEach(code => {
                if (!localMap.get(linea).has(code)) {
                    localMap.get(linea).add(code);
                    changed = true;
                }
            });
        }
    });

    // Ricostruisci l'array ordinato per linea
    const updatedRC = Array.from(localMap.entries())
        .map(([linea, codes]) => ({
            linea,
            codes: Array.from(codes).sort((a, b) => a.localeCompare(b, 'it'))
        }))
        .sort((a, b) => a.linea.localeCompare(b.linea, 'it'));
    // Trasforma di nuovo la mappa in array ordinato per linea
    const newLocalRC = Array.from(localMap.entries())
        .map(([linea, codes]) => ({
            linea,
            codes: Array.from(codes)
        }))
        .sort((a, b) => Number(a.linea) - Number(b.linea));
    // Salva solo se ci sono cambiamenti
    if (changed) {
        fs.writeFileSync(filePath, JSON.stringify(newLocalRC, null, 2), 'utf8');
    }
    console.log("["+new Date()+"] Journey codes updated.");
    return newLocalRC;
}

async function updateStopCodes(){
    //console.log("["+new Date()+"] Updating stop codes.");
    // 1. Fetch routes from the given URL
    const response = await axios.get(stopCodesUrl);
    const remoteStops = response.data.features;

    // 2. Load local stops file
    const filePath = './stoplist_new.json';
    let localStops = [];
    if (fs.existsSync(filePath)) {
        localStops = JSON.parse(fs.readFileSync(filePath, 'utf8'));
    }

    // 3. Build a Set of existing stop codes
    const localStopCodes = new Set(localStops.map(stop => stop.valore));

    // 4. Add new stops if not present
    let added = false;
    remoteStops.forEach(route => {
        var fermata = route.properties.wp_desc;
        const valore = route.properties.reached_waypoint_code;
        //Nomi stazione autostazione e garibaldi
        /*
        if (valore == "XYZ") {
            fermata = "Nome Personalizzato";
        }
        */
        if (valore == "MO6132") {
            fermata = "STAZIONE FS (Corsia 1)";
        }
        if (valore == "MO6133") {
            fermata = "STAZIONE FS (Corsia 2)";
        }
        if (valore == "MO6134") {
            fermata = "STAZIONE FS (Corsia 3)";
        }
        if (valore == "MO6119") {
            fermata = "STAZIONE FS (Corsia 4)";
        }
        if (valore == "MO6121") {
            fermata = "MODENA AUTOSTAZIONE (dir. Centro)";
        }
        if (valore == "MO5003") {
            fermata = "MODENA AUTOSTAZIONE (lato Novi Park)";
        }
        if (valore == "MO6600") {
            fermata = "MODENA AUTOSTAZIONE (davanti biglietteria)";
        }
        if (valore == "MO10") {
            fermata = "MODENA AUTOSTAZIONE (fianco biglietteria)";
        }
        if (valore == "MO6120") {
            fermata = "MODENA AUTOSTAZIONE (fianco biglietteria lato Novi Park)";
        }
        if (valore == "MO3") {
            fermata = "MODENA AUTOSTAZIONE (Corriere corsia 1)";
        }
        if (valore == "MO303") {
            fermata = "MODENA AUTOSTAZIONE (Corriere corsia 2)";
        }
        if (valore == "MO342") {
            fermata = "MODENA AUTOSTAZIONE (Corriere corsia 3)";
        }
        if (valore == "MO344") {
            fermata = "MODENA AUTOSTAZIONE (Corriere corsia 4)";
        }
        if (valore == "MO350") {
            fermata = "MODENA AUTOSTAZIONE (Corriere corsia 5)";
        }
        if (valore == "MO346") {
            fermata = "MODENA AUTOSTAZIONE (Corriere corsia 6)";
        }
        if (valore == "MO5900") {
            fermata = "GARIBALDI (dir. Centro)";
        }
        if (valore == "MO30") {
            fermata = "GARIBALDI (dir. Trento Trieste)";
        }
        if (valore == "MO9") {
            fermata = "GARIBALDI (lato Caduti in Guerra)";
        }
        if (valore == "MO5111") {
            fermata = "GARIBALDI (Storchi dir. Trento Trieste)";
        }
        if (valore == "MO5112") {
            fermata = "GARIBALDI (Storchi dir. Centro)";
        }
        if (valore == "MOPALTEC3") {
            fermata = "NAZIONI CAPOLINEA";
        }
        if (valore == "L A TORRE") {
            fermata = "LA TORRE";
        }
        if (!localStopCodes.has(valore)&&!valore=='') {
            localStops.push({ fermata, valore });
            localStopCodes.add(valore);
            added = true;
        }
    });

    // 5. Save updated file if changed
    if (added) {
        // Ordina localStops per "valore" (alfanumerico)
        localStops.sort((a, b) => (a.valore || '').localeCompare(b.valore || '', 'it'));
        fs.writeFileSync(filePath, JSON.stringify(localStops, null, 2), 'utf8');
    }

    const data = JSON.parse(fs.readFileSync('./stoplist_new.json', 'utf8'));
    console.log("["+new Date()+"] Stop codes updated.");
    return(data);
}

async function updateRouteNumbers(){
    //console.log("["+new Date()+"] Updating route numbers.");
    // 1. Fetch routes from the given URL
    const response = await axios.get(routeNumbersUrl);
    const remoteRoutes = response.data.routesdata;

    // 2. Load local routes file
    const filePath = './routelist.json';
    let localRoutes = [];
    if (fs.existsSync(filePath)) {
        localRoutes = JSON.parse(fs.readFileSync(filePath, 'utf8'));
    }

    // 3. Build a Set of existing routeCodes
    const localRouteCodes = new Set(localRoutes);

    // 4. Add new routes if not present
    let added = false;
    remoteRoutes.forEach(route => {
        if (!localRouteCodes.has(route.linea)) {
            localRoutes.push(route.linea);
            localRouteCodes.add(route.linea);
            added = true;
        }
    });

    // 5. Save updated file if changed
    if (added) {
        fs.writeFileSync(filePath, JSON.stringify(localRoutes, null, 2), 'utf8');
    }

    const data = JSON.parse(fs.readFileSync('./routelist.json', 'utf8'));
    console.log("["+new Date()+"] Route numbers updated.");
    return(data);
}

function fixBusRouteAndNameWimb(response){
    response.data.features.forEach(bus => {
        service = bus.properties;
        //Sant'Anna (Dislessia)
        if(service.route_desc=="SANT  ANNA"){
            service.route_desc="SANT'ANNA";
        }
        //S.Caterina (Dislessia)
        if(service.route_desc=="S. CATERINA"){
            service.route_desc="S.CATERINA";
        }
        //D'Avia (Dislessia)
        if(service.route_desc=="D AVIA"){
            service.route_desc="D'AVIA";
        }
        //1A Modena Est
        if(service.linea=="1"&&service.route_desc=="MODENA EST"){
            service.linea="1A";
        }
        //1A Polo Leonardo
        if(service.linea=="1"&&service.route_desc=="POLO LEONARDO"){
            service.linea="1A";
        }
        //1B Ariete
        if(service.linea=="1"&&service.route_desc=="V. ZETA - ARIETE"){
            service.linea="1B";
            service.route_desc="ARIETE";
        }
        //2A San Donnino
        if(service.linea=="2"&&service.route_desc=="SAN DONNINO"){
            service.linea="2A";
        }
        //3A Vaciglio
        if(service.linea=="3"&&service.route_desc=="VACIGLIO MATTARELLA"){
            service.linea="3A";
        }
        //3A Portorico (Domenica)
        if(service.linea=="3"&&service.route_desc=="PORTORICO"){
            service.linea="3A";
        }
        //3A Vaciglio (Domenica)
        if(service.linea=="3"&&service.route_desc=="VACIGLIO"){
            service.linea="3A";
        }
        //5 Dalla Chiesa -> La Torre
        if(service.linea=="5"&&service.route_desc=="DALLA CHIESA"){
            service.route_desc="DALLA CHIESA (LA TORRE)";
        }
        //5A Tre Olmi
        if(service.linea=="5"&&service.route_desc=="TRE OLMI"){
            service.linea="5A";
        }
        //7A STAZIONE FS -> GOTTARDI
        if(service.linea=="7A"&&service.route_desc=="STAZIONE FS"){
            service.route_desc="GOTTARDI";
        }
        //7/ Stazione FS
        if(service.linea=="7"&&service.route_desc=="STAZIONE FS"){
            service.linea="7/";
        }
        //9A Marzaglia Nuova
        if(service.linea=="9"&&service.route_desc=="MARZAGLIA"){
            service.linea="9A";
            service.route_desc="MARZAGLIA NUOVA";
        }
        //9C Rubiera
        if(service.linea=="9"&&service.route_desc=="RUBIERA"){
            service.linea="9C";
        }
        //9/ Stazione FS
        if(service.linea=="9"&&service.route_desc=="STAZIONE FS"){
            service.linea="9/";
        }
        //10A La Rocca
        if(service.linea=="10"&&service.route_desc=="LA ROCCA"){
            service.linea="10A";
        }
        //13A Carcere
        if(service.linea=="13"&&service.route_desc=="CARCERI"){
            service.linea="13A";
        }
        //13F Variante di merda
        if(service.linea=="13"&&service.route_code.includes("133")){
            service.linea="13F";
        }
        //14A Nazioni
        if(service.linea=="14"&&service.route_desc=="NAZIONI"){
            service.linea="14A";
        }
        //15/ Santi
        if(service.linea=="15"&&service.route_desc=="SANTI"){
            service.linea="15/";
        }
        //Cambia nomi modello bus
        if(service.model=="IVECO - URBANWAY MILD HYBRID CNG"){
            service.model="Iveco Urbanway Hybrid CNG";
        }
        if(service.model=="IRISBUS - PS09D2"){
            service.model="Irisbus Citelis CNG";
        }
        if(service.model=="BREDAMENARINIBUS - CITYMOOD LNG "){
            service.model="Menarinibus Citymood LNG";
        }
        if(service.model=="MENARINI BUS - M250CNG"){
            service.model="Menarinibus Citymood CNG";
        }
        if(service.model=="IVECO BUS - CBLE4/00 4C G13A4"&&service.vehicle_code>1000){
            service.model="Iveco Crossway LE CNG";
        }
        if(service.model=="IVECO - IVECOA60"){
            service.model="Iveco Daily";
        }
        if(service.model=="IVECO BUS - CBLE4/00 4C G13A4"&&service.vehicle_code<1000){
            service.model="Iveco Crossway LE Diesel";
        }
        if(service.model=="CBCW3/00 3C B1UA3 - IVECO BUS"){
            service.model="Iveco Crossway Line";
        }
        if(service.model=="IVECO BUS - CROSSWAY LINE"){
            service.model="Iveco Crossway Line";
        }
        if(service.model=="IVECO BUS - CROSSWAY 10,7 MT"){
            service.model="Iveco Crossway 10.7";
        }
        if(service.model=="IRISBUS - CROSSWAY"){
            service.model="Irisbus Crossway";
        }
        if(service.model=="MERCEDES BENZ - MERCEDES BENZ O 550U/E3-4"){
            service.model="Mercedes Integro";
        }
        if(service.model=="IRIZAR/SCANIA - I4CD2-SCN-"){
            service.model="Irizar i4 LNG";
        }
        if(service.model=="i4CD2 - IRIZAR"){
            service.model="Irizar i4 LNG";
        }
        if(service.model=="EVOBUS MB - INTEGRO/E5"){
            service.model="Mercedes Integro";
        }
        if(service.model=="SOLARIS - URBINO 12"){
            service.model="Solaris Urbino 12 CNG";
        }
        if(service.model=="MENARINI BUS - M250LNG"){
            service.model="Menarinibus Citymood LNG";
        }
        if(service.model=="VOLKSWAGEN - CRAFTER"){
            service.model="Iveco Crossway Line CNG";
        }
        if(service.model=="IVECO BUS - CBLE4/00 4C G1MA4 "){
            service.model="Iveco Crossway LE Diesel";
        }
        if(service.model=="IVECO BUS - CBLE4/00"){
            service.model="Iveco Crossway LE Diesel";
        }
        if(service.model=="IVECO - A60C17"){
            service.model="Iveco Daily";
        }
        if(service.model=="MAN - MAN R13 - EURO 6"){
            service.model="MAN Lion's Regio";
        }
        if(service.model=="SETRA - 415 NF "){
            service.model="Setra ex Bolzano";
        }
        if(service.vehicle_code>=4501&&service.vehicle_code<=4507){
            service.model="Setra ex Bolzano (2 porte)";
        }
        if(service.vehicle_code>=4508&&service.vehicle_code<=4513){
            service.model="Setra ex Bolzano (3 porte)";
        }
        if(service.vehicle_code=="4771"){
            service.model="Iveco Urbanway Hybrid CNG";
        }
        if(service.model=="IVECO BUS - CROSSWAY 12 MT"){
            service.model="Iveco Crossway Line";
        }
        if(service.model=="IVECO FRANCE - SFR 160 - CROSSWAY"){
            service.model="Irisbus Crossway Esercito";
        }
        if(service.model=="IVECO BUS - CROSSWAY LE - CNG"){
            service.model="Iveco Crossway LE CNG";
        }
        if(service.model=="SCANIA i4 LNG - IRIZAR"){
            service.model="Irizar i4 LNG";
        }
        if(service.model=="VOLVO - 8700B12B"){
            service.model="Volvo 8700 (Reggio Emilia)";
        }
        if(service.model=="EVOBUS MB - O 550 U"){
            service.model="Mercedes Integro";
        }
        if(service.model=="MERCEDES - MERCEDES SPRINTER ALTAS"){
            service.model="Sprinter";
        }
        if(service.model=="IVECO - CROSSWAY LE 14,49 MT"){
            service.model="Iveco Crossway LE 14m";
        }
        if(service.model=="IVECO - IVECOA50C17/P IRISBUS"){
            service.model="Irisbus Crossway";
        }
        if(service.model=="DUMMYBRAND - DUMMYBRAND"){
            service.model="Iveco Daily";
        }
        if(service.model=="DUMMYBRAND -DUMMYBRAND"){
            service.model="Iveco Daily";
        }
        if(service.model=="SOLARIS - URBINO U 8,9 LE"){
            service.model="Solaris 9 LE ex Piacenza";
        }
        if(service.model=="IVECO BUS - CBCW3/00 3C K13A3"){
            service.model="Iveco Crossway Line CNG";
        }
        if(service.model=="IVECO BUS - IS56AC2DA 5A1CX11"){
            service.model="Iveco Daily";
        }
        if(service.model=="IVECO BUS - CBCW3/00 3C K13A3"){
            service.model="Iveco Crossway Line CNG";
        }
        if(service.model=="MAN - LIONS CITY 19"){
            service.model="New MAN Lion's City 19G";
        }
        if(service.model=="MERCEDES BENZ - O 530/CNG/U-4"){
            service.model="Mercedes Citaro CNG";
        }
        if(service.model=="IRISBUS - CITYCLASS 491.12.27 CNG"){
            service.model="Irisbus Cityclass CNG ATCM";
        }
        //Dio filoviario arrostito sulla 750 volt
        if(service.vehicle_code>=35&&service.vehicle_code<=44){
            service.model="Solaris Trollino 12";
            service.plate_num="MO0"+service.vehicle_code;
        }
    });
}

function fixStopNames(stops){
    stops.features.forEach(element =>{
        if(element.properties.desc=="L A TORRE"){
            element.properties.desc="LA TORRE"
        }
    })
}

function fixPlate(response){
    response.data.features.forEach(element => {
        if(element.properties.plate_num.includes(" ")){
            element.properties.plate_num=element.properties.plate_num.replace(" ","")
        }
        //CW CNG senza targa
        if(element.properties.plate_num.includes("UN")){
            element.properties.plate_num=element.properties.plate_num=="Non disponibile";
        }
    });
}

//UR -> Urbano eccetera
function fixServiceTag(response){
    response.data.features.forEach(element => {
        if(element.properties.service_tag=="UR"){
            element.properties.service_tag="Urbano";
        }
        //Supposizione
        if(element.properties.service_tag=="SU"){
            element.properties.service_tag="Suburbano";
        }
        if(element.properties.service_tag=="EX"){
            element.properties.service_tag="Extraurbano";
        }
    });
}

async function getBusList(){
    const urlList = "https://setaapi.serverissimo.freeddns.org/busesinservice";
    const data = await axios.get(urlList);
    const item = data.data;
    return item;
}
//Cerca e aggiunge posti_totali alle info veicolo
async function addPostiTotali(response,idMezzo){
    var item=await getBusList();
    item.features.forEach(element =>{
        if(element.properties.vehicle_code==idMezzo){
            response.data.features.forEach(bus =>{
                bus.properties.posti_totali=element.properties.posti_totali;
            });
        }
    });
}
//Cerca e aggiunge next_stop alle info veicolo
async function addNextStop(response,idMezzo){
    var item=await getBusList();
    item.features.forEach(element =>{
        if(element.properties.vehicle_code==idMezzo){
            response.data.features.forEach(bus =>{
                bus.properties.next_stop=element.properties.next_stop;
            });
        }
    });
}

//Urbanway e Menarini LNG hanno la pedana
function fixPedana(response){
    response.data.features.forEach(element => {
        if(element.properties.model=="Iveco Urbanway Hybrid CNG"){
            element.properties.pedana=1;
        }
        if(element.properties.model=="Menarinibus Citymood LNG"){
            element.properties.pedana=1;
        }
    });
}

app.listen(port, () => {
    console.log(`API attiva su http://localhost:${port}`);
});
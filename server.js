const express = require('express');
const axios = require('axios');
const cors = require('cors');
const fs = require('fs');
const cron = require('node-cron');
const helmet = require('helmet')
const scrapeIt = require('scrape-it');

const app = express();
const port = 5001;

app.use(cors());
app.use(express.json());

//URLs declaration section
const routeNumbersUrl = "https://wimb.setaweb.it/publicmapbe/routes/getroutesinfo/MO";
const stopCodesUrl = "https://wimb.setaweb.it/publicmapbe/vehicles/map/MO";
const newsUrl = "https://www.setaweb.it/mo/news";
const lineedynUrl = "https://www.setaweb.it/mo/lineedyn";
const problemsBaseUrl = "https://www.setaweb.it/mo/news/linea/";

//Intervals for updating data
setInterval(updateStopCodes,20000);
setInterval(updateRouteCodes,600000);
cron.schedule("0 */8 * * *",updateRouteNumbers);
var i=0;
var j=0;
setInterval(updateRoutesStops,10000);

//CORS e X-Frame-Options per prevenire errori
app.use((req, res, next) => {
    res.header("Access-Control-Allow-Origin", "*"); // oppure specifica un dominio
    res.header("Access-Control-Allow-Methods", "*");
    res.header("Access-Control-Allow-Headers", "Content-Type, Authorization");
  
    // Gestione preflight OPTIONS
    if (req.method === "OPTIONS") {
        return res.sendStatus(200);
    }
    helmet({
        contentSecurityPolicy: {
            directives: {
            defaultSrc: ["'self'"],
            },
        },
        frameguard: false,
        })
    next();
});

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
    const setabuspre = await axios.get(`https://ertpl.pages.dev/seta_modena/menu/js/setabus.json`);
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
    const trimmed = routeId.split('(')[0].trim();
    res.json(await updateRouteStops(trimmed));
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
            const problems = await axios.get(`https://setaapi.serverissimo.freeddns.org/routeproblems`);
            const d = new Date();
            //Varianti
            response.data.arrival.services.forEach(service => {
                //Aggiungi problemi
                problems.data.codes.forEach(element =>{
                    if(element.num==service.service){
                        service.hasProblems=element.hasProblems;
                    }
                })
                //Aggiungi servizio senza variante (per notizie)
                service.officialService=service.service;

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
                //1S Autostazione (Scuola)
                if(service.service=="1"&&service.destination=="AUTOSTAZIONE"){
                    service.service="1S";
                }
                //1 _ -> Marinuzzi (Scuola)
                if(service.service=="1"&&service.destination=="_"){
                    service.destination="MARINUZZI";
                }
                //2A San Donnino
                if(service.service=="2"&&service.destination=="SAN DONNINO"){
                    service.service="2A";
                }
                //2/ Autostazione
                if(service.service=="2"&&service.destination=="AUTOSTAZIONE"){
                    service.service="2/";
                }
                //3A SANTA CATERINA-MONTEFIORINO (as 25/26)
                if(service.service=="3"&&service.codice_corsa.includes("339")){
                    service.service="3A";
                    service.destination="S.CATERINA-MONTEFIORINO";
                }               
                
                //3A Vaciglio
                if(service.service=="3"&&service.destination=="VACIGLIO"){
                    service.service="3A";
                }
                //3B Ragazzi 99 (as 25/26)
                if(service.service=="3"&&service.destination=="RAGAZZI DEL 99"){
                    service.service="3B";
                    service.destination="RAGAZZI 99";
                }
                //3B Nonantolana 1010 (as 25/26)
                if(service.service=="3"&&service.destination=="NONANTOLANA 1010"){
                    service.service="3B";
                }
                //3/ Stazione FS (as 25/26)
                if(service.service=="3"&&service.destination=="STAZIONE FS"){
                    service.service="3/";
                }
                //3A SANTA CATERINA-MONTEFIORINO NO MALAVOLTI (Domenica)
                if(service.service=="3"&&service.codice_corsa.includes("407")){
                    service.service="3A";
                    service.destination="S.CATERINA-MONTEFIORINO (NO MALAVOLTI)";
                }
                //3B SANTA CATERINA-MONTEFIORINO NO MALAVOLTI (Domenica)
                if(service.service=="3B"&&d.getDay()==6){
                    service.destination="NONANTOLANA 1010 (NO TORRAZZI)";
                }
                //4/ Autostazione (as 25/26)
                if(service.service=="4"&&service.destination=="AUTOSTAZIONE"){
                    service.service="4/";
                }
                //5 Dalla Chiesa -> La Torre              
                if(service.service=="5"&&service.destination=="DALLA CHIESA"){
                    service.destination="LA TORRE";
                }
                //5A Tre Olmi
                if(service.service=="5"&&service.destination=="TRE OLMI"){
                    service.service="5A";
                }
                //6A Santi (as 25/26)
                if(service.service=="6"&&service.destination=="SANTI"){
                    service.service="6A";
                }
                //6B Villanova (as 25/26)
                if(service.service=="6"&&service.destination=="VILLANOVA"){
                    service.service="6B";
                }
                //7 GOTTARDI -> POLICLINICO GOTTARDI
                if(service.service=="7"&&service.destination=="GOTTARDI"){
                    service.destination="POLICLINICO GOTTARDI";
                }
                //7A STAZIONE FS -> GOTTARDI
                if(service.service=="7A"&&service.destination=="STAZIONE FS"&&!service.codice_corsa.includes("728")){
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
                //9/ Autostazione
                if(service.service=="9"&&service.destination=="AUTOSTAZIONE"){
                    service.service="9/";
                }
                //10A La Rocca
                if(service.service=="10"&&service.destination=="LA ROCCA"){
                    service.service="10A";
                }
                //10S Liceo Sigonio
                if(service.service=="10"&&service.destination=="LICEO SIGONIO"){
                    service.service="10S";
                }
                //10/ AUTOSTAZIONE
                if(service.service=="10"&&service.destination=="AUTOSTAZIONE"){
                    service.service="10/";
                }
                //11/ Stazione FS
                if(service.service=="11"&&service.destination=="STAZIONE FS"){
                    service.service="11/";
                }
                //12A Nazioni ma sono dei coglioni di merda
                if(service.service=="12"&&(service.codice_corsa.includes("1280")||service.codice_corsa.includes("1284"))&&service.officialService!="5taxi"){
                    service.service="12A";
                    service.destination="NAZIONI";
                }
                //12/ Fanti FS
                if(service.service=="12"&&service.destination=="FANTI FS"){
                    service.service="12/";
                }
                //12S Garibaldi (Scuola)
                if(service.service=="12"&&service.destination=="GARIBALDI"){
                    service.service="12S";
                }
                //12S Largo Garibaldi (Scuola)
                if(service.service=="12"&&service.destination=="LARGO GARIBALDI"){
                    service.service="12S";
                }
                //13 S.ANNA -> SANT'ANNA (dio rincoglionito e dislessico)
                if(service.service=="13"&&service.destination=="S.ANNA"){
                    service.destination="SANT'ANNA";
                }
                //13A Carcere
                if(service.service=="13"&&service.destination=="CARCERI"){
                    service.service="13A";
                }
                //13F Variante di merda
                if(service.service=="13"&&service.codice_corsa.includes("133")){
                    service.service="13F";
                }
                //643 _ -> Polo Scolastico Sassuolo
                if(service.service=="643"&&service.destination=="_"){
                    service.destination="POLO SCOLASTICO SASSUOLO";
                }
                //Varianti vecchie AS 24/25
                /*
                //3A Vaciglio-Mattarella
                if(service.service=="3"&&service.destination=="VACIGLIO MATTARELLA"){
                    service.service="3A";
                }
                //3A Portorico (Domenica)
                if(service.service=="3"&&service.destination=="PORTORICO"){
                    service.service="3A";
                }
                //14A Nazioni
                if(service.service=="14"&&service.destination=="NAZIONI"){
                    service.service="14A";
                }
                //15/ Santi
                if(service.service=="15"&&service.destination=="SANTI"){
                    service.service="15/";
                }
                */
                //Varianti settembre 2025
                /*
                //3F MONTEFIORINO-PORTORICO (Domenica)
                if(service.service=="3"&&service.destination=="PORTORICO"){
                    service.service="3F";
                    service.destination="MONTEFIORINO-PORTORICO";
                }
                */
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
        await addNextStopCode(response,id);
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

//Forward lineedyn_linea_dett_fermate_e_orari.php
app.get('/lineedyn_linea_dett_percorsi', async (req, res) => {  
    var response = undefined;
    if(req.query.dd==undefined){
        response = await axios.get(`https://www.setaweb.it/lineedyn_linea_dett_percorsi.php?b=`+req.query.b+"&l="+req.query.l+"&v="+req.query.v);
    }else{
        response = await axios.get(`https://www.setaweb.it/lineedyn_linea_dett_percorsi.php?b=`+req.query.b+"&l="+req.query.l+"&dd="+req.query.dd+"&v="+req.query.v);
    }
    res.end(response.data);
});

//API per ottenere l'elenco notizie seta in json
app.get('/allnews', async (req, res) => {  
    const response = await axios.get(newsUrl);
    const data = scrapeIt.scrapeHTML(response.data, {
        news: {
            listItem: ".news li div div a",
            data: {
                title: ".title",
                date: ".date-title",
                link: {
                    attr: "href", convert: x => "https://www.setaweb.it/"+x,
                },
                type: {
                    selector: ".image-news", attr: "style",convert:x => {
                        if(x.includes("pericolo.png")){
                            return "Importante";
                        }if(x.includes("seta-informa.png")){
                            return "Informazione";
                        }if(x.includes("novita.png")){
                            return "Novità";
                        }if(x.includes("orari.png")){
                            return "Orari";
                        }if(x.includes("autobus-treno.png")){
                            return "Autobus Treno";
                        }if(x.includes("tessera.png")){
                            return "Biglietti";
                        }if(x.includes("lavori-in-corso.png")){
                            return "Lavori in corso";
                        }if(x.includes("controllore.png")){
                            return "Personale";
                        }
                    }
                }
            }
        }
    })
    res.json(data);
});

//API per ottenere notizia seta in json
app.get('/news', async (req, res) => {  
    const link = req.query.link;
    const response = await axios.get(link);
    const data = scrapeIt.scrapeHTML(response.data, {
        title: ".container-title",
        date: ".container-date-title",
        content: {
            selector: ".descrizione",how: "html"
        },
    })
    res.json(data);
});

//Ottenere elenco di linee con problemi
app.get('/routeproblems', async (req, res) => {  
    const response = await axios.get(lineedynUrl);
    const data = scrapeIt.scrapeHTML(response.data, {
        codes: {
            listItem: ".riga_copertura",
            data: {
                num: ".numero",
                hasProblems: {
                    selector: ".news_linea",
                    convert: x => {
                        if(x!=""){
                            return true;
                        }else{
                            return false;
                        }
                    },
                },
                shitCode: {
                    selector: ".testo_per_ricerca",
                    attr: "data-val"
                },
            }
        }
    })
    //Va bene dani
    data.codes[data.codes.length-1].num="Dani";
    data.codes[data.codes.length-1].hasProblems=true;
    data.codes[data.codes.length-1].shitCode="2F";
    res.json(data)
});

//Ottenere il dettaglio dei problemi sulla linea
app.get('/routeproblems/:id', async (req, res) => {  
    const id = req.params.id;
    const codiciDiMerda = await getCodiciDiMerda();
    var sc = "";
    codiciDiMerda.codes.forEach(element => {
        if(element.num==id){
            sc = element.shitCode;
        }
    });
    const problemsUrl = problemsBaseUrl + sc;
    try{
        const response = await axios.get(problemsUrl);
        const data = scrapeIt.scrapeHTML(response.data, {
            problems: {
                listItem: ".news-container",
                data: {
                    title: ".news-container div",
                    date: ".date-title",
                    link: {
                        closest: "a",
                        attr: "href",
                        convert: x => "https://www.setaweb.it/"+x
                    },
                },
            }
        })
        res.json(data);
    }catch(error){
        console.error(error);
        res.json({
            "error" : "Route number is not valid."
        });
    } 
});

//Codici del cazzo del loro sito di merda
app.get('/shitcodes', async (req, res) => {  
    res.json(await getCodiciDiMerda());
});

async function updateRoutesStops(){
    const routeCodes = await updateRouteCodes();
    const routeId=routeCodes[i].codes[j];
    if(routeCodes[i].codes[j+1]==undefined){
        i++;
        j=0;
    }else if(routeCodes[i+1]==undefined){
        i=0;
        j=0;
    }else{
        j++;
    }
    if(!routeId.includes("(")){
        await updateRouteStops(routeId);
    }
}

async function updateRouteStops(routeId) {
    const filePath = './routestops/'+routeId+'.json';
    try{
        const stopsPre = await axios.get(`https://wimb.setaweb.it/publicmapbe/waypoints/getroutewaypoints/`+routeId);
        const remoteStops = stopsPre.data.features;

        // 2. Load local stops file
        // Assicurati che la cartella esista
        if (!fs.existsSync('./routestops')) {
            fs.mkdirSync('./routestops', { recursive: true });
        }

        // Se il file non esiste, crealo con un array vuoto
        if (!fs.existsSync(filePath)) {
            fs.writeFileSync(filePath, '[]', 'utf8');
        }
        
        localStops = JSON.parse(fs.readFileSync(filePath, 'utf8'));
        //Dislessia
        fixStopNames(remoteStops)

        // 3. Build a Set of existing stop codes
        const localStopCodes = new Set(localStops.map(stop => stop.code));

        // 4. Add new stops if not present
        let added = false;
        remoteStops.forEach(route => {
            const desc = route.properties.desc;
            const code = route.properties.code;
            const islast = route.properties.islast;
            if (!localStopCodes.has(code)) {
                localStops.push({ desc, code, islast });
                localStopCodes.add(code);
            }
        });
        //Salva
        fs.writeFileSync(filePath, JSON.stringify(localStops, null, 2), 'utf8');
        //Compone stillExists
        const data = JSON.parse(`{
            "stillExists" : ${true},
            "stops" : ${fs.readFileSync('./routestops/'+routeId+'.json', 'utf8')}
        }`);
        
        return(data);
    }catch(error){
        if (!fs.existsSync(filePath)) {
            console.log(error);
            return({"error" : "Percorso non trovato"});
        }else{
            const data = JSON.parse(`{
                "stillExists" : ${false},
                "stops" : ${fs.readFileSync('./routestops/'+routeId+'.json', 'utf8')}
            }`);

            return(data);
        }
    }
}

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
        if (valore == "MO149") {
            fermata = "GOTTARDI (cap. 7)";
        }
        if (valore == "MO148") {
            fermata = "GOTTARDI (cap. 9)";
        }
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
        if (valore == "MO6783") {
            fermata = "POLO LEONARDO (Strada)";
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
    //console.log("["+new Date()+"] Stop codes updated.");
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
    const d = new Date();
    response.data.features.forEach(bus => {
        service = bus.properties;
        //Sant'Anna (Dislessia)
        if(service.route_desc=="SANT  ANNA"){
            service.route_desc="SANT'ANNA";
        }
        //Sant'Anna (Rincoglionimento)
        if(service.route_desc=="S.ANNA"){
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
        //La torre (Dislessia)
        if(service.route_desc=="L A TORRE"){
            service.route_desc="LA TORRE";
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
        //1S Autostazione (Scuola)
        if(service.linea=="1"&&service.route_desc=="AUTOSTAZIONE"){
            service.linea="1S";
        }
        //1 _ -> Marinuzzi (Scuola)
        if(service.linea=="1"&&service.route_desc=="_"){
            service.route_desc="MARINUZZI";
        }
        //2A San Donnino
        if(service.linea=="2"&&service.route_desc=="SAN DONNINO"){
            service.linea="2A";
        }
        //2/ Autostazione
        if(service.linea=="2"&&service.route_desc=="AUTOSTAZIONE"){
            service.linea="2/";
        }
        //3A SANTA CATERINA-MONTEFIORINO (as 25/26)
        if(service.linea=="3"&&service.route_code.includes("339")){
            service.linea="3A";
            service.route_desc="S.CATERINA-MONTEFIORINO";
        }
        //3A Vaciglio
        if(service.linea=="3"&&service.route_desc=="VACIGLIO"){
            service.linea="3A";
        }
        //3B Ragazzi del 99 (as 25/26)
        if(service.linea=="3"&&service.route_desc=="RAGAZZI DEL 99"){
            service.linea="3B";
        }
        //3B Nonantolana 1010 (as 25/26)
        if(service.linea=="3"&&service.route_desc=="NONANTOLANA 1010"){
            service.linea="3B";
        }
        //3/ Stazione FS (as 25/26)
        if(service.linea=="3"&&service.route_desc=="STAZIONE FS"){
            service.linea="3/";
        }
        //3A SANTA CATERINA-MONTEFIORINO NO MALAVOLTI (Domenica)
        if(service.linea=="3"&&service.route_code.includes("407")){
            service.linea="3A";
            service.route_desc="S.CATERINA-MONTEFIORINO (NO MALAVOLTI)";
        }
        //3B SANTA CATERINA-MONTEFIORINO NO MALAVOLTI (Domenica)
        if(service.linea=="3B"&&d.getDay()==6){
            service.route_desc="NONANTOLANA 1010 (NO TORRAZZI)";
        }
        //4/ Autostazione (as 25/26)
        if(service.linea=="4"&&service.route_desc=="AUTOSTAZIONE"){
            service.linea="4/";
        }
        //5 Dalla Chiesa -> La Torre              
        if(service.linea=="5"&&service.route_desc=="DALLA CHIESA"){
            service.route_desc="LA TORRE";
        }
        //5A Tre Olmi
        if(service.linea=="5"&&service.route_desc=="TRE OLMI"){
            service.linea="5A";
        }
        //6A Santi (as 25/26)
        if(service.linea=="6"&&service.route_desc=="SANTI"){
            service.linea="6A";
        }
        //6B Villanova (as 25/26)
        if(service.linea=="6"&&service.route_desc=="VILLANOVA"){
            service.linea="6B";
        }
        //7 GOTTARDI -> POLICLINICO GOTTARDI
        if(service.linea=="7"&&service.route_desc=="GOTTARDI"){
            service.route_desc="POLICLINICO GOTTARDI";
        }
        //7A STAZIONE FS -> GOTTARDI
        if(service.linea=="7A"&&service.route_desc=="STAZIONE FS"&&!service.route_code.includes("728")){
            service.route_desc="GOTTARDI";
        }
        //7A/ STAZIONE FS
        if(service.linea=="7A"&&service.route_desc=="STAZIONE FS"){
            service.linea="7A/";
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
        //9/ Autostazione
        if(service.linea=="9"&&service.route_desc=="AUTOSTAZIONE"){
            service.linea="9/";
        }
        //10A La Rocca
        if(service.linea=="10"&&service.route_desc=="LA ROCCA"){
            service.linea="10A";
        }
        //10S Liceo Sigonio
        if(service.linea=="10"&&service.route_desc=="LICEO SIGONIO"){
            service.linea="10S";
        }
        //10/ AUTOSTAZIONE
        if(service.linea=="10"&&service.route_desc=="AUTOSTAZIONE"){
            service.linea="10/";
        }
        //11/ Stazione FS
        if(service.linea=="11"&&service.route_desc=="STAZIONE FS"){
            service.linea="11/";
        }
        //12A Nazioni ma sono dei coglioni di merda
        if(service.linea=="12"&&service.route_code.includes("1280")||service.route_code.includes("1284")){
            service.linea="12A";
            service.route_desc="NAZIONI";
        }
        //12/ Fanti FS
        if(service.linea=="12"&&service.route_desc=="FANTI FS"){
            service.linea="12/";
        }
        //12S Garibaldi (Scuola)
        if(service.linea=="12"&&service.route_desc=="GARIBALDI"){
            service.linea="12S";
        }
        //12S Largo Garibaldi (Scuola)
        if(service.linea=="12"&&service.route_desc=="LARGO GARIBALDI"){
            service.linea="12S";
        }
        //13 S.ANNA -> SANT'ANNA (dio rincoglionito e dislessico)
        if(service.linea=="13"&&service.route_desc=="S.ANNA"){
            service.route_desc=="SANT'ANNA";
        }
        //13A Carcere
        if(service.linea=="13"&&service.route_desc=="CARCERI"){
            service.linea="13A";
        }
        //13F Variante di merda
        if(service.linea=="13"&&service.route_code.includes("133")){
            service.linea="13F";
        }
        //643 _ -> Polo Scolastico Sassuolo
        if(service.linea=="643"&&service.route_desc=="_"){
            service.route_desc="POLO SCOLASTICO SASSUOLO";
        }
        //Varianti vecchie AS 24/25
        /*
        //3A Vaciglio-Mattarella
        if(service.linea=="3"&&service.route_desc=="VACIGLIO MATTARELLA"){
            service.linea="3A";
        }
        //3A Portorico (Domenica)
        if(service.linea=="3"&&service.route_desc=="PORTORICO"){
            service.linea="3A";
        }
        //14A Nazioni
        if(service.linea=="14"&&service.route_desc=="NAZIONI"){
            service.linea="14A";
        }
        //15/ Santi
        if(service.linea=="15"&&service.route_desc=="SANTI"){
            service.linea="15/";
        }
        */
        //Varianti settembre 2025
        /*
        //3F MONTEFIORINO-Portorico (Domenica)
        if(service.linea=="3"&&service.route_desc=="PORTORICO"){
            service.linea="3F";
            service.route_desc="MONTEFIORINO-PORTORICO";
        }
        */
        //Cambia nomi modello bus (nuovo metodo per numero)
        //Lotti promiscui vengono messi in alto in modo che se c'è qualcosa di sbagliato (c'è) vengono sistemati poi dopo
        if(service.vehicle_code >= 300 && service.vehicle_code <= 891){
            service.model = "Irisbus Crossway";
        }
        if(service.vehicle_code >= 354 && service.vehicle_code <= 691){
            service.model = "Mercedes Integro O550";
        }
        if(service.vehicle_code >= 651 && service.vehicle_code <= 682){
            service.model = "Mercedes Citaro O530Ü";
        }
        if(service.vehicle_code >= 341 && service.vehicle_code <= 352){
            service.model = "MAN Lion's Regio";
        }
        if(service.vehicle_code >= 335 && service.vehicle_code <= 338){
            service.model = "Mercedes Integro O550";
        }
        if(service.vehicle_code >= 1 && service.vehicle_code <= 7){
            service.model = "Neoplan Electroliner";
        }

        if(service.vehicle_code >= 25 && service.vehicle_code <= 34){
            service.model = "CAM Busotto NGT";
        }

        if(service.vehicle_code >= 113 && service.vehicle_code <= 132){
            service.model = "Mercedes Citaro O530N Diesel";
        }

        if(service.vehicle_code >= 133 && service.vehicle_code <= 142){
            service.model = "Irisbus Cityclass CNG ATCM";
        }

        if(service.vehicle_code >= 143 && service.vehicle_code <= 146){
            service.model = "Mercedes Citaro O530N CNG";
        }

        if(service.vehicle_code >= 170 && service.vehicle_code <= 197){
            service.model = "Irisbus Citelis CNG EEV";
        }

        if(service.vehicle_code >= 198 && service.vehicle_code <= 200){
            service.model = "Solaris Urbino 12 III CNG";
        }

        if(service.vehicle_code >= 271 && service.vehicle_code <= 290){
            service.model = "MenariniBus Citymood CNG";
        }

        if(service.vehicle_code >= 4750 && service.vehicle_code <= 4763){
            service.model = "MenariniBus Citymood LNG";
        }

        if(service.vehicle_code >= 4770 && service.vehicle_code <= 4799){
            service.model = "Iveco Urbanway Mild Hybrid CNG";
        }

        if(service.vehicle_code >= 7901 && service.vehicle_code <= 7912){
            service.model = "Solaris Urbino 12 IV Hydrogen";
        }

        if(service.vehicle_code >= 209 && service.vehicle_code <= 212){
            service.model = "Solaris Urbino 18 III";
        }

        if(service.vehicle_code >= 323 && service.vehicle_code <= 324){
            service.model = "Irisbus Crossway LE";
        }

        if(service.vehicle_code >= 213 && service.vehicle_code <= 216){
            service.model = "BredaMenariniBus Avancity+ CNG 18";
        }

        if(service.vehicle_code >= 217 && service.vehicle_code <= 221){
            service.model = "MAN Lion's City G 3p ex Gamla";
        }

        if(service.vehicle_code >= 339 && service.vehicle_code <= 379){
            service.model = "Iveco Crossway LE";
        }

        if(service.vehicle_code >= 222 && service.vehicle_code <= 224){
            service.model = "MAN Lion's City G 4p ex Unibuss";
        }

        if(service.vehicle_code >= 251 && service.vehicle_code <= 270){
            service.model = "Solaris Urbino 12 III LE";
        }

        if(service.vehicle_code >= 229 && service.vehicle_code <= 232){
            service.model = "Mercedes Citaro O530G ex Koln";
        }

        if(service.vehicle_code >= 380 && service.vehicle_code <= 381){
            service.model = "MAN Lion's City L CNG ex TronderBilene";
        }

        if(service.vehicle_code >= 382 && service.vehicle_code <= 387){
            service.model = "Iveco Crossway LE Bianchi";
        }

        if(service.vehicle_code >= 4501 && service.vehicle_code <= 4507){
            service.model = "Setra S415 LE 2p ex Bolzano";
        }

        if(service.vehicle_code >= 4508 && service.vehicle_code <= 4513){
            service.model = "Setra S415 LF 3p ex Bolzano";
        }

        if(service.vehicle_code >= 4621 && service.vehicle_code <= 4628){
            service.model = "Iveco Crossway LE 14";
        }

        if(service.vehicle_code >= 4450 && service.vehicle_code <= 4466){
            service.model = "Iveco Crossway LE 12 CNG";
        }

        if(service.vehicle_code >= 4467 && service.vehicle_code <= 4468){
            service.model = "Iveco Crossway Line 12 CNG";
        }

        if(service.vehicle_code >= 4850 && service.vehicle_code <= 4852){
            service.model = "MAN Lion's City 19 CNG";
        }

        if(service.vehicle_code >= 4853 && service.vehicle_code <= 4854){
            service.model = "Mercedes Citaro O530GÜ ex Tper";
        }

        if(service.vehicle_code >= 692 && service.vehicle_code <= 699){
            service.model = "Irisbus Ares SFR117";
        }

        if(service.vehicle_code >= 354 && service.vehicle_code <= 691){
            service.model = "Mercedes Integro O550";
        }

        if(service.vehicle_code == 686){
            service.model = "Mercedes Integro O550 (Giallo)";
        }

        if(service.vehicle_code >= 325 && service.vehicle_code <= 334){
            service.model = "Irisbus Crossway ex Esercito Tedesco";
        }

        if(service.vehicle_code >= 4400 && service.vehicle_code <= 4426){
            service.model = "Iveco Crossway Line";
        }

        if(service.vehicle_code >= 4427 && service.vehicle_code <= 4449){
            service.model = "Scania Irizar i4 LNG";
        }
        if(service.vehicle_code >= 7901 && service.vehicle_code <= 7912){
            service.model = "Solaris Urbino IV Hydrogen";
        }
        
        //Dio filoviario arrostito sulla 750 volt
        if(service.vehicle_code>=35&&service.vehicle_code<=44){
            service.model="Solaris Trollino 12 IV";
            service.plate_num="MO0"+service.vehicle_code;
        }
        //Daily e altre cose
        if(service.model=="SOLARIS - URBINO U 8,9 LE"){
            service.model="Solaris 9 LE ex Piacenza";
        }
        if(service.vehicle_code>= 66 && service.vehicle_code<=67){
            service.model="Mercedes Sprinter";
        }
        if(service.vehicle_code>= 63 && service.vehicle_code<=64){
            service.model="Iveco Daily carr. Indicar";
        }
        if(service.vehicle_code==50){
            service.model="Iveco Daily";
        }
        if(service.vehicle_code==46){
            service.model="Mercedes Sprinter";
        }
        if(service.vehicle_code==961){
            service.model="Iveco Thesi";
        }
        if(service.vehicle_code==962){
            service.model="Ford Transit";
        }
        if(service.vehicle_code==954){
            service.model="Sprinter carr. Buzola";
        }
        if(service.vehicle_code==965){
            service.model="Iveco Daily";
        }
        if(service.vehicle_code==952){
            service.model="Iveco Daily carr. Gerbus";
        }
        /*
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
        */
        if(false){
            service.linea="99";
            service.route_desc="Glory to ATCM";
            service.model="Glory to ATCM";
            //service.vehicle_code="Glory to ATCM";
            service.service_tag="Glory to ATCM";
            service.plate_num="Glory to ATCM";
            service.num_passeggeri="Glory to ATCM";
            service.waypoint_code="Glory to ATCM";
            service.posti_totali="Glory to ATCM";
            service.journey_code="Glory to ATCM";
            service.delay="Glory to ATCM";
            service.occupancy_lstupd="Glory to ATCM";
            service.wp_desc="Glory to ATCM";
            service.service_code="Glory to ATCM";
            service.route_code="Glory to ATCM";
            service.reached_waypoint_code="Glory to ATCM";
            service.next_stop="Glory to ATCM";
            service.next_stop="Glory to ATCM";
            service.next_stop="Glory to ATCM";
            service.next_stop="Glory to ATCM";
            service.next_stop="Glory to ATCM";
        }
    });
}

function fixStopNames(stops){
    stops.forEach(element =>{
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
            element.properties.plate_num="Non disponibile";
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

//Cerca e aggiunge waypoint_code alle info veicolo
async function addNextStopCode(response,idMezzo){
    var item=await getBusList();
    item.features.forEach(element =>{
        if(element.properties.vehicle_code==idMezzo){
            response.data.features.forEach(bus =>{
                bus.properties.waypoint_code=element.properties.waypoint_code;
            });
        }
    });
}

//Urbanway e Menarini LNG hanno la pedana
function fixPedana(response){
    response.data.features.forEach(element => {
        if(element.properties.model=="Iveco Urbanway Mild Hybrid CNG"){
            element.properties.pedana=1;
        }
        if(element.properties.model=="MenariniBus Citymood LNG"){
            element.properties.pedana=1;
        }
    });
}

//Codici di merda delle route nel loro sito del cazzo
async function getCodiciDiMerda(){
    const response = await axios.get(lineedynUrl);
    const data = scrapeIt.scrapeHTML(response.data, {
        codes: {
            listItem: ".riga_copertura",
            data: {
                num: ".numero",
                shitCode: {
                    selector: ".testo_per_ricerca",
                    attr: "data-val"
                }
            }
        }
    })
    return data;
}

app.listen(port, () => {
    console.log(`API attiva su http://localhost:${port}`);
});
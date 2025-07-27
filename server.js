const express = require('express');
const axios = require('axios');
const cheerio = require('cheerio');
const cors = require('cors');
const fs = require('fs');

const app = express();
const port = 5001;

app.use(cors());
app.use(express.json());
// mettilo su / se lo vuoi sull'ip diretto
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
            for(var i=0;i<response.data.arrival.services.length;i++){
                //1A Modena Est
                if(response.data.arrival.services[i].service=="1"&&response.data.arrival.services[i].destination=="MODENA EST"){
                    response.data.arrival.services[i].service="1A";
                }
                //1A Polo Leonardo
                if(response.data.arrival.services[i].service=="1"&&response.data.arrival.services[i].destination=="POLO LEONARDO"){
                    response.data.arrival.services[i].service="1A";
                }
                //1B Ariete
                if(response.data.arrival.services[i].service=="1"&&response.data.arrival.services[i].destination=="ARIETE"){
                    response.data.arrival.services[i].service="1B";
                }
                //1B Ariete
                if(response.data.arrival.services[i].service=="1"&&response.data.arrival.services[i].destination=="ARIETE"){
                    response.data.arrival.services[i].service="1B";
                }
                //2A San Donnino
                if(response.data.arrival.services[i].service=="2"&&response.data.arrival.services[i].destination=="SAN DONNINO"){
                    response.data.arrival.services[i].service="2A";
                }
                //3A Vaciglio
                if(response.data.arrival.services[i].service=="3"&&response.data.arrival.services[i].destination=="VACIGLIO MATTARELLA"){
                    response.data.arrival.services[i].service="3A";
                }
                //3A Portorico (Domenica)
                if(response.data.arrival.services[i].service=="3"&&response.data.arrival.services[i].destination=="PORTORICO"){
                    response.data.arrival.services[i].service="3A";
                }
                //3A Vaciglio (Domenica)
                if(response.data.arrival.services[i].service=="3"&&response.data.arrival.services[i].destination=="VACIGLIO"){
                    response.data.arrival.services[i].service="3A";
                }
                //5A Tre Olmi
                if(response.data.arrival.services[i].service=="5"&&response.data.arrival.services[i].destination=="TRE OLMI"){
                    response.data.arrival.services[i].service="5A";
                }
                //9A Marzaglia Nuova
                if(response.data.arrival.services[i].service=="9"&&response.data.arrival.services[i].destination=="MARZAGLIA"){
                    response.data.arrival.services[i].service="9A";
                }
                //9C Rubiera
                if(response.data.arrival.services[i].service=="9"&&response.data.arrival.services[i].destination=="RUBIERA"){
                    response.data.arrival.services[i].service="9C";
                }
                //9/ Stazione FS
                if(response.data.arrival.services[i].service=="9"&&response.data.arrival.services[i].destination=="STAZIONE FS"){
                    response.data.arrival.services[i].service="9/";
                }
                //10A La Rocca
                if(response.data.arrival.services[i].service=="10"&&response.data.arrival.services[i].destination=="LA ROCCA"){
                    response.data.arrival.services[i].service="10A";
                }
                //13A Carcere
                if(response.data.arrival.services[i].service=="13"&&response.data.arrival.services[i].destination=="CARCERI"){
                    response.data.arrival.services[i].service="13A";
                }
                //14A Nazioni
                if(response.data.arrival.services[i].service=="14"&&response.data.arrival.services[i].destination=="NAZIONI"){
                    response.data.arrival.services[i].service="14A";
                }
                //15/ Santi
                if(response.data.arrival.services[i].service=="15"&&response.data.arrival.services[i].destination=="SANTI"){
                    response.data.arrival.services[i].service="15/";
                }
            }
            // Step 1: Mappa i servizi per codice_corsa divisi per tipo
            const plannedMap = new Map();
            const realtimeMap = new Map();

            response.data.arrival.services.forEach(service => {
            if (service.type === "planned") {
                plannedMap.set(service.codice_corsa, service);
            } else if (service.type === "realtime") {
                realtimeMap.set(service.codice_corsa, service);
            }
            });

            // Step 2: Filtra i servizi
            const filteredServices = response.data.arrival.services.filter(service => {
            if (service.type === "realtime") return true;
            return !realtimeMap.has(service.codice_corsa);
            });

            // Step 3: Aggiungi "delay" dove possibile
            filteredServices.forEach(service => {
            if (service.type === "realtime") {
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
                    "waypoint" : "MO3053"
                }
            });
        }
    }
});

app.listen(port, () => {
    console.log(`API attiva su http://localhost:${port}`);
});
//Stile grafico
// https://leaflet-extras.github.io/leaflet-providers/preview/
//Necessita licenza
// var Stadia_OSMBright = L.tileLayer('https://tiles.stadiamaps.com/tiles/osm_bright/{z}/{x}/{y}{r}.png', {
//     minZoom: 12,
//     attribution: '&copy; <a href="https://stadiamaps.com/">Stadia Maps</a>, &copy; <a href="https://openmaptiles.org/">OpenMapTiles</a> &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors'
// });


var Stadia_OSMBright = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    minZoom: 11,
    attribution: '© OpenStreetMap',
    referrerPolicy : 'strict-origin-when-cross-origin'
});

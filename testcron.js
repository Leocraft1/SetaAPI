const cron = require('node-cron');

const testTsk = () =>{
    console.log("Ciao, adesso sono le "+ new Date());
};

//cron.schedule("* * * * * *", testTsk);
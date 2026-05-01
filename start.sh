cd /root/seta

# Rest of script from old node
# Pull and sync from SetaAPI repository
cd /DATA/AppData/Seta-APIs
git stage *
echo Staged
git restore server.js
git restore package.json
echo Restored
git commit -m "Automatic Update Commit from Serverissimo"
echo Committed
git push https://github.com/Leocraft1/SetaAPI main
echo Pushed

git pull
echo Pulled
npm install

node server.js

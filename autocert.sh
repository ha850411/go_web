# autocert.sh -> 自動更新憑證

#go project path
projectPath=/home/gopath/go_web
myHost=easonlee.site

#備份原來的憑證
mv $projectPath/certs/fullchain.crt $projectPath/certs/fullchain.crt.bk.$(date '+%Y%m%d')
mv $projectPath/certs/privkey.key $projectPath/certs/privkey.key.bk.$(date '+%Y%m%d')

#產生憑證
certbot certonly --webroot -w /home/gopath/go_web -d $myHost --email le850411@gmail.com

#複製憑證
cp /etc/letsencrypt/live/$myHost/fullchain.pem $projectPath/certs/fullchain.crt
cp /etc/letsencrypt/live/$myHost/privkey.pem $projectPath/certs/privkey.key

#重啟 docker 重新掛載證書)
docker-compose restart
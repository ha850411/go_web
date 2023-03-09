# crontab : * * 0 0 0 每日儲存 log
sourceFile=/home/gopath/go_web/logs/run.log
targetFile=/home/gopath/go_web/logs/log_$(date '+%Y%m%d').log
cp $sourceFile $targetFile
cat /dev/null > $sourceFile
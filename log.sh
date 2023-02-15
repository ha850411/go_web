# crontab : * * 0 0 0 每日儲存 log
cp ./logs/run.log ./logs/log_$(date '+%Y%m%d').log
cat /dev/null > ./logs/run.log
go build -o go_main
nohup ./go_main >> ./logs/run.log 2>&1 &
# 每日儲存 log
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do
    DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
    SOURCE="$( readlink "$SOURCE" )"
    [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

sourceFile=$DIR/logs/run.log
targetFile=$DIR/logs/log_$(date '+%Y%m%d').log

current_month=$(date +"%Y%m")
preMonthFile=$DIR/logs/log_$(date -v -1m -j -f "%Y%m" "$current_month" +"%Y%m")*

cp $sourceFile $targetFile
cat /dev/null > $sourceFile

# 刪除前一個月 log
rm -rf $preMonthFile
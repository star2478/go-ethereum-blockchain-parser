package lib
import (
		"fmt"
		"log"
		"os"
		"runtime"
		//"time"
		//"reflect"
)

func GetBlockNumberByTime(timeBegin string, timeEnd string) int,int {
	log.Print("begin GetBlockNumberByTime, timeBegin=", timeBegin, " timeEnd=", timeEnd)
	fileName := "../blocktime"
	fd, err := os.Open(fileName)
        if err != nil {
                log.Fatal("file=",fileName," open fail!!", err);
        }
	defer fd.Close()
	readBuf := bufio.NewReader(fd)
	for {
		line, err := readBuf.ReadString('\n')
        	if err != io.EOF {
			break
		}
                line = strings.TrimSpace(line)
                arrLine = strings.Split(lineFrom, " ")
                blockTime := arrLine[0]
		blockNumber := strconv.Atoi(arrLine[1])
	}
}

package lib
import (
		"log"
		"os"
		"io"
		"strings"
		"strconv"
		"bufio"
)

func GetBlockNumberByTime(timeBegin string, timeEnd string) (int,int) {
	log.Print("begin GetBlockNumberByTime, timeBegin=", timeBegin, " timeEnd=", timeEnd)
	fileName := "blocktime"
	fd, err := os.Open(fileName)
        if err != nil {
                log.Fatal("file=",fileName," open fail!!", err);
        }
	defer fd.Close()
	blockNumberBegin := -1
	blockNumberEnd := -1
	lastScanBlockNumber := -1
	blockTime := ""
	readBuf := bufio.NewReader(fd)
	for {
		line, err := readBuf.ReadString('\n')
        	if err == io.EOF {
			break
		}
                line = strings.TrimSpace(line)
                arrLine := strings.Split(line, " ")
                blockTime = arrLine[0]
		blockNumber,_ := strconv.Atoi(arrLine[1])
		if blockNumberBegin == -1 && timeBegin <= blockTime {
			blockNumberBegin = blockNumber
		}
		if timeEnd < blockTime {
			blockNumberEnd = lastScanBlockNumber
			break
		}
		lastScanBlockNumber = blockNumber
	}
	if blockNumberBegin != -1 && timeEnd >= blockTime {
		blockNumberEnd = lastScanBlockNumber
	}
  	if blockNumberBegin == -1 || blockNumberEnd == -1 {
        	log.Fatal("Cannot find blockNumber!!! timeBegin=", timeBegin, " timeEnd=", timeEnd, ", blockNumberBegin=", blockNumberBegin, " blockNumberEnd=", blockNumberEnd);
  	}	
	log.Print("finish GetBlockNumberByTime, blockNumberBegin=", blockNumberBegin, " blockNumberEnd=",blockNumberEnd)
	return blockNumberBegin, blockNumberEnd
}

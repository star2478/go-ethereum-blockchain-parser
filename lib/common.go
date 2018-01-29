package lib
import (
		"log"
		"os"
		"os/exec"
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
        	log.Fatal("Cannot find blockNumber!!! timeBegin=", timeBegin, " timeEnd=", timeEnd, ", blockNumberBegin=", blockNumberBegin, " blockNumberEnd=", blockNumberEnd, ". You can try to run 'go run synBlockTime.go' to get current block time before this");
  	}	
	log.Print("finish GetBlockNumberByTime, blockNumberBegin=", blockNumberBegin, " blockNumberEnd=",blockNumberEnd)
	return blockNumberBegin, blockNumberEnd
}

func PathExists(path string) (bool) {
        _, err := os.Stat(path)
        if err == nil {
                return true
        }
        return false
}
func GetAndCheckDir(dirname string) string {
        curDir, _ := os.Getwd()  //
        newDir := curDir +"/" + dirname
        if PathExists(newDir) == false {
                err := os.Mkdir(newDir, os.ModePerm)  //
                if err != nil {
                        log.Fatal(newDir, " create fail! ", err)
                        //return nil
                }
        }
        return newDir
}
func ExecCmd(command string, isFatalExit bool) {

    log.Print("begin to run command=" + command);
    cmd := exec.Command("/bin/bash", "-c", command)
    err := cmd.Run()
    if err != nil {
	if isFatalExit == true {
        	log.Fatal("Fail!!!cmd="+command+" run fail!", err)
	} else {
        	log.Print("Fail!!!cmd="+command+" run fail!", err)

	}
    }
}

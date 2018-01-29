package main

import (
  "bytes"
  "sort"
  "fmt"
  "log"
  "strconv"
  "strings"
  "time"
  "runtime"
  "encoding/json"
  "os"
  "io"
  "bufio"
  "./lib"
  //"./golang-set"
//  "github.com/ethereum/go-ethereum/common"
//  "github.com/ethereum/go-ethereum/common/hexutil"
//  "github.com/deckarep/golang-set"
)

var MAX_NUM_SHARE = 10
var MAX_NUM_RETRY = 3
var APPEND_FILE_ACCOUNT_MAX_NUM = 100
var APPEND_FILE_LOOP_MAX_NUM = 1000
var APPEND_FILE_ACCOUNT_TX_MAX_NUM = 10000

type Timeline struct {
        Txtime string
        TxAmt string
        TxType string
}
type TimelineSlice [] Timeline
func (a TimelineSlice) Len() int {   // Öд Len() ·½·¨
 return len(a)
}
func (a TimelineSlice) Swap(i, j int){  // Öд Swap() ·½·¨
 a[i], a[j] = a[j], a[i]
}
func (a TimelineSlice) Less(i, j int) bool { // Öд Less() ·½·¨£¬ ´Ӵó¡ÅÐ
 return a[j].Txtime < a[i].Txtime
}
func setTxMap(m3 map[string]string, key string, value string) {
        if _, ok := m3[key]; ok {
                m3[key] += "&" + value
        } else {
                m3[key] = value
        }
}
func getFileContent(txMap map[string]string) bytes.Buffer {
	fileContent := bytes.Buffer{}
	for account, v := range txMap {
		if len(v) <= 0 {
			continue;
        	}
        	timelineList := strings.Split(v, "&")
        	timelineSlice := make([]Timeline, len(timelineList))
        	for i := range timelineList {
                	timelineItemArr := strings.Split(timelineList[i], "|")
                	timeline := &Timeline{
                    		timelineItemArr[0],
                    		timelineItemArr[1],
                    		timelineItemArr[2],
                	}
                	timelineSlice[i] = *timeline
        	}
		sort.Sort(TimelineSlice(timelineSlice))
		jsonStr, err := json.Marshal(timelineSlice)
        	if err != nil {
                	log.Print("Marshal failed!!!", err)
			continue
        	}
		fileContent.WriteString(account + "\t")
		fileContent.WriteString(strconv.Itoa(len(timelineList)) + "\t")
		fileContent.WriteString(string(jsonStr) + "\n")
	}
	return fileContent
}

func main() {
	if len(os.Args) < 3 {
                log.Fatal("Param Invalid!!! go run getTxTimelineGroupByAccount.go [timeFrom] [timeTo], eg. go run getTxTimelineGroupByAccount.go 2018-01-01-00-00-00 2018-02-01-00-00-00")
        }
  	fmt.Println("getTxTimelineGroupByAccount begin==================");
  	timeBegin := time.Now().Unix()  
  	MULTICORE := runtime.NumCPU()
  	runtime.GOMAXPROCS(MULTICORE)
	timeFrom := os.Args[1]
        timeTo := os.Args[2]
        if timeFrom >= timeTo {
                log.Fatal("timeFrom=", timeFrom, " >= timeTo=", timeTo)
        }
	dir := lib.GetAndCheckDir("tx")
  	fromFileName := dir + "/" + timeFrom + "-" + timeTo + "-from-sort"
  	toFileName := dir + "/" + timeFrom + "-" + timeTo + "-to-sort"
  	desFileName := dir + "/" + timeFrom + "-" + timeTo + "-timeline"
	fd1, err1 := os.Open(fromFileName)
	if err1 != nil {
		log.Fatal("file=",fromFileName," open fail!! First of all, you must run 'go run getTxByTime.go ", timeFrom, " ", timeTo, "' to get transactions! err=" , err1);
	}
	fd2, err2 := os.Open(toFileName)
	if err2 != nil {
		log.Fatal("file=",toFileName," open fail!! First of all, you must run 'go run getTxByTime.go ", timeFrom, " ", timeTo, "' to get transactions! err=" , err1);
	}
        desFd,err3 := os.OpenFile(desFileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err3 != nil {
		log.Fatal("file=",desFileName," open fail!!", err3);
	}
	fromLineNum := 0
	toLineNum := 0
	loopCount := 0
	readBufFrom := bufio.NewReader(fd1)
	readBufTo := bufio.NewReader(fd2)
			
	from := ""
	to := ""
	//txTime := ""
	//txType := ""
	//timelineString := ""
	var arrFrom [] string
	var arrTo [] string 
	txMap := map[string]string{}

	lineFrom, errReadFrom := readBufFrom.ReadString('\n')
	if errReadFrom != io.EOF {
		fromLineNum++
		lineFrom = strings.TrimSpace(lineFrom)
		arrFrom = strings.Split(lineFrom, "\t")
		from = arrFrom[2]
	}
	lineTo, errReadTo := readBufTo.ReadString('\n')
	if errReadTo != io.EOF {
		toLineNum++
		lineTo = strings.TrimSpace(lineTo)
		arrTo = strings.Split(lineTo, "\t")
		to = arrTo[3]
	}				
			
			
	for {		
		if errReadFrom == io.EOF || errReadTo == io.EOF {
			break
		}
		if from == to {
			tmp := from
			accountTxNumFrom :=0
			accountTxNumTo :=0
			for {
				setTxMap(txMap, from, arrFrom[0] + "|" + arrFrom[4] + "|out")
				lineFrom, errReadFrom = readBufFrom.ReadString('\n')
				fromLineNum++
				if errReadFrom == io.EOF {
					break
				}
				lineFrom = strings.TrimSpace(lineFrom)
				arrFrom = strings.Split(lineFrom, "\t")
				from = arrFrom[2]
				if from > tmp {
					break
				}
				accountTxNumFrom++
				if  accountTxNumFrom >=  APPEND_FILE_ACCOUNT_TX_MAX_NUM {
					fileContent := getFileContent(txMap)		
					wBuf:=[]byte(fileContent.String())
                			desFd.Write(wBuf)
					log.Print("write account's tx into file=", desFileName, " fromLineNum=",fromLineNum, " toLineNum=", toLineNum)
					txMap = map[string]string{}
					accountTxNumFrom = 0
				}
			}
			for {
				setTxMap(txMap, to, arrTo[0] + "|" + arrTo[4] + "|in")
				lineTo, errReadTo = readBufTo.ReadString('\n')
				toLineNum++
				if errReadTo == io.EOF {
					break
				}
				lineTo = strings.TrimSpace(lineTo)
				arrTo = strings.Split(lineTo, "\t")
				to = arrTo[3]
				if to > tmp {
					break
				}
				accountTxNumTo++
				if  accountTxNumTo >=  APPEND_FILE_ACCOUNT_TX_MAX_NUM {
					fileContent := getFileContent(txMap)		
					wBuf:=[]byte(fileContent.String())
                			desFd.Write(wBuf)
					log.Print("write account's tx into file=", desFileName, " fromLineNum=",fromLineNum, " toLineNum=", toLineNum)
					txMap = map[string]string{}
					accountTxNumTo = 0
				}
			}
		} else if from < to {
			setTxMap(txMap, from, arrFrom[0] + "|" + arrFrom[4] + "|out")
			lineFrom, errReadFrom = readBufFrom.ReadString('\n')
			fromLineNum++
			if errReadFrom == io.EOF {
				break
			}
			lineFrom = strings.TrimSpace(lineFrom)
			arrFrom = strings.Split(lineFrom, "\t")
			from = arrFrom[2]
		} else {
			setTxMap(txMap, to, arrTo[0] + "|" + arrTo[4] + "|in")
			lineTo, errReadTo = readBufTo.ReadString('\n')
			toLineNum++
			if errReadTo == io.EOF {
				break
			}
			lineTo = strings.TrimSpace(lineTo)
			arrTo = strings.Split(lineTo, "\t")
			to = arrTo[3]
		}
		loopCount++
		if len(txMap) >= APPEND_FILE_ACCOUNT_MAX_NUM || loopCount >= APPEND_FILE_LOOP_MAX_NUM {
			fileContent := getFileContent(txMap)		
			wBuf:=[]byte(fileContent.String())
                	desFd.Write(wBuf)
			log.Print("write account's tx into file=", desFileName, " fromLineNum=",fromLineNum, " toLineNum=", toLineNum)
			txMap = map[string]string{}
			loopCount = 0
		}
	}
	log.Print("==================out of loop================fromLineNum=",fromLineNum, " toLineNum=", toLineNum)

	if len(txMap) >= APPEND_FILE_ACCOUNT_MAX_NUM || loopCount >= APPEND_FILE_LOOP_MAX_NUM {
		fileContent := getFileContent(txMap)		
		wBuf:=[]byte(fileContent.String())
                desFd.Write(wBuf)
		log.Print("write account's tx into file=", desFileName, " fromLineNum=",fromLineNum, " toLineNum=", toLineNum)
		txMap = map[string]string{}
	}

	if errReadFrom != io.EOF && errReadTo == io.EOF {
		for {
			setTxMap(txMap, from, arrFrom[0] + "|" + arrFrom[4] + "|out")
			lineFrom, errReadFrom = readBufFrom.ReadString('\n')
			fromLineNum++
			if errReadFrom == io.EOF {
				break
			}
			lineFrom = strings.TrimSpace(lineFrom)
			arrFrom = strings.Split(lineFrom, "\t")
			from = arrFrom[2]
		}
	}
	if errReadFrom == io.EOF && errReadTo != io.EOF {
		for {
			setTxMap(txMap, to, arrTo[0] + "|" + arrTo[4] + "|in")
			lineTo, errReadTo = readBufTo.ReadString('\n')
			toLineNum++
			if errReadTo == io.EOF {
				break
			}
			lineTo = strings.TrimSpace(lineTo)
			arrTo = strings.Split(lineTo, "\t")
			to = arrTo[3]
		}
	}
	
	if len(txMap) > 0 {
		fileContent := getFileContent(txMap)		
		wBuf:=[]byte(fileContent.String())
                desFd.Write(wBuf)
		log.Print("write account's tx into file=", desFileName, " fromLineNum=",fromLineNum, " toLineNum=", toLineNum)
	}

        fd1.Close() 
        fd2.Close() 
        desFd.Close() 

  timeEnd := time.Now().Unix()  
  log.Print("getTxTimelineGroupByAccount finish, cost=", (timeEnd - timeBegin), "s")
}

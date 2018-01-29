package main

import (
  "bytes"
  //"sort"
  "fmt"
  "log"
  "strconv"
  "strings"
  "time"
  "runtime"
  //"encoding/json"
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
var APPEND_FILE_LOOP_MAX_NUM = 10000
var APPEND_FILE_ACCOUNT_TX_MAX_NUM = 10000
var txMap = map[string]string{}
var txMapFrom = map[string]*Count{}
var txMapTo = map[string]*Count{}

type Count struct {
        TxCount int
        TxAmtTotal float64
}
func initMap(key string) {
	
        if _, ok := txMap[key]; !ok {
                txMap[key] = ""
        } 
        if _, ok := txMapFrom[key]; !ok {
                txMapFrom[key] = &Count{0, 0}
        }
        if _, ok := txMapTo[key]; !ok {
                txMapTo[key] = &Count{0, 0}
        }
}
func incTxMap(tx map[string]*Count, key string, incCount int, incAmt float64) {
	initMap(key)
	//		log.Print(key, " :", incCount," ", incAmt)
        tx[key].TxCount += incCount
        tx[key].TxAmtTotal += incAmt
	//		log.Print(key, " ", tx[key].TxCount, " ", tx[key].TxAmtTotal, " ", incCount, " ", incAmt)
}

func main() {
	if len(os.Args) < 3 {
                log.Fatal("Param Invalid!!! go run getTxCountGroupByAccount.go [timeFrom] [timeTo], eg. go run getTxCountGroupByAccount.go 2018-01-01-00-00-00 2018-02-01-00-00-00")
        }
  	fmt.Println("getTxCountGroupByAccount begin==================");
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
        desFileName := dir + "/" + timeFrom + "-" + timeTo + "-count"
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
	readBufFrom := bufio.NewReader(fd1)
	readBufTo := bufio.NewReader(fd2)
		
	from := ""
	to := ""
	//txTime := ""
	//txType := ""
	//timelineString := ""
	var arrFrom [] string
	var arrTo [] string 
	loopCount := 0

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
			for {
				amtFrom,_ := strconv.ParseFloat(arrFrom[4],64)
				incTxMap(txMapFrom, from, 1, amtFrom)
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
			}
			for {
				amtTo,_ := strconv.ParseFloat(arrTo[4],64)
				incTxMap(txMapTo, to, 1, float64(amtTo))
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
			}
		} else if from < to {
			amtFrom,_ := strconv.ParseFloat(arrFrom[4],64)
                        incTxMap(txMapFrom, from, 1, amtFrom)
			lineFrom, errReadFrom = readBufFrom.ReadString('\n')
			fromLineNum++
			if errReadFrom == io.EOF {
				break
			}
			lineFrom = strings.TrimSpace(lineFrom)
			arrFrom = strings.Split(lineFrom, "\t")
			from = arrFrom[2]
		} else {
			amtTo,_ := strconv.ParseFloat(arrTo[4],64)
			incTxMap(txMapTo, to, 1, amtTo)
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
		if loopCount % APPEND_FILE_LOOP_MAX_NUM == 0 {
			log.Print("loopCount=", loopCount)
		} 
	}
	log.Print("==================finish loop================fromLineNum=",fromLineNum, " toLineNum=", toLineNum)

	if errReadFrom != io.EOF && errReadTo == io.EOF {
		for {
			amtFrom,_ := strconv.ParseFloat(arrFrom[4],64)
                        incTxMap(txMapFrom, from, 1, amtFrom)
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
			amtTo,_ := strconv.ParseFloat(arrTo[4],64)
			incTxMap(txMapTo, to, 1, amtTo)
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
	
	log.Print("begin to write result file")
	loopCount = 0
	fileContent := bytes.Buffer{}
	for account,_ := range txMap {
		countFrom := txMapFrom[account] 
		countTo := txMapTo[account] 
		countTotal := countFrom.TxCount + countTo.TxCount
		amtTotal := countFrom.TxAmtTotal + countTo.TxAmtTotal
		fileContent.WriteString(account)
		fileContent.WriteString(" " + strconv.Itoa(countTotal) + " " + strconv.Itoa(countFrom.TxCount) + " " + strconv.Itoa(countTo.TxCount))
		fileContent.WriteString(" " + strconv.FormatFloat(amtTotal, 'f', -1, 64) + " " + strconv.FormatFloat(countFrom.TxAmtTotal, 'f', -1, 64) + " " + strconv.FormatFloat(countTo.TxAmtTotal, 'f', -1, 64) + "\n")
		if loopCount % APPEND_FILE_LOOP_MAX_NUM == 0 {
			wBuf:=[]byte(fileContent.String())
                	desFd.Write(wBuf)
			log.Print("write account's tx into file=", desFileName, " loopCount=", loopCount)
			fileContent = bytes.Buffer{}
		}
		loopCount++
    	}
	if len(fileContent.String()) > 0 {
		wBuf:=[]byte(fileContent.String())
               	desFd.Write(wBuf)
		log.Print("write account's tx into file=", desFileName, " loopCount=", loopCount)
	}	

        fd1.Close() 
        fd2.Close() 
        desFd.Close() 

  timeEnd := time.Now().Unix()  
  log.Print("getTxCountGroupByAccount finish, cost=", (timeEnd - timeBegin), "s")
}

package main

import (
  "bytes"
  "fmt"
  "log"
  "strconv"
  "strings"
  "os/exec" 
  "math"
  "time"
  "runtime"
  "encoding/json"
  "os"
  "io"
  "bufio"
//  "./go-ethereum/common"
//  "./go-ethereum/common/hexutil"
  //"./golang-set"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/common/hexutil"
//  "github.com/deckarep/golang-set"
)

var MAX_NUM_SHARE = 10
var MAX_NUM_RETRY = 3
var APPEND_FILE_NUM = 100
var c = make(chan string, 100)

type ResultOfGetBalance struct {
    Jsonrpc    string
    Id    int
    Result	string//hexutil.Uint64
}
type RPCTransaction struct {
        BlockHash        common.Hash     `json:"blockHash"`
        BlockNumber      *hexutil.Big    `json:"blockNumber"`
        From             *common.Address  `json:"from"`
        Gas              hexutil.Uint64  `json:"gas"`
        GasPrice         *hexutil.Big    `json:"gasPrice"`
        Hash             common.Hash     `json:"hash"`
        Input            hexutil.Bytes   `json:"input"`
        Nonce            hexutil.Uint64  `json:"nonce"`
        To               *common.Address `json:"to"`
        TransactionIndex hexutil.Uint    `json:"transactionIndex"`
        Value            *hexutil.Big    `json:"value"`
        V                *hexutil.Big    `json:"v"`
        R                *hexutil.Big    `json:"r"`
        S                *hexutil.Big    `json:"s"`
}

type Block struct {
        UncleHashes  []common.Hash    `json:"uncles"`
        Hash         common.Hash      `json:"hash"`
        Transactions []RPCTransaction `json:"transactions"`
}

var tenToAny map[int]string = map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n", 24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t", 30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z", 36: ":", 37: ";", 38: "<", 39: "=", 40: ">", 41: "?", 42: "@", 43: "[", 44: "]", 45: "^", 46: "_", 47: "{", 48: "|", 49: "}", 50: "A", 51: "B", 52: "C", 53: "D", 54: "E", 55: "F", 56: "G", 57: "H", 58: "I", 59: "J", 60: "K", 61: "L", 62: "M", 63: "N", 64: "O", 65: "P", 66: "Q", 67: "R", 68: "S", 69: "T", 70: "U", 71: "V", 72: "W", 73: "X", 74: "Y", 75: "Z"}
 
 
func decimalToAny(num, n int) string {
 new_num_str := ""
 var remainder int
 var remainder_string string
 for num != 0 {
  remainder = num % n
  if 76 > remainder && remainder > 9 {
   remainder_string = tenToAny[remainder]
  } else {
   remainder_string = strconv.Itoa(remainder)
  }
  new_num_str = remainder_string + new_num_str
  num = num / n
 }
 return new_num_str
}

func findkey(in string) int {
 result := -1
 for k, v := range tenToAny {
  if in == v {
   result = k
  }
 }
 return result
}
//187b2445cfb31e42b6a
func anyToDecimal(num string, n int) float64 {
 var new_num float64
 new_num = 0.0
 nNum := len(strings.Split(num, "")) - 1
 for _, value := range strings.Split(num, "") {

  tmp := float64(findkey(value))
  if tmp != -1 {
   new_num = new_num + tmp*math.Pow(float64(n), float64(nNum))
   nNum = nNum - 1
  } else {
   break
  }
 }
 return float64(new_num)
}

func exec_shell(s string, accountNo string, retry int, fileContent *string) {
    cmd := exec.Command("/bin/bash", "-c", s)
    var out bytes.Buffer

    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        //log.Fatal("cmd="+s+" run fail!", err)
	retry--
        log.Print("Fatal!!!retry=",retry," cmd="+s+" run fail!", err)
	if(retry > 0) {
		exec_shell(s, accountNo, retry, fileContent)
    	}
    }
    
    stb := &ResultOfGetBalance{}
    //fmt.Printf("-----cmd=%s, result=%s\n", s, out.String())
    err = json.Unmarshal([]byte(out.String()), &stb)
    
    if err != nil {
	retry--;
        log.Print("umarshal fail!!!retry=",retry," accountNo=",accountNo," cmd=", s, " input=", out.String()," ", err)
	if(retry > 0) {
		exec_shell(s, accountNo, retry, fileContent)
    	}
    } 
    //fmt.Printf("----result=%s:%d\n", stb.Jsonrpc, stb.Id)
    balance := stb.Result
    balance = balance[2 : len(balance)]	// get rid of 0x prefix
    price := anyToDecimal(balance, 16)/1000000000000000000
    *fileContent += accountNo  + " " + strconv.FormatFloat(price, 'f', -1, 64) + "\n"; 
    //if balance <= 0 {
	//log.Print("block=", blockNumber, " balance len <= 0")
//	return
  //  }

    //fmt.Printf("----cmd=%s balance=%s\n", s, balance)
}

func pathExists(path string) (bool) {
        _, err := os.Stat(path)
        if err == nil {
                return true
        }
        return false
}

func getBalance(srcFileNamePrefix string, desFileNamePrefix string, sufix int) {
	srcFileName := srcFileNamePrefix + "_" + strconv.Itoa(sufix)
        desFileName := desFileNamePrefix + "_" + strconv.Itoa(sufix)
	log.Print("taskId=", desFileName, " begin! srcFileName=", srcFileName)
        fd1, err1 := os.Open(srcFileName)
        if err1 != nil {
                log.Fatal("file=",srcFileName," open fail!!", err1);
        }
        fd2,err2 := os.OpenFile(desFileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
        if err2 != nil {
                log.Fatal("file=",desFileName," open fail!!", err2);
        }
	fileContent := ""
        i := 0
        readBuf := bufio.NewReader(fd1)
        for {
                accountNo, err := readBuf.ReadString('\n')
                if err != nil {
                        if err == io.EOF {
                                break
                        }
                        return
                }
                accountNo = strings.TrimSpace(accountNo)
		command := "curl -X POST --data '{\"jsonrpc\":\"3.0\",\"method\":\"eth_getBalance\",\"params\":[\""+accountNo+"\",\"latest\"],\"id\":1}' -H \"Content-type: application/json;charset=UTF-8\"  localhost:8545";
    		exec_shell(command, accountNo, MAX_NUM_RETRY, &fileContent)
                if i%APPEND_FILE_NUM == 0 {
                        if len(fileContent) > 0 {
                                buf:=[]byte(fileContent)
                                fd2.Write(buf)
                                fileContent = ""
                                log.Print("write balance into file=", desFileName, " i=",i)
                        } else {
                                fileContent = ""
                                log.Print("fileContent len <= 0, not need to write balance into file=", desFileName, " i=",i)
                        }
                }
                i++
//break//////////////////////////////////////////////
        }
        if len(fileContent) > 0 {
                buf:=[]byte(fileContent)
                fd2.Write(buf)
                log.Print("write balance into file=", desFileName, " i=",i)
        }
        fd1.Close() 
        fd2.Close() 
	c <- desFileName
}

func main() {
  if len(os.Args) < 4 {
  	log.Fatal("Param Invalid!!! go run getBalanceMulti.go [srcFileNamePrefix] [desFileNamePrefix] [concurrentNum]")
  }
  fmt.Println("getBalance begin==================");
  timeBegin := time.Now().Unix()  
  MULTICORE := runtime.NumCPU()
  runtime.GOMAXPROCS(MULTICORE)
  srcFileNamePrefix := os.Args[1]
  desFileNamePrefix := os.Args[2]
  concurrentNum,err := strconv.Atoi(os.Args[3])
  if err != nil {
	log.Fatal("Param fail!!! concurrentNum=", os.Args[3], " must be integer! ", err);
  }
  for sufix := 1; sufix <= concurrentNum; sufix++ {
	go getBalance(srcFileNamePrefix, desFileNamePrefix, sufix)
  }
  for i := 1; i <= concurrentNum; i++ {
	taskId := <- c
        log.Print(taskId, " finish");
  }

  timeEnd := time.Now().Unix()  
  log.Print("getBalance finish, cost=", (timeEnd - timeBegin), "s")
}

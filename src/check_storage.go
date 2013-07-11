package main

import (
        "os"
        "fmt"
        "time"
	"flag"
        "math/big"
        snmp "github.com/soniah/gosnmp"
)

// oid valid for Linux system

var ERRORS map[string]int

var DISKSIZE string = "1.3.6.1.2.1.25.2.3.1.5"
var DISKUSED string = "1.3.6.1.2.1.25.2.3.1.6"
var DEVICE string = "1"
var HOST string = "127.0.0.1"
var WARNING int = 0
var CRITICAL int = 0

func check_config() {
	flag.StringVar(&DISKSIZE, "S", "1.3.6.1.2.1.25.2.3.1.5", "Disk total size oid")
	flag.StringVar(&DISKUSED, "U", "1.3.6.1.2.1.25.2.3.1.6", "Disk used size oid")
	flag.StringVar(&DEVICE, "d", "1", "Device id for oid")
	flag.StringVar(&HOST, "H", "", "Snmp server to query")
	flag.IntVar(&WARNING, "w", 80, "Value for warning")
	flag.IntVar(&CRITICAL, "c", 95, "Value for warning")
	flag.Parse()
	if flag.Lookup("H").Value.String()=="" {
		fmt.Println("Missing paramenter -H is required")
		flag.PrintDefaults()
		os.Exit(2)
	}
}

func outputWriter(x int, u *big.Int) string{
	status := "OK"
	if x > WARNING {
		if x > CRITICAL {
			status = "CRITICAL"
		} else {
			status = "WARNING"
		}
	}
		
	fmt.Printf("STORAGE %s: actual value %d%%\n", status, x)
	return status
}

func main() {
	ERRORS = map[string]int { 
		"OK":0,
		"WARNING":1,
		"CRITICAL":2,
		"UNKNOWN":3,
		"DEPENDENT":4,
	}
	check_config()
        s := &snmp.GoSNMP {
		Port:         161,
                Community:   "public",
                Version:     snmp.Version1,
		Timeout:     time.Duration(5) * time.Second,
                Target:      HOST,
        }
        size := big.NewInt(0);
        used := big.NewInt(0);
        err := s.Connect()
        if err != nil {
                fmt.Println("merda")
                os.Exit(ERRORS["UNKNOWN"])
        }
        defer s.Conn.Close()
        resp, err := s.Get([]string{DISKSIZE + "." + DEVICE, DISKUSED + "." + DEVICE})
        if err != nil {
                fmt.Println("err", err)
		os.Exit(ERRORS["UNKNOWN"])
        }
        for _, variable := range resp.Variables {
                if variable.Name == DISKSIZE+"."+DEVICE {
                        size = snmp.ToBigInt(variable.Value)
                } else {
                        used = snmp.ToBigInt(variable.Value)
                }
        }

        res := big.NewInt(0)
        mm := big.NewInt(0)
	ceil := big.NewInt(2)
        res = res.Mul(used, big.NewInt(100))
        res,mm = res.DivMod(res,size,mm)
	ceil = ceil.Div(res,ceil)
	if mm.Cmp(ceil)>0 {
		res = res.Add(res,big.NewInt(1))
	}
	status:=outputWriter(int(res.Int64()), used)
	os.Exit(ERRORS[status])
}
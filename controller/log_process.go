package log_process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log-process/types"
	"log-process/utils"
	"os"
	"regexp"
	"strconv"
	"time"
)

type LogProcess struct {
	Path        string
	InfluxDBDSN string
	Rc          chan []byte
	Wc          chan []byte
}

func (l *LogProcess) Read() {
	f, err := os.OpenFile(l.Path, os.O_RDONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file failure, %v\n", err))
	}
	defer f.Close()

	// 定位到文件末尾
	f.Seek(0, 2)
	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(time.Second)
			continue
		}
		if err != nil {
			panic(fmt.Sprintf("读取文件失败， %s\n", err))
		}

		l.Rc <- line[:len(line)-1]
	}
}

func (l *LogProcess) LogAnalysis() {
	regStr := `(\d+?\.\d+?\.\d+?\.\d+?)\s-\s-\s\[(\d+)/(\w+)/(\d+):(\d+:\d+:\d+)\s\+\d+\]\s\"(\w+)\s(/.*?)\s.*?\"` +
		`\s(\d+?)\s(.*?)\s(.*?)\s(.*?)\s(\d+)\s(\d+?)\s`
	reg := regexp.MustCompile(regStr)
	if reg == nil {
		panic("reg MustCompile failure")
	}

	for b := range l.Rc {
		str := string(b)
		// 匹配字段
		s := reg.FindStringSubmatch(str)

		// 去除第一个元素
		resultList := s[1:]

		// date
		dd := resultList[1]
		mm := utils.GenerateMonth(resultList[2])
		yy := resultList[3]
		logDate := fmt.Sprintf("%s-%s-%s %s", yy, mm, dd, resultList[4])

		var lf types.LogFormat
		lf.Ip = resultList[0]
		lf.Date = logDate
		lf.Method = resultList[5]
		lf.RequestPath = resultList[6]
		lf.RequestSize, _ = strconv.Atoi(resultList[7])
		lf.UpstreamAddr = resultList[8]
		lf.UpstreamTime, _ = strconv.ParseFloat(resultList[9], 64)
		lf.ResponseTime, _ = strconv.ParseFloat(resultList[10], 64)
		lf.ResponseStatus, _ = strconv.Atoi(resultList[11])
		lf.RequestBodySize, _ = strconv.Atoi(resultList[12])

		// json
		logBin, _ := json.Marshal(lf)

		l.Wc <- logBin
	}
}

func (l *LogProcess) WriteToInfluxDB() {
	for b := range l.Wc {
		// test unmarshal json
		var str types.LogFormat
		fmt.Println(string(b))
		err := json.Unmarshal(b, &str)
		if err != nil {
			panic(fmt.Sprintf("unmarshal failure, %s\n", err))
		}

		fmt.Printf("successfully, %v\n", str)

		// todo: unmarshal json and write to influxdb
	}
}

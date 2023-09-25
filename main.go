package main

import (
	log_process "log-process/controller"
)

func main() {
	var lp = log_process.LogProcess{
		// Path 后续使用flag --log 传入
		//Path:        "D:/software/nginx-1.24.0/logs/access.log", // my windows
		Path:        "/usr/local/nginx/logs/access.log", // my linux
		InfluxDBDSN: "//todo",
		Rc:          make(chan []byte, 10),
		Wc:          make(chan []byte, 10),
	}

	go lp.Read()
	go lp.LogAnalysis()
	go lp.WriteToInfluxDB()

	select {}
}

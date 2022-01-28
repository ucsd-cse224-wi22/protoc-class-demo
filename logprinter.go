package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	
	pb "weblog/weblog"
	"google.golang.org/protobuf/proto"
)

func main() {
	logfp, err := os.Open("log.db")
	if err != nil {
		log.Panic(err)
	}
	defer logfp.Close()

	for {
		loglenbytes := make([]byte, 4)
		_, err := io.ReadFull(logfp, loglenbytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		loglen := binary.BigEndian.Uint32(loglenbytes)

		logbytes := make([]byte, loglen)
		_, err = io.ReadFull(logfp, logbytes)
		if err != nil {
			log.Panic(err)
		}

		logentry := &pb.WebRequest{}
		if err := proto.Unmarshal(logbytes, logentry); err != nil {
			log.Fatal(err)
		}
		log.Println(logentry)
	}
}

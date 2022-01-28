package main

import (
	"encoding/binary"
	"log"
	"net/http"
	"os"
	"time"

	pb "weblog/weblog"
	"google.golang.org/protobuf/proto"
)

func logRequest(logentry *pb.WebRequest) {
	weblogwire, err := proto.Marshal(logentry)
	if err != nil {
		log.Fatal(err)
	}

	logfile, err2 := os.OpenFile("log.db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err2 != nil {
		log.Fatal(err2)
	}

	log.Printf("Writing out a log entry of size: %d\n", uint32(len(weblogwire)))
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(weblogwire)))

	logfile.Write(buf)
	logfile.Write(weblogwire)

	logfile.Close()
}

func WithLogging(h http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
        start := time.Now()

        h.ServeHTTP(rw, r) // serve the original request

		weblog := &pb.WebRequest{Ts: start.Unix(),
									Method: r.Method,
									Url: r.RequestURI,
									Delay: uint64(time.Since(start).Microseconds())}

		logRequest(weblog)

		log.Printf("ts: %v, method: %s, url: %s, duration: %v\n",
			weblog.GetTs(), weblog.GetMethod(), weblog.GetUrl(), weblog.GetDelay())
    })
}

func main() {
	log.Fatal(http.ListenAndServe(":8080",
		WithLogging(http.FileServer(http.Dir("htdocs")))))
}

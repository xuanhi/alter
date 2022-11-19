//File: main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.xuanhi/alter/module"
	"github.xuanhi/alter/utils/zaplog"
)

func timeHandler(format string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
	return http.HandlerFunc(fn)
}

func alterWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			zaplog.Sugar.Errorln("Read failed", err)
		}
		defer r.Body.Close()
		//fmt.Println("json:", string(b))
		notification := &module.Notification{}
		err = json.Unmarshal(b, notification)
		if err != nil {
			zaplog.Sugar.Errorln("json format error:", err)
		}

		zaplog.Sugar.Info("notification:", notification)

	} else {
		zaplog.Sugar.Warnln("ONly support Post")
		fmt.Fprintf(w, "Only support post")
	}
}

func main() {
	zaplog.InitLogger()
	defer zaplog.SyncLogger()
	// Note that we skip creating the ServeMux...

	var format string = time.RFC1123
	th := timeHandler(format)

	// We use http.Handle instead of mux.Handle...
	http.Handle("/time", th)

	http.HandleFunc("/", alterWebhook)

	zaplog.Sugar.Info("Listening...")
	// And pass nil as the handler to ListenAndServe.
	http.ListenAndServe(":5001", nil)
}

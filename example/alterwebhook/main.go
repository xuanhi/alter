// File: main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.xuanhi/alter/module"
	"github.xuanhi/alter/notifier"
	"github.xuanhi/alter/utils/zaplog"
	"go.uber.org/zap"
)

var ctx context.Context
var ChatID string

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
		fmt.Println("json:", string(b))
		notification := &module.Notification{}
		err = json.Unmarshal(b, notification)
		if err != nil {
			zaplog.Sugar.Errorln("json format error:", err)
		}

		zaplog.Sugar.Info("notification:", notification)

		zaplog.Sugar.Infof("当前有%d个告警", len(notification.Alerts))

		for i := 0; i < len(notification.Alerts); i++ {
			text, err := notifier.AlterMsgCard(*notification, i)

			if err != nil {
				zaplog.Sugar.Errorf("消息卡片制作失败", zap.Error(err))
			}
			// cha_id oc_c3c09fe9ac4a8cac3995227476889d5b
			_, err = notifier.SendAlterMsg(ctx, ChatID, text)

			if err != nil {
				zaplog.Sugar.Errorf("发送飞书消息错误", zap.Error(err))
			}
		}

	} else {
		zaplog.Sugar.Warnln("ONly support Post")
		fmt.Fprintf(w, "Only support post")
	}
}

func main() {
	ctx = context.Background()

	zaplog.InitLogger()
	defer zaplog.SyncLogger()
	// Note that we skip creating the ServeMux...
	ChatID = notifier.GetChatIdByfirst(ctx)
	zaplog.Sugar.Infof("获取到群id:%s", ChatID)
	var format string = time.RFC1123
	th := timeHandler(format)

	// We use http.Handle instead of mux.Handle...
	http.Handle("/time", th)

	http.HandleFunc("/", alterWebhook)

	zaplog.Sugar.Info("Listening...")
	// And pass nil as the handler to ListenAndServe.
	http.ListenAndServe(":5001", nil)
}

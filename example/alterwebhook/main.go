// File: main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func alertWebhookReferName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clustername := vars["name"]
	zaplog.Sugar.Infof("当前请求方法: %s", r.Method)
	if r.Method == "POST" {
		b, err := io.ReadAll(r.Body)

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

		//zaplog.Sugar.Info("notification:", notification)

		zaplog.Sugar.Infof("当前有%d个告警", len(notification.Alerts))
		token, err := notifier.GetTenantAccessToken(ctx)
		if err != nil {
			zaplog.Sugar.Errorln("failed to get tenant access token", err)
		}

		if clustername == "all" {
			chatitem, err := notifier.GetChatIdItems(ctx)
			if err != nil {
				zaplog.Sugar.Errorf("获取不到chatiditems:%v", err)
				return
			}
			for i := 0; i < len(notification.Alerts); i++ {
				text, err := notifier.AlterMsgCard(*notification, i)

				if err != nil {
					zaplog.Sugar.Errorf("消息卡片制作失败", zap.Error(err))
				}
				// cha_id oc_c3c09fe9ac4a8cac3995227476889d5b

				notifier.SendAlterMsgWithMutilBytoken(ctx, *chatitem, text, token)

			}
		} else {
			for i := 0; i < len(notification.Alerts); i++ {
				text, err := notifier.AlterMsgCard(*notification, i)

				if err != nil {
					zaplog.Sugar.Errorf("消息卡片制作失败", zap.Error(err))
				}
				// cha_id oc_c3c09fe9ac4a8cac3995227476889d5b

				err = notifier.SendAlterMsgByNameAndtoken(ctx, clustername, text, token)
				if err != nil {
					zaplog.Sugar.Errorf("通过名字没有发现chatid:%v", err)
					break
				}

			}
		}

	} else {
		zaplog.Sugar.Warnln("ONly support Post")
		fmt.Fprintf(w, "Only support post")
	}

}

func alterWebhook(w http.ResponseWriter, r *http.Request) {
	zaplog.Sugar.Infof("当前请求方法: %s", r.Method)
	if r.Method == "POST" {
		b, err := io.ReadAll(r.Body)

		if err != nil {
			zaplog.Sugar.Errorln("Read failed", err)
		}
		defer r.Body.Close()
		zaplog.Sugar.Infoln("获取到json数据: ", string(b))
		//fmt.Println("json:", string(b))
		notification := &module.Notification{}
		err = json.Unmarshal(b, notification)
		if err != nil {
			zaplog.Sugar.Errorln("json format error:", err)
			return
		}

		//zaplog.Sugar.Info("notification:", notification)

		zaplog.Sugar.Infof("当前有%d个告警", len(notification.Alerts))
		token, err := notifier.GetTenantAccessToken(ctx)
		if err != nil {
			zaplog.Sugar.Errorln("failed to get tenant access token", err)
		}

		for i := 0; i < len(notification.Alerts); i++ {
			text, err := notifier.AlterMsgCard(*notification, i)

			if err != nil {
				zaplog.Sugar.Errorf("消息卡片制作失败", zap.Error(err))
			}
			// cha_id oc_c3c09fe9ac4a8cac3995227476889d5b
			_, err = notifier.SendAlterMsgBytoken(ctx, ChatID, text, token)

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
	zaplog.Sugar.Infof("群id默认值:%s", ChatID)
	var format string = time.RFC1123
	th := timeHandler(format)

	r := mux.NewRouter()

	r.Handle("/time", th)
	r.HandleFunc("/", alterWebhook)
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		items, _ := notifier.GetChatIdItems(ctx)

		fmt.Fprintf(w, "items: %v\n", items.Items)
	})
	r.HandleFunc("/cluster/{name}", alertWebhookReferName)
	r.HandleFunc("/webhook/event", notifier.ReceiveEvent(ctx))

	zaplog.Sugar.Info("Listening...")

	http.ListenAndServe(":5001", r)

	// // We use http.Handle instead of mux.Handle...
	// http.Handle("/time", th)

	// http.HandleFunc("/", alterWebhook)

	// zaplog.Sugar.Info("Listening...")
	// // And pass nil as the handler to ListenAndServe.
	// http.ListenAndServe(":5001", nil)
}

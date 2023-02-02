package notifier

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.xuanhi/alter/utils/zaplog"
)

func TestGetTenantAccessToken(t *testing.T) {
	zaplog.InitLogger()
	//	defer zaplog.SyncLogger()
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		zaplog.Sugar.Errorln("failed to get tenant access token", err)
		fmt.Println("failed to get tenant access token", err)
		t.Fail()
		return
	}
	fmt.Println("我的token: ", token)
}

func TestGetChatIdByfirst(t *testing.T) {
	zaplog.InitLogger()
	ctx := context.Background()

	chat_id := GetChatIdByfirst(ctx)

	fmt.Println("我的chat_id:", chat_id)
}

func TestGetChatIdByName(t *testing.T) {
	zaplog.InitLogger()
	ctx := context.Background()
	chat_id, err := GetChatIdByName(ctx, "线上服务器警报群1")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("我的chat_id:", chat_id)
}

func TestSendAlterMsg(t *testing.T) {
	zaplog.InitLogger()
	ctx := context.Background()
	a, err := SendAlterMsg(ctx, "oc_c3c09fe9ac4a8cac3995227476889d5b", "{\"text\": \"xhh test content \"}")
	if err != nil {
		zaplog.Sugar.Errorln("failed ", err)
		t.Fail()
		return
	}
	fmt.Println("test:", a)
}

func TestAlterMsgCard(t *testing.T) {
	//fmt.Println(time.Now())
	//formatTimestr := "02 Jan 2006 15:04:05 GMT"
	formate := "2006-01-02 15:04:05"
	//    DateFmt:="2006-01-02"
	//2022-11-22 05:21:50.432 +0000 UTC
	//reltime, err := time.Parse(formate, "2022-11-22 05:21:50.432")
	loc, _ := time.LoadLocation("Europe/Berlin")
	reltime2, err := time.ParseInLocation(formate, "2022-11-22 05:21:50.432", loc)
	if err != nil {
		fmt.Println("错误", err)
		t.Fail()
	}
	//fmt.Println(reltime.Format(formate))
	fmt.Println(reltime2.Format(formate))
	//start, _ := time.Parse(formate, "2022-11-22 05:21:50.432 +0000 UTC")
	//fmt.Println(start.Add(8 * time.Hour).String())
}

func TestABC(t *testing.T) {
	fmt.Println("==", os.Getenv("AppID"))

}

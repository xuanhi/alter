package notifier

import (
	"context"
	"fmt"
	"testing"

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

	chat_id, err := GetChatIdByfirst(ctx)
	if err != nil {
		zaplog.Sugar.Errorln("failed to get chat_id", err)
		t.Fail()
		return
	}
	fmt.Println("我的chat_id:", chat_id)
}

func TestSendAlterMsg(t *testing.T) {
	zaplog.InitLogger()
	ctx := context.Background()
	a, err := SendAlterMsg(ctx, "oc_c3c09fe9ac4a8cac3995227476889d5b")
	if err != nil {
		zaplog.Sugar.Errorln("failed ", err)
		t.Fail()
		return
	}
	fmt.Println("test:", a)
}

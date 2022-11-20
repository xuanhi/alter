package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.xuanhi/alter/module"
	"github.xuanhi/alter/utils/zaplog"
)

// APIPath
var (
	TenantAccessTokenURL = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	GetChatId            = "https://open.feishu.cn/open-apis/im/v1/chats?page_size=20"
	SendMsgUrl           = "https://open.feishu.cn/open-apis/im/v1/messages"
)

// GetTenantAccessToken get tenant access token for app
// Refer to: https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/tenant_access_token_internal
func GetTenantAccessToken(ctx context.Context) (string, error) {
	cli := &http.Client{}
	reqBody := module.TenantAccessTokenRequest{
		APPID:     "cli_a3f3305dc1f9d013",
		APPSecret: "Ud5TZ38Sh6u1XUBfd4xOhhqClgoFf68l",
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", TenantAccessTokenURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		return "", err
	}
	resp, err := cli.Do(req)
	if err != nil {
		zaplog.Sugar.Errorln("failed to get token:")
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	zaplog.Sugar.Infoln(string(body))

	tokenResp := &module.TenantAccessTokenResponse{}
	err = json.Unmarshal(body, tokenResp)
	if err != nil {
		return "", err
	}
	return tokenResp.TenantAccessToken, nil
}

//获取机器人所在群的chart_id
func GetChatIdByfirst(ctx context.Context) (string, error) {
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		zaplog.Sugar.Errorln("failed to get tenant access token", err)
		return "", err
	}
	cli := &http.Client{}

	req, err := http.NewRequest("GET", GetChatId, nil)
	if err != nil {
		zaplog.Sugar.Errorln("get chatid failed", err)
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := cli.Do(req)
	if err != nil {
		zaplog.Sugar.Errorln("获取群id失败", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zaplog.Sugar.Errorln("读取resp body 失败", err)
		return "", err
	}

	ChatIdDatas := &module.ChatIdDatas{}
	err = json.Unmarshal(body, ChatIdDatas)
	if err != nil {
		zaplog.Sugar.Errorln("json 解析错误", err)
		return "", err
	}
	if ChatIdDatas.Code != 0 {
		zaplog.Sugar.Warnf("没能获取到群id,code:%v,msg:%v", ChatIdDatas.Code, ChatIdDatas.Msg)
	}

	zaplog.Sugar.Infof("获取到chat_id数组,返回第一个chat_id : %v -- 群名字: %v", ChatIdDatas.Data.Items[0].Chat_id, ChatIdDatas.Data.Items[0].Name)
	return ChatIdDatas.Data.Items[0].Chat_id, nil
}

//发送告警消息
func SendAlterMsg(ctx context.Context, chatID string) (*module.MessageItem, error) {
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		zaplog.Sugar.Errorln("failed to get tenant access token", err)
		return nil, err
	}
	cli := &http.Client{}

	MessageReques := struct {
		ReceiveID string `json:"receive_id"`
		Content   string `json:"content"`
		MsgType   string `json:"msg_type"`
	}{
		ReceiveID: chatID,
		Content:   "{\"text\": \"xhh test content \"}",
		MsgType:   "text",
	}

	reqBytes, err := json.Marshal(MessageReques)

	if err != nil {
		zaplog.Sugar.Errorln("json 解析错误", err)
		return nil, err
	}
	req, err := http.NewRequest("POST", SendMsgUrl, strings.NewReader(string(reqBytes)))
	if err != nil {
		zaplog.Sugar.Errorf("new request failed")
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	q := req.URL.Query()
	q.Add("receive_id_type", "chat_id")
	req.URL.RawQuery = q.Encode()

	var logID string
	resp, err := cli.Do(req)

	if err != nil {
		zaplog.Sugar.Errorf("create message failed, err=%v", err)
		return nil, err
	}
	if resp != nil && resp.Header != nil {
		logID = resp.Header.Get("x-tt-logid")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	zaplog.Sugar.Infoln("data: ", string(body))
	if err != nil {
		zaplog.Sugar.Errorln("read body failed")
		return nil, err
	}
	createMessageResp := &module.CreateMessageResponse{}
	err = json.Unmarshal(body, createMessageResp)
	if err != nil {
		zaplog.Sugar.Errorf("failed to unmarshal")
		return nil, err
	}
	if createMessageResp.Code != 0 {
		zaplog.Sugar.Errorf("failed to create message, code: %v, msg: %v, log_id: %v", createMessageResp.Code, createMessageResp.Message, logID)
		return nil, fmt.Errorf("create message failed")
	}
	zaplog.Sugar.Infof("succeed create message, msg_id: %v", createMessageResp.Data.MessageID)
	return createMessageResp.Data, nil
}

package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.xuanhi/alter/module"
	"github.xuanhi/alter/utils/zaplog"
	"go.uber.org/zap"
)

// APIPath
var (
	TenantAccessTokenURL = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	GetChatId            = "https://open.feishu.cn/open-apis/im/v1/chats?page_size=20"
	SendMsgUrl           = "https://open.feishu.cn/open-apis/im/v1/messages"
	AppID                string
	Secret               string
)

// GetTenantAccessToken get tenant access token for app
// Refer to: https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/tenant_access_token_internal
func GetTenantAccessToken(ctx context.Context) (string, error) {
	AppID = os.Getenv("AppID")
	Secret = os.Getenv("Secret")
	//zaplog.Sugar.Infof("系统环境变量%v", os.Environ())
	if AppID == "" || Secret == "" {
		AppID = "cli_a3f3305dc1f9d013"
		Secret = "Ud5TZ38Sh6u1XUBfd4xOhhqClgoFf68l"
		zaplog.Sugar.Infof("使用测试密钥AppID: %s,Secret: %s", AppID, Secret)
	} else {
		zaplog.Sugar.Infof("密钥AppID: %s,Secret: %s", AppID, Secret)
	}

	cli := &http.Client{}
	reqBody := module.TenantAccessTokenRequest{
		APPID:     AppID,
		APPSecret: Secret,
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

// 获取机器人所在群的chart_id,只返回第一个
func GetChatIdByfirst(ctx context.Context) string {
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		zaplog.Sugar.Errorln("failed to get tenant access token", err)
		return ""
	}
	cli := &http.Client{}

	req, err := http.NewRequest("GET", GetChatId, nil)
	if err != nil {
		zaplog.Sugar.Errorln("get chatid failed", err)
		return ""
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := cli.Do(req)
	if err != nil {
		zaplog.Sugar.Errorln("获取群id失败", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zaplog.Sugar.Errorln("读取resp body 失败", err)
		return ""
	}

	ChatIdDatas := &module.ChatIdDatas{}
	err = json.Unmarshal(body, ChatIdDatas)
	if err != nil {
		zaplog.Sugar.Errorln("json 解析错误", err)
		return ""
	}
	if ChatIdDatas.Code != 0 {
		zaplog.Sugar.Warnf("没能获取到群id,code:%v,msg:%v", ChatIdDatas.Code, ChatIdDatas.Msg)
	}

	zaplog.Sugar.Infof("获取到chat_id数组,返回第一个chat_id : %v -- 群名字: %v", ChatIdDatas.Data.Items[0].Chat_id, ChatIdDatas.Data.Items[0].Name)
	return ChatIdDatas.Data.Items[0].Chat_id
}

// 发送告警消息
func SendAlterMsg(ctx context.Context, chatID, altermsg string) (*module.MessageItem, error) {
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
		Content:   altermsg,
		MsgType:   "interactive",
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

// 告警消息卡片制作
func AlterMsgCard(altermsg module.Notification, id int) (string, error) {
	status := "red"
	if altermsg.Alerts[id].Status == "resolved" {
		status = "green"
	}

	headconfig := map[string]bool{
		"wide_screen_mod": true,
	}
	headconfig2 := map[string]interface{}{
		"template": status,
		"title": module.AlterCon{
			Content: fmt.Sprintf("🔺%s  alertname:%s", altermsg.Alerts[id].Status, altermsg.Alerts[id].Labels["alertname"]),
			Tag:     "plain_text",
		},
	}
	// var lablecontent []string
	// for k, v := range altermsg.Alerts[id].Labels {
	// 	lablecontent = append(lablecontent, fmt.Sprintf("%s: %s", k, v))
	// 	lablecontent=fmt.Sprintf("%s: %s\n",k,v)
	// }

	headconfig3 := []map[string]interface{}{
		{
			"tag": "hr",
		}, {
			"tag": "div",
			"text": module.AlterCon{
				Tag: "lark_md",
				Content: fmt.Sprintf("**当前时间：%s**\n**告警类型:** %s\n**告警级别**: %s\n**故障节点:** %s",
					time.Now().Format("2006-01-02 15:04:05"), altermsg.Alerts[id].Labels["alertname"], altermsg.Alerts[id].Labels["severity"], altermsg.Alerts[id].Labels["instance"]),
			},
		}, {
			"tag": "hr",
		}, {
			"tag": "div",
			"text": module.AlterCon{
				Tag: "lark_md",
				//				Content: strings.Join(lablecontent, "\n"),
				Content: fmt.Sprintf("**告警主题: %s**\n\n**告警详情:** %s", altermsg.Alerts[id].Annotations["summary"], altermsg.Alerts[id].Annotations["description"]),
			},
		}, {
			"tag": "hr",
		}, {
			"tag": "div",
			"text": module.AlterCon{
				Tag:     "lark_md",
				Content: fmt.Sprintf("**故障时间:** %s \n**恢复时间:** %s ", altermsg.Alerts[id].StartsAt.Add(8*time.Hour).Format("2006-01-02 15:04:05"), altermsg.Alerts[id].EndsAt.Add(8*time.Hour).Format("2006-01-02 15:04:05")),
			},
		}, {
			"tag": "hr",
		},
	}

	content := module.AlterContent{
		Config:   headconfig,
		Header:   headconfig2,
		Elements: headconfig3,
	}
	contentbyte, err := json.Marshal(content)
	if err != nil {
		zaplog.Sugar.Errorf("json解析错误", zap.Error(err))
		return "", err
	}
	return string(contentbyte), nil

}

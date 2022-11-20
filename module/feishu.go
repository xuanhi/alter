package module

type TenantAccessTokenRequest struct {
	APPID     string `json:"app_id"`
	APPSecret string `json:"app_secret"`
}

type TenantAccessTokenResponse struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}

//群信息
type ChatId struct {
	Chat_id       string `json:"chat_id"`
	Description   string `json:"description"`
	Name          string `json:"name"`
	Owner_id      string `json:"owner_id"`
	Owner_id_type string `json:"owner_id_type"`
}

//群信息items
type ChatIdItems struct {
	Items []ChatId `json:"items"`
}

type ChatIdDatas struct {
	Code int         `json:"code"`
	Data ChatIdItems `json:"data"`
	Msg  string      `json:"msg"`
}

//发送消息post参数
type CreateMessageRequest struct {
	ReceiveID string `json:"receive_id"`
	Content   string `json:"content"`
	MsgType   string `json:"msg_type"`
}

//sendMsg后返回的响应结构体
type CreateMessageResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    *MessageItem `json:"data"`
}

type MessageItem struct {
	MessageID  string         `json:"message_id,omitempty"`
	RootID     string         `json:"root_id,omitempty"`
	ParentID   string         `json:"parent_id,omitempty"`
	MsgType    string         `json:"msg_type,omitempty"`
	CreateTime string         `json:"create_time,omitempty"`
	UpdateTime string         `json:"update_time,omitempty"`
	Deleted    bool           `json:"deleted,omitempty"`
	ChatID     string         `json:"chat_id,omitempty"`
	Sender     *MessageSender `json:"sender,omitempty"`
	Body       *MessageBody   `json:"body,omitempty"`
}

type MessageBody struct {
	Content string `json:"content,omitempty"`
}

type MessageSender struct {
	ID         string `json:"id,omitempty"`
	IDType     string `json:"id_type,omitempty"`
	SenderType string `json:"sender_type"`
	TenantKey  string `json:"tenant_key"`
}

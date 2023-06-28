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

// 群信息
type ChatId struct {
	Chat_id       string `json:"chat_id"`
	Description   string `json:"description"`
	Name          string `json:"name"`
	Owner_id      string `json:"owner_id"`
	Owner_id_type string `json:"owner_id_type"`
}

// 群信息items
type ChatIdItems struct {
	Items []ChatId `json:"items"`
}

type ChatIdDatas struct {
	Code int         `json:"code"`
	Data ChatIdItems `json:"data"`
	Msg  string      `json:"msg"`
}

// 发送消息post参数
type CreateMessageRequest struct {
	ReceiveID string `json:"receive_id"`
	Content   string `json:"content"`
	MsgType   string `json:"msg_type"`
}

// sendMsg后返回的响应结构体
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

//事件API

type ReceiveEventEncrypt struct {
	Encrypt string `json:"encrypt" form:"encrypt"`
}
type DecryptToken struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}
type Event struct {
	Schema string      `json:"schema"`
	Header Header      `json:"header"`
	Event  interface{} `json:"event"`
}
type Header struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	CreateTime string `json:"create_time"`
	Token      string `json:"token"`
	AppID      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
}
type ReceiveMessageEvent struct {
	Schema string       `json:"schema"`
	Header Header       `json:"header"`
	Event  MessageEvent `json:"event"`
}
type MessageEvent struct {
	Sender  Sender  `json:"sender"`
	Message Message `json:"message"`
}
type Sender struct {
	SenderID   map[string]string `json:"sender_id"`
	SenderType string            `json:"sender_type"`
	TenantKey  string            `json:"tenant_key"`
}

type Message struct {
	MessageID   string     `json:"message_id"`
	RootID      string     `json:"root_id"`
	ParentID    string     `json:"parent_id"`
	CreateTime  string     `json:"create_time"`
	ChatID      string     `json:"chat_id"`
	ChatType    string     `json:"chat_type"`
	MessageType string     `json:"message_type"`
	Content     string     `json:"content"`
	Mentions    []*Mention `json:"mentions,omitempty"`
}
type Mention struct {
	Key       string  `json:"key,omitempty"`
	ID        *UserID `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	TenantKey string  `json:"tenant_key,omitempty"`
}
type UserID struct {
	UserID  string `json:"user_id,omitempty"`
	OpenID  string `json:"open_id,omitempty"`
	UnionID string `json:"union_id,omitempty"`
}

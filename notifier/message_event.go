package notifier

import (
	"context"
	"strings"

	"github.xuanhi/alter/module"
	"github.xuanhi/alter/utils/zaplog"
)

func HandleReceiveMessageEvent(ctx context.Context, event *module.ReceiveMessageEvent) error {
	msg := event.Event.Message
	chatID := msg.ChatID
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		zaplog.Sugar.Errorf("failed to get tenant access token :%v", err)
		return err
	}
	switch msg.MessageType {
	case "text":
		if strings.Contains(msg.Content, "/solve") {
			createMsgRequest := &module.CreateMessageRequest{
				ReceiveID: chatID,
				Content:   "{\"text\":\"问题已解决，辛苦了 \\n\"}",
				MsgType:   "text",
			}
			resp, err := SendMessage(ctx, token, createMsgRequest)
			if err != nil {
				zaplog.Sugar.Errorf("failed to send msg")
				return err
			}
			zaplog.Sugar.Infof("succeed send msg, msg_id: %v", resp.MessageID)
		}
	case "post":
		if strings.Contains(msg.Content, "/solve") {
			createMsgRequest := &module.CreateMessageRequest{
				ReceiveID: chatID,
				Content:   "{\"text\":\"问题已解决，辛苦了 \\n\"}",
				MsgType:   "text",
			}
			resp, err := SendMessage(ctx, token, createMsgRequest)
			if err != nil {
				zaplog.Sugar.Errorf("failed to send msg")
				return err
			}
			zaplog.Sugar.Infof("succeed send msg, msg_id: %v", resp.MessageID)
		}

	default:
		zaplog.Sugar.Infof("unhandled message type, msg_type: %v", msg.MessageType)

	}
	return nil

}

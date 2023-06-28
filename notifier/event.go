package notifier

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.xuanhi/alter/module"
	"github.xuanhi/alter/utils/zaplog"
)

func ReceiveEvent(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		zaplog.Sugar.Infof("当前请求方法: %s", r.Method)
		if r.Method == "POST" {
			var req = &module.ReceiveEventEncrypt{}
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				zaplog.Sugar.Errorln("Read failed", err)
			}
			defer r.Body.Close()
			zaplog.Sugar.Infoln("获取到json数据: ", string(b))
			err = json.Unmarshal(b, req)
			if err != nil {
				zaplog.Sugar.Errorf("unmarshal failed %v", err)
				return
			}
			decryptStr, err := Decrypt(req.Encrypt, EncryptKey)
			if err != nil {
				zaplog.Sugar.Errorf("decrypt error: %v", err)
				return
			}
			zaplog.Sugar.Infof("receive decrypt event: %v", decryptStr)
			decryptToken := &module.DecryptToken{}
			err = json.Unmarshal([]byte(decryptStr), decryptToken)
			if err != nil {
				zaplog.Sugar.Errorf("Unmarshal failed again :%v", err)
				return
			}
			if decryptToken.Challenge != "" {
				challenge := struct {
					Challenge string `json:"challenge"`
				}{
					Challenge: decryptToken.Challenge,
				}
				c, err := json.Marshal(challenge)
				if err != nil {
					zaplog.Sugar.Errorf("json 解析失败:%v", err)
				}
				w.WriteHeader(200)
				w.Write(c)
				return
			}
			event := &module.Event{}
			err = json.Unmarshal([]byte(decryptStr), event)
			if err != nil {
				zaplog.Sugar.Errorf("Unmarshal failed, maybe Challenge : %v", err)
				return
			}
			zaplog.Sugar.Infof("receive event, event: %v", event)
			eventType := event.Header.EventType
			zaplog.Sugar.Infof("header: %v", event.Header)
			zaplog.Sugar.Infof("eventType: %v", eventType)
			switch eventType {
			case "im.message.receive_v1":
				receiveMsgEvent := &module.ReceiveMessageEvent{}
				err = json.Unmarshal([]byte(decryptStr), receiveMsgEvent)
				if err != nil {
					zaplog.Sugar.Errorf("Unmarshal failed, maybe Challenge")
					return
				}
				go func() {
					err = HandleReceiveMessageEvent(ctx, receiveMsgEvent)
					if err != nil {
						zaplog.Sugar.Errorf("handle receive message event failed: %v", err)
					}
				}()
			default:
				zaplog.Sugar.Info("unhandled event")
			}
			w.WriteHeader(200)
			respdata := struct {
				Message string `json:"message"`
			}{
				Message: "ok",
			}
			data, err := json.Marshal(respdata)
			if err != nil {
				zaplog.Sugar.Errorf("json 解析失败:%v", err)
			}
			w.Write(data)

		} else {
			zaplog.Sugar.Warnln("ONly support Post")
			fmt.Fprintf(w, "Only support post")
		}
	}
}

func Decrypt(encrypt string, key string) (string, error) {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", fmt.Errorf("base64StdEncode Error[%v]", err)
	}
	if len(buf) < aes.BlockSize {
		return "", errors.New("cipher  too short")
	}
	keyBs := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keyBs[:sha256.Size])
	if err != nil {
		return "", fmt.Errorf("AESNewCipher Error[%v]", err)
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(buf)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return string(buf[n : m+1]), nil
}

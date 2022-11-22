package module

import "time"

type Alert struct {
	Status string `json:"status"`
	//普通标签
	Labels map[string]string `json:"labels"`
	//普通注解
	Annotations map[string]string `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
}

type Notification struct {
	Version  string `json:"version"`
	GroupKey string `json:"groupKey"`
	Status   string `json:"status"`
	Receiver string `json:"receiver"`
	//分组标签
	GroupLabels map[string]string `json:"groupLabels"`
	//普通标签
	CommonLabels map[string]string `json:"commonLabels"`
	//普通注解标签
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

// 飞书通知模板
type AlterContent struct {
	Config   map[string]bool          `json:"config"`
	Header   map[string]interface{}   `json:"header"`
	Elements []map[string]interface{} `json:"elements"`
}

type AlterCon struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type AlterTex struct {
	Tag  string   `json:"tag"`
	Text AlterCon `json:"text"`
}

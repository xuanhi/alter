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

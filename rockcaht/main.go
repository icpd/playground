package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	UserID    = "zT3WtGyaxXr8WQCgw"
	MessageID = "4HNPjWtgt4tweMdfy"
	AuthToken = "9i-g2_wQrqLOIDpbujKyL6bR1HRi_CaV_U4viqvs2h6"
)

var Cookie = fmt.Sprintf("rc_uid=%s; rc_token=%s", UserID, AuthToken)

func main() {
	resp, err := http.Get("https://raw.githubusercontent.com/Rocketchat/Rocket.Chat/7442ffc0226192910affe0ba874e7712cbbe80bc/packages/livechat/src/lib/emoji/emojis.ts")
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return
	}

	re := regexp.MustCompile(`':([^']+):'`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	ticker := time.Tick(5 * time.Second)

	for _, match := range matches[1160:] {
		if len(match) > 1 {
		retry:
			<-ticker
			emoji := match[1]
			fmt.Println(emoji)
			rst := reaction(emoji)

			fmt.Println(rst)
			if strings.Contains(rst, "too many requests") {
				goto retry
			}
		}
	}
}

func reaction(emoji string) string {
	message := fmt.Sprintf(`{"messageId":"%s","reaction":":%s:"}`, MessageID, emoji)
	req, _ := http.NewRequest("POST", "http://chat.test.com/api/v1/chat.react", bytes.NewBuffer([]byte(message)))
	req.Header.Add("Host", "chat.test.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Auth-Token", AuthToken)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Rocket.Chat/4.0.1 Chrome/mChromeVersion Electron/mElectronVersion Safari/mSafariVersion")
	req.Header.Add("X-User-Id", UserID)
	req.Header.Add("Origin", "http://chat.test.com")
	req.Header.Add("Referer", "http://chat.test.com/group/CDN-te4-zhan4-dui4")
	req.Header.Add("Accept-Language", "zh-CN")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", Cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	respBody := string(body)
	return respBody
}

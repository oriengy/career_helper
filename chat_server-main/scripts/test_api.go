package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type client struct {
	baseURL string
	http    *http.Client
}

func newClient(baseURL string) *client {
	return &client{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *client) post(path string, headers map[string]string, body any, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(raw))
	}

	if out == nil {
		return nil
	}
	return json.Unmarshal(raw, out)
}

func main() {
	baseURL := flag.String("base", "http://localhost:8082", "base URL")
	phone := flag.String("phone", "13800000000", "phone number")
	verificationCode := flag.String("code", "1234", "verification code")
	flag.Parse()

	c := newClient(*baseURL)

	baseHeaders := map[string]string{
		"Connect-Protocol-Version": "1",
	}

	fmt.Println("1) PhoneLogin -> token")
	var loginResp struct {
		Token string `json:"token"`
	}
	err := c.post("/user.UserService/PhoneLogin", baseHeaders, map[string]any{
		"phone":            *phone,
		"verificationCode": *verificationCode,
	}, &loginResp)
	if err != nil {
		fmt.Println("PhoneLogin failed:", err)
		os.Exit(1)
	}
	if loginResp.Token == "" {
		fmt.Println("PhoneLogin failed: token is empty")
		os.Exit(1)
	}
	fmt.Println("   token ok")

	authHeaders := map[string]string{
		"Connect-Protocol-Version": "1",
		"Authorization":            "Bearer " + loginResp.Token,
	}

	fmt.Println("2) GetConfig (public)")
	var cfgResp any
	err = c.post("/config.ConfigService/GetConfig", baseHeaders, map[string]any{
		"keys":     []string{},
		"app":      "",
		"platform": "",
		"env":      "",
		"version":  "",
	}, &cfgResp)
	if err != nil {
		fmt.Println("GetConfig failed:", err)
		os.Exit(1)
	}
	fmt.Println("   config ok")

	fmt.Println("3) CreateChatSession (auth)")
	var chatResp struct {
		ChatSession struct {
			ID string `json:"id"`
		} `json:"chatSession"`
	}
	err = c.post("/chat.ChatService/CreateChatSession", authHeaders, map[string]any{
		"profile": map[string]any{
			"name":   "TestUser",
			"imName": "TestIM",
			"avatar": "",
			"gender": "female",
		},
	}, &chatResp)
	if err != nil {
		fmt.Println("CreateChatSession failed:", err)
		os.Exit(1)
	}
	if chatResp.ChatSession.ID == "" {
		fmt.Println("CreateChatSession failed: session id is empty")
		os.Exit(1)
	}
	sessionID := chatResp.ChatSession.ID
	fmt.Println("   sessionId:", sessionID)

	fmt.Println("4) CreateChatMessage (auth)")
	var msgResp struct {
		Messages []any `json:"messages"`
	}
	err = c.post("/message.ChatMessageService/CreateChatMessage", authHeaders, map[string]any{
		"messages": []map[string]any{
			{
				"sessionId": sessionID,
				"role":      "SELF",
				"msgType":   "HISTORY",
				"content":   "Hello",
				"tags":      []string{},
			},
			{
				"sessionId": sessionID,
				"role":      "FRIEND",
				"msgType":   "HISTORY",
				"content":   "Hi there",
				"tags":      []string{},
			},
		},
	}, &msgResp)
	if err != nil {
		fmt.Println("CreateChatMessage failed:", err)
		os.Exit(1)
	}
	fmt.Println("   messages created:", len(msgResp.Messages))

	fmt.Println("5) ListChatMessages (auth)")
	var listResp struct {
		Messages []any `json:"messages"`
	}
	err = c.post("/message.ChatMessageService/ListChatMessages", authHeaders, map[string]any{
		"sessionId": sessionID,
		"pageSize":  50,
	}, &listResp)
	if err != nil {
		fmt.Println("ListChatMessages failed:", err)
		os.Exit(1)
	}
	fmt.Println("   messages listed:", len(listResp.Messages))

	fmt.Println("6) ListProfiles (auth)")
	var profResp struct {
		Profiles []any `json:"profiles"`
	}
	err = c.post("/profile.ProfileService/ListProfiles", authHeaders, map[string]any{
		"searchName": "",
		"ids":        []string{},
		"pageSize":   "20",
	}, &profResp)
	if err != nil {
		fmt.Println("ListProfiles failed:", err)
		os.Exit(1)
	}
	fmt.Println("   profiles listed:", len(profResp.Profiles))

	fmt.Println("7) Translate (auth, optional - requires AI keys)")
	var transResp struct {
		Content string `json:"content"`
	}
	err = c.post("/translate.TranslateService/Translate", authHeaders, map[string]any{
		"content": "你好",
		"from":    "FEMALE",
		"to":      "MALE",
		"history": "用户: 你好\n朋友: 在忙吗",
	}, &transResp)
	if err != nil {
		fmt.Println("   translate skipped/fail:", err)
	} else {
		fmt.Println("   translate ok:", transResp.Content)
	}

	fmt.Println("8) SendConsultMessage (auth, optional - requires AI keys)")
	var consultResp any
	err = c.post("/message.ChatMessageService/SendConsultMessage", authHeaders, map[string]any{
		"sessionId": sessionID,
		"content":   "我该怎么回复对方？",
	}, &consultResp)
	if err != nil {
		fmt.Println("   consult skipped/fail:", err)
	} else {
		fmt.Println("   consult ok")
	}

	fmt.Println("Done.")
}

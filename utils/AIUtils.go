package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"yatori-go-coreV3/models/ctype"
)

// AIChatMessages ChatGLMChat struct that holds the chat messages.
type AIChatMessages struct {
	Messages []Message `json:"messages"`
}

// Message struct represents individual messages.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AiMut AI锁，防止同时过多调用
var AiMut sync.Mutex

// AggregationAIApi 聚合所有AI接口，直接通过aiType判断然后返回内容
func AggregationAIApi(url,
	model string,
	aiType ctype.AiType,
	aiChatMessages AIChatMessages,
	apiKey string) (string, error) {
	AiMut.Lock()
	defer AiMut.Unlock()
	switch aiType {
	case ctype.ChatGLM:
		return ChatGLMChatReplyApi(model, apiKey, aiChatMessages, 3, nil)
	case ctype.XingHuo:
		return XingHuoChatReplyApi(model, apiKey, aiChatMessages, 3, nil)
	case ctype.TongYi:
		return TongYiChatReplyApi(model, apiKey, aiChatMessages, 3, nil)
	case ctype.DouBao:
		return DouBaoChatReplyApi(model, apiKey, aiChatMessages, 3, nil)
	case ctype.Other:
		return OtherChatReplyApi(url, model, apiKey, aiChatMessages, 3, nil)
	default:
		return "", errors.New(string("AI Type: " + aiType))
	}
}

// AICheck AI可用性检测
func AICheck(url, model, apiKey string, aiType ctype.AiType) error {
	AiMut.Lock()
	defer AiMut.Unlock()
	aiChatMessages := AIChatMessages{
		Messages: []Message{
			{
				Role:    "user",
				Content: "你好",
			},
		},
	}

	if aiType == "" {
		return errors.New("AI Type: " + "请先填写AIType参数")
	}
	if apiKey == "" {
		return errors.New("无效apiKey，请检查apiKey是否正确填写")
	}

	switch aiType {
	case ctype.ChatGLM:
		_, err := ChatGLMChatReplyApi(model, apiKey, aiChatMessages, 3, nil)
		return err
	case ctype.XingHuo:
		_, err := XingHuoChatReplyApi(model, apiKey, aiChatMessages, 3, nil)
		return err
	case ctype.TongYi:
		_, err := TongYiChatReplyApi(model, apiKey, aiChatMessages, 3, nil)
		return err
	case ctype.DouBao:
		_, err := DouBaoChatReplyApi(model, apiKey, aiChatMessages, 3, nil)
		return err
	case ctype.Other:
		_, err := OtherChatReplyApi(url, model, apiKey, aiChatMessages, 3, nil)
		return err
	default:
		return errors.New(string("AI Type: " + aiType))
	}
}

// TongYiChatReplyApi 通义千问API
func TongYiChatReplyApi(
	model,
	apiKey string,
	aiChatMessages AIChatMessages,
	retryNum int, /*最大重连次数*/
	lastErr error,
) (string, error) {
	if retryNum < 0 { //重连次数用完直接返回
		return "", lastErr
	}
	if model == "" {
		model = "qwen-plus"
	}
	client := &http.Client{
		Timeout: 30 * time.Second, // Set connection and read timeout
	}

	url := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	requestBody := map[string]interface{}{
		"model":    model,
		"messages": aiChatMessages.Messages,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return TongYiChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to execute HTTP request: %v", err))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return TongYiChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to read response body: %v", err))
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		time.Sleep(100 * time.Millisecond)
		return TongYiChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to parse JSON response: %v  response body: %s", err, body))
	}

	choices, ok := responseMap["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("AI回复内容未找到，AI返回信息：" + string(body))
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to parse message from response")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content field missing or not a string in response")
	}

	return content, nil
}

// ChatGLM API
func ChatGLMChatReplyApi(
	model,
	apiKey string,
	aiChatMessages AIChatMessages,
	retryNum int, /*最大重连次数*/
	lastErr error,
) (string, error) {
	if model == "" {
		model = "glm-4"
	}
	if retryNum < 0 { //重连次数用完直接返回
		return "", lastErr
	}
	client := &http.Client{
		Timeout: 30 * time.Second, // Set connection and read timeout
	}

	url := "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	requestBody := map[string]interface{}{
		"model":    model,
		"messages": aiChatMessages.Messages,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return ChatGLMChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to execute HTTP request: %v", err))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return ChatGLMChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to read response body: %v", err))
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		time.Sleep(100 * time.Millisecond)
		return ChatGLMChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to parse JSON response: %v   response body: %s", err, body))
	}

	choices, ok := responseMap["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("AI回复内容未找到，AI返回信息：" + string(body))
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to parse message from response")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content field missing or not a string in response")
	}

	return content, nil
}

// 星火API
func XingHuoChatReplyApi(model,
	apiKey string,
	aiChatMessages AIChatMessages,
	retryNum int, /*最大重连次数*/
	lastErr error,
) (string, error) {
	if retryNum < 0 { //重连次数用完直接返回
		return "", lastErr
	}
	if model == "" {
		model = "generalv3.5" //默认模型
	}
	client := &http.Client{
		Timeout: 30 * time.Second, // Set connection and read timeout
	}

	url := "https://spark-api-open.xf-yun.com/v1/chat/completions"
	requestBody := map[string]interface{}{
		"model":    model,
		"messages": aiChatMessages.Messages,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return XingHuoChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to marshal JSON data: %v", err))
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return XingHuoChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to create HTTP request: %v", err))

	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return XingHuoChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to execute HTTP request: %v", err))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return XingHuoChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, err)
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		time.Sleep(100 * time.Millisecond)
		return XingHuoChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to parse JSON response: %v   response body: %s", err, body))
	}

	choices, ok := responseMap["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		//防止傻逼星火认为频繁调用报错的问题，踏马老子都加锁了还频繁调用，我频繁密码了
		if strings.Contains(responseMap["error"].(map[string]interface{})["message"].(string), "AppIdQpsOverFlow") {
			time.Sleep(100 * time.Millisecond)
			return XingHuoChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, err)
		}
		log.Printf("unexpected response structure: %v", responseMap)
		return "", fmt.Errorf("AI回复内容未找到，AI返回信息：" + string(body))
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to parse message from response")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content field missing or not a string in response")
	}

	return content, nil
}

// DouBaoChatReplyApi 豆包API
func DouBaoChatReplyApi(model,
	apiKey string,
	aiChatMessages AIChatMessages,
	retryNum int, /*最大重连次数*/
	lastErr error,
) (string, error) {
	if retryNum < 0 { //重连次数用完直接返回
		return "", lastErr
	}
	client := &http.Client{
		Timeout: 40 * time.Second, // Set connection and read timeout
	}

	url := "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
	requestBody := map[string]interface{}{
		"model":    model,
		"messages": aiChatMessages.Messages,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return DouBaoChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to execute HTTP request: %v", err))

	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return DouBaoChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to read response body: %v", err))
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		time.Sleep(100 * time.Millisecond)
		return DouBaoChatReplyApi(model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to parse JSON response: %v    response body: %s", err, body))
	}

	choices, ok := responseMap["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Printf("unexpected response structure: %v", responseMap)
		return "", fmt.Errorf("AI回复内容未找到，AI返回信息：" + string(body))
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to parse message from response")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content field missing or not a string in response")
	}

	return content, nil
}

// OtherChatReplyApi 其他支持CHATGPT API格式的AI模型接入
func OtherChatReplyApi(url,
	model,
	apiKey string,
	aiChatMessages AIChatMessages,
	retryNum int, /*最大重连次数*/
	lastErr error,
) (string, error) {
	if retryNum < 0 { //重连次数用完直接返回
		return "", lastErr
	}
	client := &http.Client{
		Timeout: 40 * time.Second, // Set connection and read timeout
	}
	requestBody := map[string]interface{}{
		"model":    model,
		"messages": aiChatMessages.Messages,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return OtherChatReplyApi(url, model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to execute HTTP request: %v", err))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		return OtherChatReplyApi(url, model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to read response body: %v", err))
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		time.Sleep(100 * time.Millisecond)
		return OtherChatReplyApi(url, model, apiKey, aiChatMessages, retryNum-1, fmt.Errorf("failed to parse JSON response: %v    response body: %s", err, body))
	}

	choices, ok := responseMap["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("AI回复内容未找到，AI返回信息：" + string(body))
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to parse message from response")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content field missing or not a string in response")
	}

	return content, nil
}

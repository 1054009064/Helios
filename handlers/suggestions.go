package handlers

import (
	"context"
	"fmt"
	"helios/config"
	"helios/lib"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
)

type Suggestion struct {
	Text string `json:"text"`
}

type SuggestionsResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
}

func SuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")

	// 如果没有查询参数，返回空建议
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		response := SuggestionsResponse{
			Suggestions: []Suggestion{},
		}
		sonic.ConfigDefault.NewEncoder(w).Encode(response)
		return
	}

	// 获取第一个 API 站点
	apiSites := config.GlobalConfig.SiteList
	if len(apiSites) == 0 {
		w.Header().Set("Content-Type", "application/json")
		response := SuggestionsResponse{
			Suggestions: []Suggestion{},
		}
		sonic.ConfigDefault.NewEncoder(w).Encode(response)
		return
	}

	firstSite := apiSites[0]

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 使用 channel 来接收结果或超时
	resultChan := make(chan []Suggestion, 1)

	go func() {
		// 从第一个源搜索
		results, err := lib.SearchFromAPI(firstSite, query, 1)
		if err != nil {
			resultChan <- []Suggestion{}
			return
		}

		// 应用黄色过滤
		results = lib.FilterYellowContent(results)

		// 转换为建议格式
		suggestions := make([]Suggestion, 0, len(results))
		for _, result := range results {
			suggestions = append(suggestions, Suggestion{
				Text: result.Title,
			})
		}

		resultChan <- suggestions
	}()

	var suggestions []Suggestion
	select {
	case suggestions = <-resultChan:
		// 成功获取结果
	case <-ctx.Done():
		// 超时，返回空结果
		fmt.Printf("搜索建议超时\n")
		suggestions = []Suggestion{}
	}

	w.Header().Set("Content-Type", "application/json")
	response := SuggestionsResponse{
		Suggestions: suggestions,
	}
	sonic.ConfigDefault.NewEncoder(w).Encode(response)
}

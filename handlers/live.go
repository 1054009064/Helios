package handlers

import (
	"helios/config"
	"net/http"

	"github.com/bytedance/sonic"
)

type LiveSourceResponse struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	EPG      string `json:"epg"`
	UA       string `json:"ua,omitempty"`
	From     string `json:"from"`
	Disabled bool   `json:"disabled"`
}

type LiveSourcesAPIResponse struct {
	Success bool                 `json:"success"`
	Data    []LiveSourceResponse `json:"data"`
}

func LiveSourcesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sources := make([]LiveSourceResponse, 0, len(config.GlobalConfig.Lives))
	for _, live := range config.GlobalConfig.Lives {
		sources = append(sources, LiveSourceResponse{
			Key:      live.Key,
			Name:     live.Name,
			URL:      live.Url,
			EPG:      live.Epg,
			UA:       live.UA,
			From:     "config", // 来源为配置文件
			Disabled: false,    // 默认不禁用
		})
	}

	response := LiveSourcesAPIResponse{
		Success: true,
		Data:    sources,
	}

	w.Header().Set("Content-Type", "application/json")
	sonic.ConfigDefault.NewEncoder(w).Encode(response)
}

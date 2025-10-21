package handlers

import (
	"helios/config"
	"net/http"

	"github.com/bytedance/sonic"
)

func ResourcesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 将 map 转换为 list
	apiSites := config.GlobalConfig.APISites
	resources := make([]interface{}, 0, len(apiSites))
	
	for _, site := range apiSites {
		resources = append(resources, site)
	}

	w.Header().Set("Content-Type", "application/json")
	sonic.ConfigDefault.NewEncoder(w).Encode(resources)
}

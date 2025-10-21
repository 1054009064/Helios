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

	w.Header().Set("Content-Type", "application/json")
	sonic.ConfigDefault.NewEncoder(w).Encode(config.GlobalConfig.SiteList)
}

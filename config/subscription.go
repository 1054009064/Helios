package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"helios/models"

	"github.com/btcsuite/btcutil/base58"
)

var GlobalConfig models.Config

func FetchSubscription(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch subscription: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("subscription request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	decodedData := base58.Decode(string(body))
	if len(decodedData) == 0 {
		return fmt.Errorf("failed to decode base58 data")
	}

	// 先用 encoding/json 解析以保持键的顺序
	var rawConfig struct {
		APISites json.RawMessage `json:"api_site"`
		Lives    json.RawMessage `json:"lives"`
	}

	if err := json.Unmarshal(decodedData, &rawConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config JSON: %v", err)
	}

	// 构建最终的 config
	var config models.Config
	config.APISites = make(map[string]models.APISite)
	config.SiteList = make([]models.APISite, 0)
	config.Lives = make([]models.LiveSource, 0)

	// 使用 json.Decoder 按顺序解析 api_site 对象
	dec := json.NewDecoder(bytes.NewReader(rawConfig.APISites))

	// 读取开始的 {
	if _, err := dec.Token(); err != nil {
		return fmt.Errorf("failed to parse api_site: %v", err)
	}

	// 按顺序读取每个键值对
	for dec.More() {
		// 读取键
		token, err := dec.Token()
		if err != nil {
			return fmt.Errorf("failed to read key: %v", err)
		}
		key := token.(string)

		// 读取值
		var site models.APISite
		if err := dec.Decode(&site); err != nil {
			return fmt.Errorf("failed to decode site: %v", err)
		}

		site.Key = key
		config.APISites[key] = site
		config.SiteList = append(config.SiteList, site)
	}

	if len(config.APISites) == 0 {
		return fmt.Errorf("no API sites found in subscription")
	}

	// 解析 lives 字段（如果存在）
	if len(rawConfig.Lives) > 0 {
		liveDec := json.NewDecoder(bytes.NewReader(rawConfig.Lives))

		// 读取开始的 {
		if _, err := liveDec.Token(); err != nil {
			return fmt.Errorf("failed to parse lives: %v", err)
		}

		// 按顺序读取每个键值对
		for liveDec.More() {
			// 读取键
			token, err := liveDec.Token()
			if err != nil {
				return fmt.Errorf("failed to read live key: %v", err)
			}
			key := token.(string)

			// 读取值
			var live models.LiveSource
			if err := liveDec.Decode(&live); err != nil {
				return fmt.Errorf("failed to decode live source: %v", err)
			}

			live.Key = key
			config.Lives = append(config.Lives, live)
		}

		log.Printf("Loaded %d live sources", len(config.Lives))
	}

	GlobalConfig = config
	log.Printf("Subscription config loaded successfully. API sites: %d, Live sources: %d", len(config.APISites), len(config.Lives))
	return nil
}

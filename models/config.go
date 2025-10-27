package models

type Config struct {
	APISites map[string]APISite `json:"api_site"`
	SiteList []APISite          `json:"site_list"`
	Lives    []LiveSource       `json:"lives"`
}

type APISite struct {
	Key    string `json:"key"`
	API    string `json:"api"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

type LiveSource struct {
	Key  string `json:"key"`
	Url  string `json:"url"`
	Name string `json:"name"`
	UA   string `json:"ua"`
	Epg  string `json:"epg"`
}

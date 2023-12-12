package main

type Config struct {
	DomainName []string             `yaml:"domainName"`
	Domains    map[string]int64     `yaml:"Domains"`
	GroupID    GroupID              `yaml:"groupID"`
	BotToken   string               `yaml:"botToken"`
	GroupAuth  map[string]GroupAuth `yaml:"groupAuth"`
	Redis      Redis                `yaml:"redis"`
}

type GroupID struct {
	AdminGroupID int64 `yaml:"adminGroupID"` //后台控制群
	UserGroupID  int64 `yaml:"userGroupID"`  //用户发码群
}

type GroupAuth struct {
	ID     int64  `yaml:"id"`
	AuthID string `yaml:"authID"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Model struct {
	H5               []string `yaml:"H5"`
	Admin            []string `yaml:"后台"`
	Agent            []string `yaml:"代理"`
	App              []string `yaml:"app"`
	DownloadUrl      []string `yaml:"app下载地址"`
	SpareDownloadUrl []string `yaml:"app备用下载地址"`
	Service          string   `yaml:"美洽客服"`
	ServiceAccount   string   `yaml:"美洽账号"`
	ServiceUrl       string   `yaml:"客服链接"`
	SpareH5          []string `yaml:"备用域名"`
	BlockDomain      []string `yaml:"空白域名"`
}

const (
	ICEX = 1 << iota
	M1F
	MIAX
	TGX
	VGX
	ISE
	BitBank
	SZ
	Shop
	Aquis
	Voya
	JinSha
	ShangPuJing
	Jason
)

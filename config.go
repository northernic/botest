package main

type Config struct {
	DomainName  []string         `yaml:"domainName"`
	Domains     map[string]int64 `yaml:"Domains"`
	GroupID     int64            `yaml:"groupID"`
	BotToken    string           `yaml:"botToken"`
	ShangPuJing string           `yaml:"ShangPuJing"`
	JinSha      string           `yaml:"JinSha"`

	ICEX         Model                `yaml:"欧美ICEX"`
	M1F          Model                `yaml:"欧美M1F"`
	LSEX         Model                `yaml:"欧美LSEX"`
	MIAX         Model                `yaml:"欧美MIAX"`
	TGX          Model                `yaml:"欧美TGX"`
	VGX          Model                `yaml:"欧美VGX"`
	ISE          Model                `yaml:"欧美ISE"`
	BitBank      Model                `yaml:"比特银行"`
	Shop         Model                `yaml:"跨境电商"`
	JinSha1      Model                `yaml:"2.1金沙项目"`
	ShangPuJing1 Model                `yaml:"2.1上普京项目"`
	Voya         Model                `yaml:"voya"`
	Aquis        Model                `yaml:"Aquis"`
	Jason        Model                `yaml:"律师事务所"`
	GroupAuth    map[string]GroupAuth `yaml:"groupAuth"`
}

type GroupAuth struct {
	ID     int64  `yaml:"id"`
	AuthID string `yaml:"authID"`
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

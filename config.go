package main

type Config struct {
	DomainName []string `yaml:"domainName"`
	GroupID    int64    `yaml:"groupID"`
	BotToken   string   `yaml:"botToken"`
	Alternate  struct {
		M1F struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"M1F"`
		ISE struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"ISE"`
		ICEX struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"ICEX"`
		MIAX struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"MIAX"`
		TGX struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"TGX"`
		VGX struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"VGX"`
		Szzg struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"Szzg"`
		Shop struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"Shop"`
		JinSha struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"JinSha"`
		ShangPuJing struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"ShangPuJing"`
		LSEX struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
			DomainName    []string `yaml:"domainName"`
		} `yaml:"LSEX"`
	} `json:"alternate"`

	ShangPuJing string `yaml:"ShangPuJing"`

	JinSha string `yaml:"JinSha"`

	ICEX         ICEX        `yaml:"欧美ICEX"`
	M1F          M1F         `yaml:"欧美M1F"`
	LSEX         LSEX        `yaml:"欧美LSEX"`
	MIAX         MIAX        `yaml:"欧美MIAX"`
	TGX          TGX         `yaml:"欧美TGX"`
	VGX          VGX         `yaml:"欧美VGX"`
	ISE          ISE         `yaml:"欧美ISE"`
	BitBank      BitBank     `yaml:"比特银行"`
	SZ           SZ          `yaml:"数字中国"`
	Shop         Shop        `yaml:"跨境电商"`
	JinSha1      JinSha      `yaml:"2.1金沙项目"`
	ShangPuJing1 ShangPuJing `yaml:"2.1上普京项目"`
	LuHai        LuHai       `yaml:"陆海新通道"`
}

type ICEX struct {
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
}
type M1F struct {
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
}
type LSEX struct {
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
}
type MIAX struct {
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
}
type TGX struct {
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
}
type VGX struct {
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
}

type ISE struct {
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
}

type BitBank struct {
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
}

type SZ struct {
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
}

type Shop struct {
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
}

type JinSha struct {
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
}

type ShangPuJing struct {
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
}

type LuHai struct {
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
}

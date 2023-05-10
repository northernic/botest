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

	//ShangPuJing struct {
	//	Infos []struct {
	//		ID        int      `yaml:"id"`
	//		AdminName []string `yaml:"adminName"`
	//	} `yaml:"infos"`
	//	BeiYong []string `yaml:"beiYong"`
	//} `yaml:"ShangPuJing"`

	ShangPuJing string `yaml:"ShangPuJing"`

	JinSha string `yaml:"JinSha"`
}

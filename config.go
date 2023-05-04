package main

type Config struct {
	DomainName []string `yaml:"domainName"`
	GroupID    int64    `yaml:"groupID"`
	BotToken   string   `yaml:"botToken"`
	Alternate  struct {
		M1F struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
		} `yaml:"M1F"`
		ISE struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
		} `yaml:"ISE"`
		ICEX struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
		} `yaml:"ICEX"`
		MIAX struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
		} `yaml:"MIAX"`
		TGX struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
		} `yaml:"TGX"`
		VGX struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
		} `yaml:"VGX"`
		Szzg struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
		} `yaml:"Szzg"`
		Shop struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
		} `yaml:"Shop"`
		JinSha struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
		} `yaml:"JinSha"`
		ShangPuJing struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
		} `yaml:"ShangPuJing"`
		LSEX struct {
			Name          string   `yaml:"name"`
			NewDomainName []string `yaml:"newDomainName"`
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

	JinSha struct {
		Infos []struct {
			ID        int      `yaml:"id"`
			AdminName []string `yaml:"adminName"`
		} `yaml:"infos"`
		BeiYong []string `yaml:"beiYong"`
	} `yaml:"JinSha"`
}

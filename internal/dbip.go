package internal

import (
	"fmt"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

var ispName = map[string]string{
	"SeeHu Technology":          "视虎科技",
	"ChinaCache":                "蓝汛",
	"west.cn":                   "西部数码",
	"Shaanxi BC & TV Network":   "陕西广电",
	"Jinke Information Network": "金科信息网",
	"Tokyo Metropolis CDS International Interconnection Nodes": "东京都北京首都在线(CDS)国际互联节点",
	"Weiyi Network":          "唯一网络",
	"New Generation Network": "新一代",
	"Beijing ZhongXinHengYe": "北京众鑫恒业科技有限责任公司",
	"SITA":                   "国际航空电讯集团公司(SITA)",
	"BGP multi line":         "BGP多线",
	"Anchang Network":        "安畅网络",
	"Yingtong Network":       "盈通网络",
	"Baidu Cloud":            "百度云",
	"Great Wall Broadband":   "长城宽带",
	"China Unicom":           "联通",
	"Y-Link Network":         "云林网络",
	"China Mobile":           "移动",
	"Backbone Network":       "骨干网",
	"Cnix":                   "皓宽网络",
	"SPEEDTOP":               "速必拓网络科技有限公司",
	"KUANCOM":                "京宽网络",
	"Yan Da Zheng Yang":      "燕大正洋",
	"Wasu":                   "华数",
	"Aiwei Broadband":        "爱为宽带",
	"Guangzhou Shanghang Information Technology CO., Ltd": "广州尚航信息科技有限公司",
	"SIEMENS":                         "西门子公司",
	"Tencent Cloud":                   "腾讯云",
	"Haoyun telecom":                  "浩云电信",
	"FIBRLINK":                        "中电飞华",
	"Dalian University of Technology": "大连理工大学",
	"Yungu Technology":                "夽谷科技",
	"Telstra":                         "太平洋电信",
	"Shanda Group":                    "盛大网络",
	"BGP data center":                 "BGP数据中心",
	"UCloud":                          "优刻云",
	"NOVA net":                        "南凌科技",
	"Aliyun":                          "阿里云",
	"Bangrun Technology":              "邦润科技",
	"jdcloud":                         "京东云",
	"GAINET":                          "景安",
	"Hewlett-Packard":                 "惠普",
	"Zhujiang Broadband":              "珠江宽频",
	"Kingsoft Cloud":                  "金山云",
	"China Broadnet":                  "广电",
	"Wangsu":                          "网宿",
	"cloudvsp":                        "天地祥云",
	"Shuhuashi Technology":            "舒华士",
	"AnLai Communication":             "安莱信息通信",
	"Quanjie Technology":              "全捷科技发展有限公司",
	"Baidu Spider":                    "百度蜘蛛",
	"PubYun":                          "公云PubYun",
	"CERNET":                          "教育网",
	"Youtian Broadband":               "油田宽带",
	"Cisco":                           "思科",
	"California":                      "加州",
	"San Xin Shi Dai":                 "三信时代",
	"niaoyun":                         "小鸟云",
	"Alibaba":                         "阿里巴巴",
	"Oray dandelion":                  "蒲公英",
	"Aadata":                          "钜讯网络",
	"sina":                            "新浪",
	"Ocn":                             "东方有线",
	"Ningbo Gao Fang":                 "宁波高防",
	"He Nei Shi Dai":                  "广东省河内时代网络科技有限公司深圳分公司",
	"Stdaily":                         "科技网",
	"Shitong Broadband":               "视通宽带",
	"Department of Education of Anhui Province": "安徽省教育厅",
	"Base Station WiFi":                         "基站WiFi",
	"HaoYun":                                    "浩云",
	"SUNING":                                    "苏宁",
	"National Node":                             "国际节点",
	"Tong Mei Group":                            "同煤集团",
	"topway":                                    "天威视讯",
	"Woodnn":                                    "网鼎科技",
	"Shixun Broadband":                          "视讯宽带",
	"Spider":                                    "蜘蛛",
	"CHINACOMM":                                 "中电华通",
	"Zonergy":                                   "中兴能源",
	"NETEASE":                                   "网易",
	"Beijing Wang Yun Wu Xian Technology":       "北京网云无限科技有限公司BGP多线数据中心",
	"Heng Hui Technology":                       "恒慧通信",
	"Xi'an Jiaotong University":                 "西安交通大学",
	"KJNet":                                     "宽捷网络",
	"PCCW":                                      "电讯盈科",
	"Yangchen Weiye Technology CO., Ltd":        "阳晨伟业科技有限公司",
	"Guangdong Wang Cheng CO., Ltd":             "广东网城在线有限公司",
	"Tian Ying Information Technology":          "天盈信息技术",
	"Xiang Da Xin":                              "祥达信",
	"SINNET":                                    "光环新网",
	"Wensu Network":                             "稳速网络",
	"IDCS":                                      "天互数据电信",
	"NET263":                                    "263网络通信",
	"Beijing Teletron":                          "电信通",
	"Microsoft Cloud":                           "微软云",
	"Cable Network":                             "有线",
	"Vnet":                                      "世纪互联",
	"ground telecom":                            "润迅通信",
	"Hutchison Whampoa Limited (HWL)":           "和记黄埔",
	"China FAW Group Corporation":               "一汽",
	"NETEASE Cloud":                             "网易云",
	"Huatong Broadband":                         "华通宽带",
	"Beijing AiDi Communicate":                  "北京爱迪通信科技有限公司",
	"China Pingmei Shenma":                      "平煤神马集团",
	"Netbank":                                   "网银互联",
	"China Telecom":                             "电信",
	"Ruijiang Technology":                       "睿江科技",
	"HUAWEI Cloud":                              "华为云",
	"Beijing Gehua CATV Network Co., Ltd":       "歌华有线",
	"Link-Net telecom":                          "临网通讯",
	"Liaoning Fangyi Tech CO., Ltd":             "辽宁方翊科技有限公司",
	"CATV":                                      "有线电视",
	"Xin Fei Jin Xin":                           "新飞金信",
	"Linux Pathshala":                           "Linux Pathshala数据中心",
	"QingCloud":                                 "青云QingCloud",
	"Unknown":                                   "未知",
	"China Railcom":                             "中移铁通",
	"Bolu Telecom":                              "博路电信",
	"Google Cloud":                              "谷歌云",
	"KINPONET":                                  "KINPONET软银数据中心",
	"KNET":                                      "北龙中网",
	"Baidu":                                     "百度",
	"Founder Broadband":                         "方正宽带",
	"LanDui Network":                            "蓝队网络",
	"Abbott Laboratories":                       "北芝加哥雅培公司",
	"CNISP":                                     "互联网服务商联盟",
	"Huayu Broadband":                           "华宇宽带",
	"Ping'an Technology":                        "平安科技",
	"MOS":                                       "美团云",
	"Hong Kong Broadband Network":               "香港宽频",
	"Kaopu Cloud":                               "靠谱云",
	"BIH":                                       "互联港湾",
	"Topnew Info":                               "铜牛",
	"Juyou Network":                             "聚友网络",
	"Amazon Cloud":                              "亚马逊云",
	"Weisai Network":                            "维赛网络",
	"Dr. Peng Telecom and Media Group":          "鹏博士",
}

func NewDB(mmdb string) (*geoip2.Reader, error) {
	return geoip2.Open(mmdb)
}

func GetIPInfo(targetIP string, db *geoip2.Reader) (*geoip2.LocationISP, error) {
	ip := net.ParseIP(targetIP)
	if ip == nil {
		return nil, fmt.Errorf("%s is not a valid ip address", targetIP)
	}
	return db.LocationISP(ip)
}

type IPInfo struct {
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	City          string  `json:"city"`
	Postal        string  `json:"zip"`
	TimeZone      string  `json:"timezone"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	ISP           string  `json:"isp"`
	UserType      string  `json:"user_type"`
}

func GetIPInfoFromLocationISP(info *geoip2.LocationISP, lang string) *IPInfo {
	ipInfo := &IPInfo{
		ContinentCode: info.Continent.Code,
		CountryCode:   info.Country.IsoCode,
		Postal:        info.Postal.Code,
		TimeZone:      info.Location.TimeZone,
		Latitude:      info.Location.Latitude,
		Longitude:     info.Location.Longitude,
	}
	// language
	var secondLang string
	if lang == "en" {
		secondLang = "zh-CN"
	} else if lang == "zh-CN" {
		secondLang = "en"
	} else {
		lang = "zh-CN"
		secondLang = "en"
	}
	// isp
	ispNameCn := make([]string, 0)
	for _, isp := range strings.Split(info.Traits.ISP, "/") {
		if ispCn, ok := ispName[isp]; ok {
			ispNameCn = append(ispNameCn, ispCn)
		} else {
			ispNameCn = append(ispNameCn, isp)
		}
	}
	if lang == "zh-CN" {
		ipInfo.ISP = strings.Join(ispNameCn, "/")
	} else {
		ipInfo.ISP = info.Traits.ISP
	}
	// continent
	if continent, ok := info.Continent.Names[lang]; ok {
		ipInfo.Continent = continent
	} else if continent, ok := info.Continent.Names[secondLang]; ok {
		ipInfo.Continent = continent
	} else {
		for _, v := range info.Continent.Names {
			ipInfo.Continent = v
			break
		}
	}
	// country
	if country, ok := info.Country.Names[lang]; ok {
		ipInfo.Country = country
	} else if country, ok := info.Country.Names[secondLang]; ok {
		ipInfo.Country = country
	} else {
		for _, v := range info.Country.Names {
			ipInfo.Country = v
			break
		}
	}
	// province
	if len(info.Subdivisions) > 0 {
		if v, ok := info.Subdivisions[0].Names[lang]; ok {
			ipInfo.Region = v
		} else if v, ok := info.Subdivisions[0].Names[secondLang]; ok {
			ipInfo.Region = v
		} else {
			for _, v := range info.Subdivisions[0].Names {
				ipInfo.Region = v
				break
			}
		}
		ipInfo.RegionCode = info.Subdivisions[0].IsoCode
	}
	// city
	if len(info.Subdivisions) > 1 {
		if v, ok := info.Subdivisions[1].Names[lang]; ok {
			ipInfo.City = v
		} else if v, ok := info.Subdivisions[1].Names[secondLang]; ok {
			ipInfo.City = v
		} else {
			for _, v := range info.Subdivisions[1].Names {
				ipInfo.City = v
				break
			}
		}
	}
	// user type
	var userType string
	if lang == "zh-CN" {
		usertype, ok := map[string]string{
			"hosting":     "数据中心",
			"corporate":   "商业公司",
			"business":    "商业公司",
			"consumer":    "家庭住宅",
			"cellular":    "蜂窝网络",
			"residential": "家庭住宅",
		}[info.Traits.UserType]
		if !ok {
			userType = info.Traits.UserType
		} else {
			userType = usertype
		}
	} else {
		userType = info.Traits.UserType
	}
	ipInfo.UserType = userType

	return ipInfo
}

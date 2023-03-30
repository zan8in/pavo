package hunter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zan8in/pavo/pkg/retryhttpclient"
)

type (
	HunterOptions struct {
		page     int    `json:"page"`
		size     int    `json:"size"`
		Key      string `json:"key"`
		queryAPI string `json:"query_api"`
	}

	HunterResultList struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    Data   `json:"data"`
	}
	Data struct {
		AccountType  string `json:"account_type"`
		Total        int    `json:"total"`
		Time         int    `json:"time"`
		ConsumeQuota string `json:"consume_quota"`
		RestQuota    string `json:"rest_quota"`
		SyntaxPrompt string `json:"syntax_prompt"`
		Arr          []Arr  `json:"arr"`
	}
	Arr struct {
		Url        string      `json:"url"`
		Ip         string      `json:"ip"`
		Port       int         `json:"port"`
		WebTitle   string      `json:"web_title"`
		Domain     string      `json:"domain"`
		Protocol   string      `json:"protocol"`
		Components []Component `json:"components"`
	}
	Component struct {
		Name    string `json:"name"`
		Version string `json:"version`
	}
)

var (
	defaultPage = 1
	defaultSize = 100
)

func New(options *HunterOptions) (*HunterOptions, error) {
	if len(options.Key) == 0 {
		return options, fmt.Errorf("hunter requires a Key")
	}

	return options, nil
}

func (hunter *HunterOptions) Query(query string) (*HunterResultList, error) {
	var (
		hunterResultList *HunterResultList
		err              error
	)
	defer hunter.ReSet()

	if len(query) == 0 {
		return hunterResultList, fmt.Errorf("hunter search content cannot be empty")
	}

	hunter.queryAPI = fmt.Sprintf("%s&search=%s", hunter.queryAPI, base64.StdEncoding.EncodeToString([]byte(query)))

	// fmt.Println(hunter.queryAPI)

	body, err := retryhttpclient.Get(hunter.queryAPI)
	if err != nil {
		return hunterResultList, err
	}

	// fmt.Println(string(body))

	if err = json.Unmarshal(body, &hunterResultList); err != nil {
		return hunterResultList, err
	}

	return hunterResultList, nil
}

func (hunter *HunterOptions) ReSet() {
	hunter.queryAPI = fmt.Sprintf("https://hunter.qianxin.com/openApi/search?api-key=%s", hunter.Key)
	if hunter.page == 0 {
		hunter.SetPage(defaultPage)
	}
	if hunter.size == 0 {
		hunter.SetSize(defaultSize)
	}
}

func (hunter *HunterOptions) SetSize(size int) {
	hunter.size = size
	hunter.queryAPI = fmt.Sprintf("%s&page_size=%d", hunter.queryAPI, hunter.size)
}

func (hunter *HunterOptions) SetPage(page int) {
	hunter.page = page
	hunter.queryAPI = fmt.Sprintf("%s&page=%d", hunter.queryAPI, hunter.page)
}

func (hunter *HunterOptions) HunterResultList2Slice(list *HunterResultList) [][]string {
	if list != nil {
		if len(list.Data.Arr) > 0 {
			arr := make([][]string, len(list.Data.Arr))

			for k, d := range list.Data.Arr {
				arr[k] = append(arr[k], []string{d.Url, d.Ip, strconv.Itoa(d.Port), d.WebTitle, d.Domain, d.Protocol}...)
			}

			return arr
		}
	}
	return [][]string{}
}

package fofa

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zan8in/goflags"
	"github.com/zan8in/pavo/pkg/retryhttpclient"
)

type (
	FofaOptions struct {
		fields   goflags.StringSlice `json:"fields"`
		page     int                 `json:"page"`
		size     int                 `json:"size"`
		full     bool                `json:"full"`
		Email    string              `json:"email"`
		Key      string              `json:"key"`
		queryAPI string              `json:"query_api"`
	}

	FofaResultList struct {
		Error          bool
		ConsumedFpoint int
		Size           int
		Page           int
		Mode           string
		Query          string
		Results        [][]string
	}

	FofaResult2 struct {
		Ip              string
		Port            int
		Protocol        string
		Country         string
		CountryName     string
		Region          string
		City            string
		Longitude       string
		Latitude        string
		AsNumber        string
		AsOrganization  string
		Host            string
		Domain          string
		Os              string
		Server          string
		Icp             string
		Title           string
		Jarm            string
		Header          string
		Banner          string
		Cert            string
		BaseProtocol    string
		Product         string
		ProductCategory string
		Version         string
		Lastupdatetime  string
		Cname           string
		IconHash        string
		CertsValid      string
		CnameDomain     string
		Body            string
		Icon            string
		Fid             string
		Structinfo      string
	}
)

func New(options *FofaOptions) (*FofaOptions, error) {
	if len(options.Email) == 0 || len(options.Key) == 0 {
		return options, fmt.Errorf("fofa requires an email address or Key")
	}

	return options, nil
}

func (fofa *FofaOptions) Query(qbase64 string) (FofaResultList, error) {
	var (
		fofaResultList FofaResultList
		err            error
	)
	defer fofa.ReSet()

	if len(qbase64) == 0 {
		return fofaResultList, fmt.Errorf("qbase64 cannot be empty")
	}

	fofa.queryAPI = fmt.Sprintf("%s&qbase64=%s", fofa.queryAPI, base64.StdEncoding.EncodeToString([]byte(qbase64)))

	body, err := retryhttpclient.Get(fofa.queryAPI)
	if err != nil {
		return fofaResultList, err
	}

	if err = json.Unmarshal(body, &fofaResultList); err != nil {
		return fofaResultList, err
	}

	return fofaResultList, nil
}

func (fofa *FofaOptions) ReSet() {
	fofa.queryAPI = fmt.Sprintf("https://fofa.info/api/v1/search/all?email=%s&key=%s", fofa.Email, fofa.Key)
	fofa.SetFields([]string{"host", "title", "ip", "port", "domain", "protocol", "server"})
	fofa.SetFull(true)
}

func (fofa *FofaOptions) SetSize(size int) {
	fofa.size = size
	fofa.queryAPI = fmt.Sprintf("%s&size=%d", fofa.queryAPI, fofa.size)
}

func (fofa *FofaOptions) SetPage(page int) {
	fofa.page = page
	fofa.queryAPI = fmt.Sprintf("%s&page=%d", fofa.queryAPI, fofa.page)
}

func (fofa *FofaOptions) SetFull(full bool) {
	fofa.full = full
	fofa.queryAPI = fmt.Sprintf("%s&full=%T", fofa.queryAPI, fofa.full)
}

func (fofa *FofaOptions) SetFields(fields []string) {
	fofa.fields = fields
	fofa.queryAPI = fmt.Sprintf("%s&fields=%s", fofa.queryAPI, strings.Join(fofa.fields, ","))
}

package dadata

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"time"
)

func (c Client) GetByINN(ctx context.Context, inn string) (res []Suggestion, err error) {
	const to = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, to)
	defer cancel()

	tok, err := c.tp.Get(ctx)
	if err != nil {
		return
	}

	b, err := json.Marshal(daDataReq{
		Query: inn,
		Count: 20,
	})
	if err != nil {
		return
	}

	dReq := fasthttp.AcquireRequest()
	dReq.Header.SetMethod(fasthttp.MethodPost)
	dReq.Header.SetContentType("application/json")
	dReq.Header.Set("Authorization", "Token "+tok)
	dReq.SetRequestURI("https://suggestions.dadata.ru/suggestions/api/4_1/rs/findById/party")
	dReq.SetBody(b)
	defer fasthttp.ReleaseRequest(dReq)

	dRes := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(dRes)

	cl := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		MaxConnDuration:          to,
		MaxConnWaitTimeout:       to,
		MaxIdleConnDuration:      to,
		ReadTimeout:              to,
		WriteTimeout:             to,
	}

	err = cl.DoTimeout(dReq, dRes, to)
	if err != nil {
		return
	}
	if dRes.StatusCode() != fasthttp.StatusOK {
		err = errors.New(string(dRes.Body()))
		return
	}

	var body daDataRes
	err = json.Unmarshal(dRes.Body(), &body)

	res = body.Suggestions
	return
}

type Suggestion struct {
	Value string `json:"value"`
	Data  data   `json:"data"`
}

type daDataReq struct {
	Query string `json:"query"`
	Count uint32 `json:"count"`
}

type daDataRes struct {
	Suggestions []Suggestion `json:"suggestions"`
}

type management struct {
	Name string `json:"name"`
	Post string `json:"post"`
}

type state struct {
	Status           string `json:"status"`
	RegistrationDate int64  `json:"registration_date"`
	LiquidationDate  int64  `json:"liquidation_date"`
}

type opf struct {
	Type  string `json:"type"`
	Code  string `json:"code"`
	Full  string `json:"full"`
	Short string `json:"short"`
}

type name struct {
	FullWithOpf  string `json:"full_with_opf"`
	ShortWithOpf string `json:"short_with_opf"`
	Full         string `json:"full"`
	Short        string `json:"short"`
}

type metro struct {
	Name     string  `json:"name"`
	Line     string  `json:"line"`
	Distance float64 `json:"distance"`
}

type AddressData struct {
	CityFiasID         string  `json:"city_fias_id"`
	CityKladrID        string  `json:"city_kladr_id"`
	CityType           string  `json:"city_type"`
	CityTypeFull       string  `json:"city_type_full"`
	City               string  `json:"city"`
	SettlementFiasID   string  `json:"settlement_fias_id"`
	SettlementKladrID  string  `json:"settlement_kladr_id"`
	SettlementType     string  `json:"settlement_type"`
	SettlementTypeFull string  `json:"settlement_type_full"`
	Settlement         string  `json:"settlement"`
	Metro              []metro `json:"metro"`
}

type address struct {
	UnrestrictedValue string      `json:"unrestricted_value"`
	Data              AddressData `json:"data"`
}

type data struct {
	Hid         string     `json:"hid"`
	Kpp         string     `json:"kpp"`
	Management  management `json:"management"`
	BranchType  string     `json:"branch_type"`
	BranchCount int        `json:"branch_count"`
	Type        string     `json:"type"`
	State       state      `json:"state"`
	Opf         opf        `json:"opf"`
	Name        name       `json:"name"`
	Inn         string     `json:"inn"`
	Ogrn        string     `json:"ogrn"`
	Okpo        string     `json:"okpo"`
	Okato       string     `json:"okato"`
	Oktmo       string     `json:"oktmo"`
	Okogu       string     `json:"okogu"`
	Okfs        string     `json:"okfs"`
	Address     address    `json:"address"`
	OgrnDate    int64      `json:"ogrn_date"`
	OkvedType   string     `json:"okved_type"`
}

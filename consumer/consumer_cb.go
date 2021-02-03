package consumer

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-org-producer/protocol"
	"time"
)

func (c Consumer) cb(rawMsg *stan.Msg) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ack := func() {
			e := rawMsg.Ack()
			if e != nil {
				c.logger.Error().Err(e).Send()
			}
		}

		var msg protocol.OrgMessage
		err := json.Unmarshal(rawMsg.Data, &msg)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}
	}()
}

type daDataRes struct {
	Data Data `json:"data"`
}

type Management struct {
	Name string `json:"name"`
	Post string `json:"post"`
}

type State struct {
	Status           string `json:"status"`
	RegistrationDate int64  `json:"registration_date"`
	LiquidationDate  int64  `json:"liquidation_date"`
}

type Opf struct {
	Type  string `json:"type"`
	Code  string `json:"code"`
	Full  string `json:"full"`
	Short string `json:"short"`
}

type Name struct {
	FullWithOpf  string `json:"full_with_opf"`
	ShortWithOpf string `json:"short_with_opf"`
	Full         string `json:"full"`
	Short        string `json:"short"`
}

type Metro struct {
	Name     string  `json:"name"`
	Line     string  `json:"line"`
	Distance float64 `json:"distance"`
}

type AddressData struct {
	CityFiasID         string  `json:"city_fias_id"`
	CityKladrID        string  `json:"city_kladr_id"`
	CityWithType       string  `json:"city_with_type"`
	CityType           string  `json:"city_type"`
	CityTypeFull       string  `json:"city_type_full"`
	City               string  `json:"city"`
	SettlementFiasID   string  `json:"settlement_fias_id"`
	SettlementKladrID  string  `json:"settlement_kladr_id"`
	SettlementWithType string  `json:"settlement_with_type"`
	SettlementType     string  `json:"settlement_type"`
	SettlementTypeFull string  `json:"settlement_type_full"`
	Settlement         string  `json:"settlement"`
	Metro              []Metro `json:"metro"`
}

type Address struct {
	UnrestrictedValue string      `json:"unrestricted_value"`
	Data              AddressData `json:"data"`
}

type Data struct {
	Hid         string      `json:"hid"`
	Kpp         string      `json:"kpp"`
	Management  Management  `json:"management"`
	BranchType  string      `json:"branch_type"`
	BranchCount int         `json:"branch_count"`
	Type        string      `json:"type"`
	State       State       `json:"state"`
	Opf         Opf         `json:"opf"`
	Name        Name        `json:"name"`
	Inn         string      `json:"inn"`
	Ogrn        string      `json:"ogrn"`
	Okpo        string      `json:"okpo"`
	Okato       string      `json:"okato"`
	Oktmo       string      `json:"oktmo"`
	Okogu       string      `json:"okogu"`
	Okfs        string      `json:"okfs"`
	Okved       string      `json:"okved"`
	Okveds      interface{} `json:"okveds"`
	Authorities interface{} `json:"authorities"`
	Documents   interface{} `json:"documents"`
	Licenses    interface{} `json:"licenses"`
	Finance     interface{} `json:"finance"`
	Address     Address     `json:"address"`
	OgrnDate    int64       `json:"ogrn_date"`
	OkvedType   string      `json:"okved_type"`
}

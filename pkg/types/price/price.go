package price

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Price struct {
	// цена в копейках
	FullPrice uint
}

func (p *Price) Kopecks() uint {
	return p.FullPrice % 100
}

func (p *Price) Rubles() uint {
	return p.FullPrice / 100
}

func (p *Price) FullPriceKopecks() uint {
	return p.FullPrice
}

func (p *Price) String() string {
	return fmt.Sprintf("%d.%d ₽", p.Rubles(), p.Kopecks())
}

func FromString(str string) (Price, bool) {
	parts := strings.Split(str, ".")
	if len(parts) != 2 {
		parts = strings.Split(str, ",")
		if len(parts) != 2 {
			return Price{}, false
		}
	}
	rubPart, err := strconv.Atoi(parts[0])
	if err != nil {
		return Price{}, false
	}
	if rubPart < 0 {
		return Price{}, false
	}

	kopPart, err := strconv.Atoi(parts[1])
	if err != nil {
		return Price{}, false
	}
	if kopPart < 0 {
		return Price{}, false
	}

	return Price{
		FullPrice: uint(rubPart)*100 + uint(kopPart),
	}, true
}

func FromFloat(num float64) (Price, bool) {
	if num < 0 {
		return Price{}, false
	}

	rubPart := uint(num)
	kopPart := uint((num - float64(rubPart)) * 100)

	return Price{
		FullPrice: rubPart*100 + kopPart,
	}, true
}

func (p *Price) UnmarshalJSON(data []byte) error {
	str := string(data[1 : len(data)-1])
	price, ok := FromString(str)
	if ok {
		p = &price
		return nil
	}

	return fmt.Errorf("Parsing price error")
}

func (p *Price) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

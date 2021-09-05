package hotelProviderInterface

import (
	"github.com/tidwall/gjson"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/logger"
	"io/ioutil"
	"path/filepath"
)

type basicInformationProvider struct {
}

func (p *basicInformationProvider) GetCities() ([]dto.CityDto, error) {

	filePath, err := filepath.Abs(filepath.Join("assets", "cities.json"))
	if err != nil {
		return nil, err
	}

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	jsonString := string(jsonData)
	if !gjson.Valid(jsonString) {
		logger.WithData(map[string]interface{}{
			"response": jsonString,
		}).Error(common.JsonDataIsNotValid.Error())
		return nil, common.JsonDataIsNotValid
	}
	citiesArray := gjson.Parse(jsonString).Array()

	cities := make([]dto.CityDto, 0, len(citiesArray))

	for _, i := range citiesArray {
		cities = append(cities, dto.CityDto{
			Id:      i.Get("id").String(),
			Name:    i.Get("name").String(),
			Country: i.Get("country").String(),
			State:   i.Get("state").String(),
			BaseId:  i.Get("baseId").Int(),
		})
	}

	return cities, nil
}

func NewBasicInformationProvider() core.BasicInformationProvider {
	return &basicInformationProvider{}
}

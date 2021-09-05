package core

import "hotel-engine/core/dto"

//Feeder ...
type Feeder interface {
	Feed(hotels dto.ElasticUpdateRequest) error
	Seed(handle func(index string)) error
	Alias(index string) error
	Close() error
}

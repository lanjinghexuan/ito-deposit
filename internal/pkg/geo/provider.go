package geo

import (
	"github.com/google/wire"
)

// ProviderSet is geo service providers.
var ProviderSet = wire.NewSet(NewGeoService)
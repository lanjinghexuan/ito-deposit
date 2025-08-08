package service

import (
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService, NewUserService, NewOrderService, NewHomeService, NewDepositService, NewCityService, NewNearbyService, NewAdminService, NewGroupService, NewCabinetCellService)

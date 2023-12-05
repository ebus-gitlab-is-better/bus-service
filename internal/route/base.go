package route

import "github.com/google/wire"

// ProviderSet is riute providers.
var ProviderSet = wire.NewSet(NewBusRouter, NewRouteRouter, NewDriverRoute)

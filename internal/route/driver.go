package route

import (
	"bus-service/internal/biz"
	"context"

	"github.com/gin-gonic/gin"
)

type DriverRoute struct {
	uc *biz.DriverUseCase
}

func NewDriverRoute(uc *biz.DriverUseCase) *DriverRoute {
	return &DriverRoute{uc: uc}
}

func (r *DriverRoute) Register(router *gin.RouterGroup) {
	router.GET("/", r.getDrivers)
}

type ListDriverDTO struct {
	Drivers []*biz.Driver `json:"drivers"`
}

// @Summary	Get drivers
// @Accept		json
// @Produce	json
// @Tags		drivers
// @Success	200	{object}	route.ListDriverDTO
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/drivers/ [get]
func (r *DriverRoute) getDrivers(c *gin.Context) {
	drivers, err := r.uc.GetDrivers(context.TODO())
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, &ListDriverDTO{
		Drivers: drivers,
	})
}

package route

import (
	"bus-service/internal/biz"
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BusRouter struct {
	uc *biz.BusUseCase
	v  *validator.Validate
}

func NewBusRouter(uc *biz.BusUseCase) *BusRouter {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &BusRouter{
		uc: uc,
		v:  validate,
	}
}

func (r *BusRouter) Register(router *gin.RouterGroup) {
	router.POST("/", r.create)
	router.GET("/:id", r.getById)
	router.PUT("/:id", r.update)
	router.DELETE("/:id", r.delete)
	router.GET("/", r.list)
}

type BusDTO struct {
	RouteID  *uint32 `validate:"required"`
	DriverID *string
	Number   string `validate:"required"`
	Status   string `validate:"required"`
}

// @Summary	Create bus
// @Accept		json
// @Produce	json
// @Tags		bus
// @Param		dto	body	route.BusDTO	true	"dto"
// @Success	200
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/bus/ [post]
func (r *BusRouter) create(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	dto := BusDTO{}

	err = json.Unmarshal(body, &dto)
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	err = r.v.Struct(dto)
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	err = r.uc.Create(context.TODO(), &biz.BusDTO{
		RouteID:  dto.RouteID,
		DriverID: dto.DriverID,
		Status:   "Не запущен",
		Number:   dto.Number,
	})

	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(200)
}

// @Summary	Update bus
// @Accept		json
// @Produce	json
// @Tags		bus
// @Param		id	path	int	true	"Bus ID"	Format(uint64)
// @Param		dto	body	route.BusDTO	true	"dto"
// @Success	200
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/bus/{id} [put]
func (r *BusRouter) update(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	dto := BusDTO{}

	err = json.Unmarshal(body, &dto)
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	err = r.v.Struct(dto)
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	err = r.uc.Update(context.TODO(), &biz.BusDTO{
		RouteID:  dto.RouteID,
		DriverID: dto.DriverID,
		Status:   "Не запущен",
		Number:   dto.Number,
		Id:       uint32(idUint),
	})

	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(200)
}

// @Summary	Delete bus
// @Accept		json
// @Produce	json
// @Tags		bus
// @Param		id	path	int	true	"Bus ID"	Format(uint64)
// @Success	200
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/bus/{id} [delete]
func (r *BusRouter) delete(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}
	err = r.uc.Delete(context.TODO(), uint32(idUint))
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(200)
}

// @Summary	Get bus by id
// @Accept		json
// @Produce	json
// @Tags		bus
// @Param		id	path	int	true	"Bus ID"	Format(uint64)
// @Success	200	{object}	biz.Bus
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/bus/{id} [get]
func (r *BusRouter) getById(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}
	bus, err := r.uc.GetById(context.TODO(), uint32(idUint))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}
	c.JSON(200, bus)
}

type ListBuses struct {
	Buses []*biz.Bus
	Count int64
}

// @Summary	List buses
// @Accept		json
// @Produce	json
// @Tags		bus
// @Success	200	{object}	route.ListBuses
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/bus/ [get]
func (r *BusRouter) list(c *gin.Context) {
	buses, total, err := r.uc.List(context.TODO())
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}
	c.JSON(200, &ListBuses{
		Buses: buses,
		Count: total,
	})
}

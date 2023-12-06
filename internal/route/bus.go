package route

import (
	"bus-service/internal/biz"
	"context"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BusRouter struct {
	uc  *biz.BusUseCase
	v   *validator.Validate
	ucS *biz.ShiftUseCase
}

func NewBusRouter(uc *biz.BusUseCase, ucS *biz.ShiftUseCase) *BusRouter {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &BusRouter{
		uc:  uc,
		v:   validate,
		ucS: ucS,
	}
}

func (r *BusRouter) Register(router *gin.RouterGroup) {
	router.POST("/", r.create)
	router.GET("/:id", r.getById)
	router.PUT("/:id", r.update)
	router.DELETE("/:id", r.delete)
	router.GET("/", r.list)
	router.POST("/:id/start", r.start)
	router.POST("/:id/charge", r.charge)
	router.POST("/:id/stop", r.stop)
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

// @Summary	Водитель начинает смену
// @Accept		json
// @Produce	json
// @Tags		bus
// @Success	200
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/bus/{id}start [post]
func (r *BusRouter) start(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}
	userD, ok := c.Get("user")
	if !ok {
		return
	}
	user, ok := userD.(*gocloak.UserInfo)
	if !ok {
		return
	}
	err = r.ucS.Create(context.TODO(), &biz.Shift{
		StartTime: time.Now(),
		DriverID:  *user.Sub,
	})
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	bus, err := r.uc.GetById(context.TODO(), uint32(idUint))
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	err = r.uc.Update(context.TODO(), &biz.BusDTO{
		Id:       bus.Id,
		RouteID:  bus.RouteID,
		Number:   bus.Number,
		Status:   "В работе",
		DriverID: user.Sub,
	})
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(200)
}

// @Summary	Водитель заканчивает смену
// @Accept		json
// @Produce	json
// @Tags		bus
// @Success	200
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/bus/{id}/stop [post]
func (r *BusRouter) stop(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}
	userD, ok := c.Get("user")
	if !ok {
		return
	}
	user, ok := userD.(*gocloak.UserInfo)
	if !ok {
		return
	}
	shift, err := r.ucS.GetByDriverID(context.TODO(), *user.Sub)
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	endTime := time.Now()
	err = r.ucS.Update(context.TODO(), &biz.Shift{
		Id:        shift.Id,
		StartTime: shift.StartTime,
		EndDate:   &endTime,
		DriverID:  shift.DriverID,
	})
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	bus, err := r.uc.GetById(context.TODO(), uint32(idUint))
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	err = r.uc.Update(context.TODO(), &biz.BusDTO{
		Id:       bus.Id,
		RouteID:  bus.RouteID,
		Number:   bus.Number,
		Status:   "Не в работе",
		DriverID: nil,
	})
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(200)
}

// @Summary	Автобус на зарядке
// @Accept		json
// @Produce	json
// @Tags		bus
// @Success	200
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/bus/{id}/charge [post]
func (r *BusRouter) charge(c *gin.Context) {
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
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	err = r.uc.Update(context.TODO(), &biz.BusDTO{
		Id:       bus.Id,
		RouteID:  bus.RouteID,
		Number:   bus.Number,
		Status:   "На зарядке",
		DriverID: bus.Driver.Id,
	})
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(200)
}

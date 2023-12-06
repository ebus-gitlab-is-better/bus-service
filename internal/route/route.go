package route

import (
	mapS "bus-service/api/map/v1"
	"bus-service/internal/biz"
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RouteRouter struct {
	uc        *biz.RouteUseCase
	mapClient mapS.MapClient
	v         *validator.Validate
}

func NewRouteRouter(uc *biz.RouteUseCase, mapClient mapS.MapClient) *RouteRouter {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &RouteRouter{uc: uc, mapClient: mapClient, v: validate}
}

func (r *RouteRouter) Register(router *gin.RouterGroup) {
	router.POST("/", r.create)
	router.GET("/:id", r.getById)
	// router.PUT("/:id", r.update)
	router.DELETE("/:id", r.delete)
	router.GET("/", r.list)
}

type StationDTO struct {
	ID   uint32
	Name string  `validate:"required"`
	Lat  float64 `validate:"required"`
	Lon  float64 `validate:"required"`
}

type RouteDTO struct {
	Number string `validate:"required"`
	// Path     string
	Stations []StationDTO `validate:"required"`
}

// @Summary	Create route
// @Accept		json
// @Produce	json
// @Tags		route
// @Param		dto	body	route.RouteDTO	true	"dto"
// @Success	200
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/route/ [post]
func (r *RouteRouter) create(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	dto := RouteDTO{}

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
	stations := make([]biz.Stations, 0)
	for _, station := range dto.Stations {
		stations = append(stations, biz.Stations{
			Lat:  station.Lat,
			Name: station.Name,
			Lon:  station.Lon,
		})
	}
	points := make([]*mapS.Point, 0)
	for _, station := range stations {
		points = append(points, &mapS.Point{
			Lat: float32(station.Lat),
			Lon: float32(station.Lon),
		})
	}
	req, err := r.mapClient.GetPath(context.TODO(), &mapS.GetPathRequest{
		Points: points,
	})
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	err = r.uc.Create(context.TODO(), &biz.Route{
		Number:   dto.Number,
		Path:     req.Shape,
		Time:     req.Time,
		Stations: stations,
		Lengths:  req.Lengths,
		Length:   req.Length,
	})

	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(200)
}

// @Summary	Update route
// @Accept		json
// @Produce	json
// @Tags		route
// @Param		id	path	int	true	"Route ID"	Format(uint64)
// @Param		dto	body	route.RouteDTO	true	"dto"
// @Success	200
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/route/{id} [put]
// func (r *RouteRouter) update(c *gin.Context) {
// 	id := c.Param("id")
// 	idUint, err := strconv.Atoi(id)

// 	if err != nil {
// 		c.AbortWithStatusJSON(400, gin.H{
// 			"error": "parse id error",
// 		})
// 		return
// 	}

// 	body, err := io.ReadAll(c.Request.Body)

// 	if err != nil {
// 		c.JSON(400, &gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	dto := RouteDTO{}

// 	err = json.Unmarshal(body, &dto)
// 	if err != nil {
// 		c.AbortWithStatusJSON(400, &gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	err = r.v.Struct(dto)
// 	if err != nil {
// 		c.AbortWithStatusJSON(400, &gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	stations := make([]biz.Stations, 0)
// 	for _, station := range dto.Stations {
// 		stations = append(stations, biz.Stations{
// 			ID:   uint(station.ID),
// 			Lat:  station.Lat,
// 			Name: station.Name,
// 			Lon:  station.Lon,
// 		})
// 	}
// 	points := make([]*mapS.Point, 0)
// 	for _, station := range stations {
// 		points = append(points, &mapS.Point{
// 			Lat: float32(station.Lat),
// 			Lon: float32(station.Lon),
// 		})
// 	}
// 	req, err := r.mapClient.GetPath(context.TODO(), &mapS.GetPathRequest{
// 		Points: points,
// 	})
// 	if err != nil {
// 		c.AbortWithStatusJSON(400, &gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	err = r.uc.Update(context.TODO(), &biz.Route{
// 		Id:       uint32(idUint),
// 		Number:   dto.Number,
// 		Path:     req.Shape,
// 		Stations: stations,
// 	})

// 	if err != nil {
// 		c.AbortWithStatusJSON(400, &gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.Status(200)
// }

// @Summary	Delete route
// @Accept		json
// @Produce	json
// @Tags		route
// @Param		id	path	int	true	"Route ID"	Format(uint64)
// @Success	200
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/route/{id} [delete]
func (r *RouteRouter) delete(c *gin.Context) {
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

// @Summary	Get route
// @Accept		json
// @Produce	json
// @Tags		route
// @Param		id	path	int	true	"Route ID"	Format(uint64)
// @Success	200 {object} biz.Route
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/route/{id} [get]
func (r *RouteRouter) getById(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}

	route, err := r.uc.GetById(context.TODO(), uint32(idUint))

	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, route)
}

type ListRoute struct {
	Routes []*biz.Route
	Count  int64
}

// @Summary	List route
// @Accept		json
// @Produce	json
// @Tags		route
// @Success	200 {object} biz.Route
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/route/ [get]
func (r *RouteRouter) list(c *gin.Context) {
	routes, total, err := r.uc.List(context.TODO())
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}
	c.JSON(200, &ListRoute{
		Routes: routes,
		Count:  total,
	})
}

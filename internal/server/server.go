package server

import (
	"aleksandersh.github.io/planning-poker-server/internal/activity/activitydata"
	"aleksandersh.github.io/planning-poker-server/internal/controller"
	"aleksandersh.github.io/planning-poker-server/internal/rooms/roomsdata"
	"aleksandersh.github.io/planning-poker-server/internal/rooms/roomsdomain"
	"aleksandersh.github.io/planning-poker-server/internal/users/usersdata"
	"aleksandersh.github.io/planning-poker-server/internal/users/usersdomain"
	"github.com/gin-gonic/gin"
)

func Start(address string) {
	router := gin.Default()

	ar := activitydata.NewRepository()
	us := usersdomain.NewService(usersdata.NewRepo(), ar)
	uc := controller.NewUsersController(us)

	router.POST("/v1/users/register", uc.Register)

	ah := controller.NewAuthHelper(us)
	rr := roomsdata.NewRepo()
	rs := roomsdomain.NewRoomsService(rr, ar)
	rc := controller.NewRoomsController(ah, rs)

	router.POST("/v1/rooms", rc.Post)
	router.GET("/v1/rooms/:room_id", rc.Get)
	router.DELETE("/v1/rooms/:room_id", rc.Delete)
	router.POST("/v1/rooms/:room_id/join", rc.Join)
	router.GET("/v1/rooms/:room_id/state", rc.GetState)

	gs := roomsdomain.NewGamesService(rr, ar)
	gc := controller.NewGamesController(ah, gs)

	router.POST("/v1/games", gc.Post)
	router.POST("/v1/games/:game_id/complete", gc.Complete)
	router.POST("/v1/games/:game_id/reset", gc.Reset)
	router.POST("/v1/games/:game_id/send-card", gc.SendCard)
	router.POST("/v1/games/:game_id/drop-card", gc.DropCard)

	router.Run(address)
}

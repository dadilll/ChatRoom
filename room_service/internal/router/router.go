package router

import (
	"crypto/rsa"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"room_service/internal/handler"
	middleware "room_service/internal/middelware"
	"room_service/internal/models"
	"room_service/internal/service"
	"room_service/internal/storage"
	"room_service/pkg/logger"
)

type Dependencies struct {
	RoomHandler   *handler.RoomHandler
	RoleHandler   *handler.RoleHandler
	MemberHandler *handler.MemberHandler
	InviteHandler *handler.InviteHandler
	PublicKey     *rsa.PublicKey
	RoleService   service.RoleService
}

func setupDependencies(logger logger.Logger, db *sqlx.DB, publicKey *rsa.PublicKey) *Dependencies {
	roomStorage := storage.NewRoomStorage(db.DB)
	roleStorage := storage.NewRoleStorage(db.DB)
	memberStorage := storage.NewRoomMemberStorage(db.DB)
	inviteStorage := storage.NewPostgresInviteStorage(db.DB)
	inviteService := service.NewInviteService(inviteStorage)

	roomService := service.NewRoomService(logger, roomStorage)
	roleService := service.NewRoleService(logger, roleStorage, roomStorage)
	memberService := service.NewMemberService(memberStorage, roleService, logger, roomStorage, inviteService)
	inviteHandler := handler.NewInviteHandler(inviteService)

	return &Dependencies{
		RoomHandler:   handler.NewRoomHandler(roomService, logger),
		RoleHandler:   handler.NewRoleHandler(roleService, logger),
		MemberHandler: handler.NewMemberHandler(memberService, logger),
		InviteHandler: inviteHandler,
		PublicKey:     publicKey,
		RoleService:   roleService,
	}
}

func setupRoomRoutes(g *echo.Group, h *handler.RoomHandler, requireManageRoom echo.MiddlewareFunc) {
	g.POST("/rooms", h.CreateRoom)
	g.GET("/rooms/:room_id", h.GetRoom)
	g.PUT("/rooms/:room_id", h.UpdateRoom, requireManageRoom)
	g.DELETE("/rooms/:room_id", h.DeleteRoom, requireManageRoom)
}

func setupRoleRoutes(
	g *echo.Group,
	h *handler.RoleHandler,
	requireManageRoles echo.MiddlewareFunc,
) {
	g.POST("/rooms/:room_id/roles", h.CreateRole, requireManageRoles)
	g.GET("/roles/:role_id", h.GetRole)
	g.GET("/rooms/:room_id/roles", h.GetRoomRoles)
	g.PUT("/roles/:role_id", h.UpdateRole, requireManageRoles)
	g.DELETE("/roles/:role_id", h.DeleteRole, requireManageRoles)

	g.PUT("/rooms/:room_id/members/:user_id/role", h.AssignRole, requireManageRoles)
	g.DELETE("/rooms/:room_id/members/:user_id/role", h.RemoveRole, requireManageRoles)
	g.GET("/rooms/:room_id/members/:user_id/role", h.GetUserRole)
}

func setupMemberRoutes(
	g *echo.Group,
	h *handler.MemberHandler,
) {
	g.POST("/rooms/:room_id/members/add", h.AddMember)
	g.GET("/rooms/:room_id/members", h.ListMembers)
	g.DELETE("/rooms/:room_id/members/leave", h.RemoveMember)
}

func setupInviteRoutes(
	g *echo.Group,
	h *handler.InviteHandler,
) {
	g.POST("/invites", h.CreateInvite)
	g.GET("/invites/:user_id", h.GetUserInvites)
	g.POST("/invites/:invite_id/accept", h.AcceptInvite)
	g.POST("/invites/:invite_id/decline", h.DeclineInvite)
	g.DELETE("/invites/:invite_id", h.DeleteInvite)
}

func SetupRoutes(e *echo.Echo, logger logger.Logger, db *sqlx.DB, publicKey *rsa.PublicKey) {
	deps := setupDependencies(logger, db, publicKey)
	api := e.Group("/api/v1", middleware.AuthMiddleware(deps.PublicKey))

	requireManageRoles := middleware.RequirePermission(
		deps.RoleService,
		logger,
		models.PermissionManageRoles,
	)
	requireManageRoom := middleware.RequirePermission(
		deps.RoleService,
		logger,
		models.PermissionManageRoom,
	)

	setupRoomRoutes(api, deps.RoomHandler, requireManageRoom)
	setupRoleRoutes(api, deps.RoleHandler, requireManageRoles)
	setupInviteRoutes(api, deps.InviteHandler)
	setupMemberRoutes(api, deps.MemberHandler)
}

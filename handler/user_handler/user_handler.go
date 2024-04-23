package userhandler

import "github.com/gin-gonic/gin"

type UserHandler interface {
	FindAll(ctx *gin.Context)
	FindOneByUsername(ctx *gin.Context)
	Create(ctx *gin.Context)
	AssignReassignRoleToUser(ctx *gin.Context)
}

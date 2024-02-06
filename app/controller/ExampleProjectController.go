package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang-structure-template-with-di/app/helper"
	"golang-structure-template-with-di/app/models/request"
	"golang-structure-template-with-di/app/service"
	"golang-structure-template-with-di/libraries/httpResponse"
)

type ExampleProjectController struct {
	authService           *service.AuthService           `di.inject:"AuthService"`
	ExampleProjectService *service.ExampleProjectService `di.inject:"ExampleProjectService"`
	validator             *validator.Validate            `di.inject:"validator"`
}

func (global *ExampleProjectController) ExampleController(c *gin.Context) {

	/*-- Check user permission with decoded JWT Token --*/
	checkUserPermission := global.authService.UserHasPermissions(c, "602")
	checkUserRoles := global.authService.UserHasRoles(c, "superadmin")
	if !checkUserRoles && !checkUserPermission {
		httpResponse.Forbidden()
		return
	}

	//*-- Binding and validating data from request body --*/
	var requestData request.CreateExampleProject
	errorValidation := helper.BindAndValidate(c, &requestData, global.validator)
	if errorValidation != nil {
		httpResponse.BadRequest(errorValidation)
	}

	dataResp := global.ExampleProjectService.ExampleProject(requestData)
	httpResponse.HttpCreated(c, "Success Create ContainerProject", dataResp)

}

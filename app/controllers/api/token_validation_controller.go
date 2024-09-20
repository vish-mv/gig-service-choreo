package api

import (
	"GIG/app/constants/info_messages"
	"GIG/app/controllers"
	"github.com/revel/revel"
)

type TokenValidationController struct {
	*revel.Controller
}

// swagger:operation GET /token/validate User validate
//
// Authenticate Validation of User
//
// This API allows to validate a token
//
// ---
// produces:
// - application/json
//
// security:
//   - Bearer: []
//   - ApiKey: []
//
// responses:
//   '200':
//     description: login success
//     schema:
//         "$ref": "#/definitions/Response"
//   '403':
//     description: input validation error
//     schema:
////       "$ref": "#/definitions/Response"
//   '500':
//     description: server error
//     schema:
//       "$ref": "#/definitions/Response"
func (c TokenValidationController) ValidateToken() revel.Result {
	return c.RenderJSON(controllers.BuildSuccessResponse(info_messages.TokenIsValid, 200))
}

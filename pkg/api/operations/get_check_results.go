// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"
	"strconv"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"

	models "github.com/supergiant/robot/pkg/models"
)

// GetCheckResultsHandlerFunc turns a function with the right signature into a get check results handler
type GetCheckResultsHandlerFunc func(GetCheckResultsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetCheckResultsHandlerFunc) Handle(params GetCheckResultsParams) middleware.Responder {
	return fn(params)
}

// GetCheckResultsHandler interface for that can handle valid get check results params
type GetCheckResultsHandler interface {
	Handle(GetCheckResultsParams) middleware.Responder
}

// NewGetCheckResults creates a new http.Handler for the get check results operation
func NewGetCheckResults(ctx *middleware.Context, handler GetCheckResultsHandler) *GetCheckResults {
	return &GetCheckResults{Context: ctx, Handler: handler}
}

/*GetCheckResults swagger:route GET /check getCheckResults

Returns list of check results produced by installed plugins

*/
type GetCheckResults struct {
	Context *middleware.Context
	Handler GetCheckResultsHandler
}

func (o *GetCheckResults) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetCheckResultsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetCheckResultsOKBody get check results o k body
// swagger:model GetCheckResultsOKBody
type GetCheckResultsOKBody struct {

	// existing checks
	CheckResults []*models.CheckResult `json:"CheckResults"`

	// total count
	TotalCount int64 `json:"TotalCount,omitempty"`
}

// Validate validates this get check results o k body
func (o *GetCheckResultsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCheckResults(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetCheckResultsOKBody) validateCheckResults(formats strfmt.Registry) error {

	if swag.IsZero(o.CheckResults) { // not required
		return nil
	}

	for i := 0; i < len(o.CheckResults); i++ {
		if swag.IsZero(o.CheckResults[i]) { // not required
			continue
		}

		if o.CheckResults[i] != nil {
			if err := o.CheckResults[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getCheckResultsOK" + "." + "CheckResults" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetCheckResultsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetCheckResultsOKBody) UnmarshalBinary(b []byte) error {
	var res GetCheckResultsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

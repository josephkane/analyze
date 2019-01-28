// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetPluginParams creates a new GetPluginParams object
// no default values defined in spec.
func NewGetPluginParams() GetPluginParams {

	return GetPluginParams{}
}

// GetPluginParams contains all the bound params for the get plugin operation
// typically these are obtained from a http.Request
//
// swagger:parameters getPlugin
type GetPluginParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The id of the plugin to retrieve
	  Required: true
	  In: path
	*/
	PluginID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetPluginParams() beforehand.
func (o *GetPluginParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rPluginID, rhkPluginID, _ := route.Params.GetOK("pluginId")
	if err := o.bindPluginID(rPluginID, rhkPluginID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindPluginID binds and validates parameter PluginID from path.
func (o *GetPluginParams) bindPluginID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.PluginID = raw

	return nil
}

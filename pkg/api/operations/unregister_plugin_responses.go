// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/supergiant/analyze/pkg/models"
)

// UnregisterPluginNoContentCode is the HTTP code returned for type UnregisterPluginNoContent
const UnregisterPluginNoContentCode int = 204

/*UnregisterPluginNoContent plugin is removed from registry

swagger:response unregisterPluginNoContent
*/
type UnregisterPluginNoContent struct {
}

// NewUnregisterPluginNoContent creates UnregisterPluginNoContent with default headers values
func NewUnregisterPluginNoContent() *UnregisterPluginNoContent {

	return &UnregisterPluginNoContent{}
}

// WriteResponse to the client
func (o *UnregisterPluginNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// UnregisterPluginNotFoundCode is the HTTP code returned for type UnregisterPluginNotFound
const UnregisterPluginNotFoundCode int = 404

/*UnregisterPluginNotFound Not Found

swagger:response unregisterPluginNotFound
*/
type UnregisterPluginNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewUnregisterPluginNotFound creates UnregisterPluginNotFound with default headers values
func NewUnregisterPluginNotFound() *UnregisterPluginNotFound {

	return &UnregisterPluginNotFound{}
}

// WithPayload adds the payload to the unregister plugin not found response
func (o *UnregisterPluginNotFound) WithPayload(payload *models.Error) *UnregisterPluginNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the unregister plugin not found response
func (o *UnregisterPluginNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UnregisterPluginNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*UnregisterPluginDefault error

swagger:response unregisterPluginDefault
*/
type UnregisterPluginDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewUnregisterPluginDefault creates UnregisterPluginDefault with default headers values
func NewUnregisterPluginDefault(code int) *UnregisterPluginDefault {
	if code <= 0 {
		code = 500
	}

	return &UnregisterPluginDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the unregister plugin default response
func (o *UnregisterPluginDefault) WithStatusCode(code int) *UnregisterPluginDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the unregister plugin default response
func (o *UnregisterPluginDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the unregister plugin default response
func (o *UnregisterPluginDefault) WithPayload(payload *models.Error) *UnregisterPluginDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the unregister plugin default response
func (o *UnregisterPluginDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UnregisterPluginDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
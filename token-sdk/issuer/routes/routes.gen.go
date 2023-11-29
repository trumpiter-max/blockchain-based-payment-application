// Package routes provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package routes

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

// Amount The amount to issue, transfer or redeem.
type Amount struct {
	// Code the code of the token
	Code string `json:"code"`

	// Value value in base units (usually cents)
	Value int64 `json:"value"`
}

// Counterparty The counterparty in a Transfer or Issuance transaction.
type Counterparty struct {
	Account string `json:"account"`

	// Node The node that holds the recipient account
	Node string `json:"node"`
}

// Error defines model for Error.
type Error struct {
	// Message High level error message
	Message string `json:"message"`

	// Payload Details about the error
	Payload string `json:"payload"`
}

// TransferRequest Instructions to issue or transfer tokens to an account
type TransferRequest struct {
	// Amount The amount to issue, transfer or redeem.
	Amount Amount `json:"amount"`

	// Counterparty The counterparty in a Transfer or Issuance transaction.
	Counterparty Counterparty `json:"counterparty"`

	// Message optional message that will be sent and stored with the transfer transaction
	Message *string `json:"message,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse = Error

// HealthSuccess defines model for HealthSuccess.
type HealthSuccess struct {
	// Message ok
	Message string `json:"message"`
}

// IssueSuccess defines model for IssueSuccess.
type IssueSuccess struct {
	Message string `json:"message"`

	// Payload Transaction id
	Payload string `json:"payload"`
}

// IssueJSONRequestBody defines body for Issue for application/json ContentType.
type IssueJSONRequestBody = TransferRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Returns 200 if the service is healthy
	// (GET /healthz)
	Healthz(ctx echo.Context) error
	// Issue tokens of any kind to an account
	// (POST /issuer/issue)
	Issue(ctx echo.Context) error
	// Returns 200 if the service is ready to accept calls
	// (GET /readyz)
	Readyz(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Healthz converts echo context to params.
func (w *ServerInterfaceWrapper) Healthz(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Healthz(ctx)
	return err
}

// Issue converts echo context to params.
func (w *ServerInterfaceWrapper) Issue(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Issue(ctx)
	return err
}

// Readyz converts echo context to params.
func (w *ServerInterfaceWrapper) Readyz(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Readyz(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/healthz", wrapper.Healthz)
	router.POST(baseURL+"/issuer/issue", wrapper.Issue)
	router.GET(baseURL+"/readyz", wrapper.Readyz)

}

type ErrorResponseJSONResponse Error

type HealthSuccessJSONResponse struct {
	// Message ok
	Message string `json:"message"`
}

type IssueSuccessJSONResponse struct {
	Message string `json:"message"`

	// Payload Transaction id
	Payload string `json:"payload"`
}

type HealthzRequestObject struct {
}

type HealthzResponseObject interface {
	VisitHealthzResponse(w http.ResponseWriter) error
}

type Healthz200JSONResponse struct{ HealthSuccessJSONResponse }

func (response Healthz200JSONResponse) VisitHealthzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Healthz503JSONResponse struct{ ErrorResponseJSONResponse }

func (response Healthz503JSONResponse) VisitHealthzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response)
}

type IssueRequestObject struct {
	Body *IssueJSONRequestBody
}

type IssueResponseObject interface {
	VisitIssueResponse(w http.ResponseWriter) error
}

type Issue200JSONResponse struct{ IssueSuccessJSONResponse }

func (response Issue200JSONResponse) VisitIssueResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type IssuedefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response IssuedefaultJSONResponse) VisitIssueResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type ReadyzRequestObject struct {
}

type ReadyzResponseObject interface {
	VisitReadyzResponse(w http.ResponseWriter) error
}

type Readyz200JSONResponse struct{ HealthSuccessJSONResponse }

func (response Readyz200JSONResponse) VisitReadyzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Readyz503JSONResponse struct{ ErrorResponseJSONResponse }

func (response Readyz503JSONResponse) VisitReadyzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Returns 200 if the service is healthy
	// (GET /healthz)
	Healthz(ctx context.Context, request HealthzRequestObject) (HealthzResponseObject, error)
	// Issue tokens of any kind to an account
	// (POST /issuer/issue)
	Issue(ctx context.Context, request IssueRequestObject) (IssueResponseObject, error)
	// Returns 200 if the service is ready to accept calls
	// (GET /readyz)
	Readyz(ctx context.Context, request ReadyzRequestObject) (ReadyzResponseObject, error)
}

type StrictHandlerFunc = strictecho.StrictEchoHandlerFunc
type StrictMiddlewareFunc = strictecho.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// Healthz operation middleware
func (sh *strictHandler) Healthz(ctx echo.Context) error {
	var request HealthzRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.Healthz(ctx.Request().Context(), request.(HealthzRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Healthz")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(HealthzResponseObject); ok {
		return validResponse.VisitHealthzResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Issue operation middleware
func (sh *strictHandler) Issue(ctx echo.Context) error {
	var request IssueRequestObject

	var body IssueJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.Issue(ctx.Request().Context(), request.(IssueRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Issue")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(IssueResponseObject); ok {
		return validResponse.VisitIssueResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Readyz operation middleware
func (sh *strictHandler) Readyz(ctx echo.Context) error {
	var request ReadyzRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.Readyz(ctx.Request().Context(), request.(ReadyzRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Readyz")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(ReadyzResponseObject); ok {
		return validResponse.VisitReadyzResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RXT2/cthP9KgP+focEIFbrOC1Q3ZI2QIxeAsdFi8Y+cKnRilmKVEjKrmrouxck9Wel",
	"1a7j1Cic02pFcubNzJvH0T3huqy0QuUsSe+JQVtpZTH8eWeMNpfdG/+Ca+VQOf/IqkoKzpzQKvlstfLv",
	"LC+wZP7p/wZzkpL/JaP1JK7aJFglbdtSkqHlRlTeCEmjO+gRkJaS98ikKz7WnKO1jwKAf7GykgF0iday",
	"LZKU6J03WhldoXEixjis3s/Q6B2hxDWVP2idEWpLPGSDX2phMCPpp+HszbBRbz4jd0vBdUFMwruwtsZv",
	"ie5oCDO8lFSskZplh+FdGaYs4/4fiOyrQx0tPiZobQBnxW1pF08I4U2p6xj4DGWBwMIaOA3C54uC89Bz",
	"NBAMZojlitD9inOdeVzvfrv8g1Byy2SNJD1brw+KHzd6pzmrpRvPTFG4AsFvBZ2Df3Z6h+owZYOreRTh",
	"NQgFG2YRaiWchRe1rZmUDXDfHC8JJbk2JfMYhHI/viaUlEKJsi5Juh5cCeVwi+agPCGQ3v9hZSj52ecQ",
	"TcWMa5bTzPd2eKwMrvby7KnKFMeY/MibWdYZ57GIhEnBPRwV66DvFJqzw9YbDizwVg2VmeP0K+AK5qDQ",
	"MrOhIAa5qAQqB73Nh/isYsL67Uspizp1TEsinw/bIiWP0Jj3YluAxFuUMLf39Z38CzompAW20bUL6Qi2",
	"nqSlKelJcIlfarQLLXqhrDN1IIQdmtRTZmjT0C5hjam9As3IMCjAqbuj04mWEj4j9KlTE/K39ITmhwcm",
	"+zpEnt0JKWGDYAPBVAbWaYMZ3AlXRD0YIh2b48H0TwKgffzLqipUrhcyH+XQt8GeKHqAURW7xK/gLeM7",
	"zGDTAINMeDyb2mEGErMtGnqtKoMWza1QW6iMuGW8gdr6f3+i0fCr0ndhK3wwWud2BVeFsPDmwwVkmAsl",
	"wiWSG62chdeQiTxH43MVbHK0FO4KwQtAxguoJIs4ul3XymgZ1DFQl2vbWIflCq7VtbrS4EwDwoGuHQWJ",
	"keAhctNRzeoSIa9VFiim1SDTtUVjV/A7czzWacOk1zB7rbbooK4y5rMQKoo4VrK7FQvh69ysRh0ULpJY",
	"u8InOjKZjtfQtRpuh4AlQ+uMbrzlcEU54byMkKuww+s1GhtrebZae2bqChWrBEnJ+Wq9Og+96YrQIAmr",
	"M+G0STq/NrkXWRvmAjTeEEk/zRnSHSGU1EaSlBTOVWmSSM2ZLLR16U/r9TphlUhuzxLS3rT0iJtkLzH2",
	"6X0WYc772xveYhABLw1h+rnwkvq+W6fT+fTVen2s9Yd9yXSGbCn5YX3+8Knp6BumlbosmWlISi7R1UZZ",
	"eLVeg4hU65gOwkKMxbe0Y1ufoDEWS268pSTw1sSfMM5puxB0IDmJooHWvdVZ82Qz+FzVF2a3ZV0f+mwi",
	"5aOuOVNj+y1lmszCAUw3kn1DpY5yMyb+BDXPJtTcL3rA1/e2zoGpBnZCZQe56Kve+br5TwBNQXiOhaFr",
	"aOIHWtbPa0EHN3pzAsyrfTB0boUzI7UNZjKmTpg5P+j/KdivEbZniDiJl8AzBj4VsTAivNjURr3saESO",
	"RfYI9X+Ohelno++kNFdLg/Nk5giFMsiy5videRmXv4crM0QSwuQcKwecSWmPXqAnBfWR0wf9d4JMnxmH",
	"uoQd2lPADTK3dIHRbgD2V1kRPhog9BAMdwclipU4Jmcpa5OSdl/n1sXXO2xiML1F7977CwP6aD64XbC+",
	"RReawA/peIummYzp3YcDl8hMjGWHWFlg/QA/OujJcejiY4c8Dm/dNxTLhPINMAIcedjetP8EAAD//+/z",
	"pA07FQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

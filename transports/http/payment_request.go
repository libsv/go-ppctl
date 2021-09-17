package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/libsv/pptcl"
)

// paymentRequestHandler is an http handler that supports BIP-270 requests.
type paymentRequestHandler struct {
	svc pptcl.PaymentRequestService
}

// NewPaymentRequestHandler will create and return a new PaymentRequestHandler.
func NewPaymentRequestHandler(svc pptcl.PaymentRequestService) *paymentRequestHandler {
	return &paymentRequestHandler{
		svc: svc,
	}
}

// RegisterRoutes will setup all routes with an echo group.
func (h *paymentRequestHandler) RegisterRoutes(g *echo.Group) {
	g.GET(RoutePaymentRequest, h.createPaymentRequest)
}

// createPaymentRequest will setup and return a new payment request.
// @Summary Request to pay an invoice and receive back outputs to use when constructing the payment transaction
// @Description Creates a payment request based on a payment id (the identifier for an invoice).
// @Tags Payment
// @Accept json
// @Produce json
// @Param paymentID path string true "Payment ID"
// @Success 201 {object} pptcl.PaymentRequest "contains outputs, merchant data and expiry information, used by the payee to construct a transaction"
// @Failure 404 {object} pptcl.ClientError "returned if the paymentID has not been found"
// @Failure 400 {object} validator.ErrValidation "returned if the user input is invalid, usually an issue with the paymentID"
// @Failure 500 {string} string "returned if there is an unexpected internal error"
// @Router /api/v1/payment/{paymentID} [GET].
func (h *paymentRequestHandler) createPaymentRequest(e echo.Context) error {
	var args pptcl.PaymentRequestArgs
	if err := e.Bind(&args); err != nil {
		return errors.Wrap(err, "failed to bind request")
	}
	resp, err := h.svc.CreatePaymentRequest(e.Request().Context(), args)
	if err != nil {
		return errors.WithStack(err)
	}
	return e.JSON(http.StatusCreated, resp)
}

package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/generated/api"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"strings"
)

func getMsg(msg []string, code int) string {
	if msg != nil && len(msg) > 0 {
		return strings.Join(msg, ". ")
	}

	return defaultStatusCodeMessages[code]
}

func Error(w http.ResponseWriter, r *http.Request, code int, err error, msg ...string) {
	ctx, span := otel.Tracer("error").Start(r.Context(), "error record")
	defer span.End()
	span.SetAttributes(attribute.String("error-from-server", err.Error()))
	span.SetAttributes(attribute.Int("http-code", code))

	r = r.WithContext(ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errMsgByte := make([]byte, 0)

	switch code {
	case http.StatusInternalServerError:
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		errMsg := api.Error{
			Message: getMsg(msg, code),
		}
		errMsgByte, err = json.Marshal(errMsg)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			w.Write([]byte(`{"error": "internal server error"}`))
			return
		}
	default:
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errMsg := api.Error400{
				Errors: make(map[string][]string),
			}

			for _, validationError := range validationErrors {
				fieldName := util.ToSnakeCase(validationError.StructField())
				errMsg.Errors[fieldName] = []string{
					getMessageForValidationError(validationError.Tag(), validationError.Param(), validationError.Value()),
				}
			}

			errMsgByte, err = json.Marshal(errMsg)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				w.Write([]byte(`{"error": "internal server error"}`))
				return
			}
		} else {
			errMsg := api.Error{
				Message: getMsg(msg, code),
			}
			errMsgByte, err = json.Marshal(errMsg)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				w.Write([]byte(`{"error": "internal server error"}`))
				return
			}
		}
	}

	span.SetAttributes(attribute.String("error-response", string(errMsgByte)))
	w.Write(errMsgByte)
}

func getMessageForValidationError(tag string, param string, value interface{}) string {
	switch tag {
	case "required":
		return fmt.Sprintf("REQUIRED")
	case "email":
		return fmt.Sprintf("INVALID EMAIL")
	case "min":
		return fmt.Sprintf("must be at least %s", param)
	case "max":
		return fmt.Sprintf("must be at most %s", param)
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", param)
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", param)
	case "eq":
		return fmt.Sprintf("must be equal to %s", param)
	case "ne":
		return fmt.Sprintf("must not be equal to %s", param)
	case "gt":
		return fmt.Sprintf("must be greater than %s", param)
	case "lt":
		return fmt.Sprintf("must be less than %s", param)
	case "oneof":
		return fmt.Sprintf("must be one of %s", param)
	case "unique":
		return fmt.Sprintf("must be unique")
	case "uuid":
		return fmt.Sprintf("must be a valid UUID")
	case "url":
		return fmt.Sprintf("must be a valid URL")
	// Custom validation
	default:
		log.Warn().Msg(fmt.Sprintf("Unknown validation tag %s", tag))
		return ""
	}
}

var defaultStatusCodeMessages = map[int]string{
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusMethodNotAllowed:    "Method Not Allowed",
	http.StatusConflict:            "Conflict",
	http.StatusInternalServerError: "Internal Status Error",
	http.StatusUnprocessableEntity: "Unprocessable Entity",
}

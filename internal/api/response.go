package api

import (
	"context"
	"encoding/json"
	"net/http"

	errs "github.com/mathif92/auth-service/internal/errors"
	"github.com/pkg/errors"
)

func Respond(ctx context.Context, w http.ResponseWriter, body interface{}, statusCode int) error {
	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonBody); err != nil {
		return err
	}

	return nil
}

func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {
	if serviceErr, ok := errors.Cause(err).(*errs.Error); ok {
		errResp := errs.ErrorResponse{Message: serviceErr.Message}

		if err := Respond(ctx, w, errResp, serviceErr.StatusCode); err != nil {
			return err
		}

		return nil
	}

	// Another kind of error was returned by the handler, so we return an InternalServerError in that case
	errResp := errs.ErrorResponse{Message: err.Error()}
	if err := Respond(ctx, w, errResp, http.StatusInternalServerError); err != nil {
		return err
	}

	return nil
}

package api

import (
	"encoding/json"
	"net/http"

	"github.com/puppetlabs/leg/encoding/transfer"
	utilapi "github.com/puppetlabs/leg/httputil/api"
	"github.com/puppetlabs/relay-core/pkg/metadataapi/errors"
	"github.com/puppetlabs/relay-core/pkg/metadataapi/server/middleware"
)

type PostEventRequestEnvelope struct {
	Data map[string]transfer.JSONInterface `json:"data"`
	Key  string                            `json:"key"`
}

func (s *Server) PostEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	managers := middleware.Managers(r)
	em := managers.Events()

	var env PostEventRequestEnvelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		utilapi.WriteError(ctx, w, errors.NewAPIMalformedRequestError().WithCause(err))
		return
	}

	data := make(map[string]interface{}, len(env.Data))
	for k, v := range env.Data {
		data[k] = v.Data
	}

	if _, err := em.Emit(ctx, data, env.Key); err != nil {
		utilapi.WriteError(ctx, w, ModelWriteError(err))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

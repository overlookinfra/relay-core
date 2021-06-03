package api

import (
	"net/http"

	"github.com/puppetlabs/leg/encoding/transfer"
	utilapi "github.com/puppetlabs/leg/httputil/api"
	"github.com/puppetlabs/relay-core/pkg/expr/evaluate"
	"github.com/puppetlabs/relay-core/pkg/expr/model"
	"github.com/puppetlabs/relay-core/pkg/manager/resolve"
	"github.com/puppetlabs/relay-core/pkg/metadataapi/errors"
	"github.com/puppetlabs/relay-core/pkg/metadataapi/server/middleware"
)

func (s *Server) GetEnvironmentVariable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	managers := middleware.Managers(r)
	name, _ := middleware.Var(r, "name")

	environment, err := managers.Environment().Get(ctx)
	if err != nil {
		utilapi.WriteError(ctx, w, ModelReadError(err))
		return
	}

	value, ok := environment.Value[name]
	if !ok {
		utilapi.WriteError(ctx, w, errors.NewModelNotFoundError())
		return
	}

	eval := evaluate.NewEvaluator(
		evaluate.WithConnectionTypeResolver{ConnectionTypeResolver: resolve.NewConnectionTypeResolver(managers.Connections())},
		evaluate.WithSecretTypeResolver{SecretTypeResolver: resolve.NewSecretTypeResolver(managers.Secrets())},
		evaluate.WithParameterTypeResolver{ParameterTypeResolver: resolve.NewParameterTypeResolver(managers.Parameters())},
		evaluate.WithOutputTypeResolver{OutputTypeResolver: resolve.NewOutputTypeResolver(managers.StepOutputs())},
	)

	rv, rerr := model.EvaluateAll(ctx, eval, value)
	if rerr != nil {
		utilapi.WriteError(ctx, w, errors.NewExpressionEvaluationError(rerr.Error()))
		return
	}

	utilapi.WriteObjectOK(ctx, w, rv)
}

func (s *Server) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	managers := middleware.Managers(r)

	environment, err := managers.Environment().Get(ctx)
	if err != nil {
		utilapi.WriteError(ctx, w, ModelReadError(err))
		return
	}

	eval := evaluate.NewEvaluator(
		evaluate.WithConnectionTypeResolver{ConnectionTypeResolver: resolve.NewConnectionTypeResolver(managers.Connections())},
		evaluate.WithSecretTypeResolver{SecretTypeResolver: resolve.NewSecretTypeResolver(managers.Secrets())},
		evaluate.WithParameterTypeResolver{ParameterTypeResolver: resolve.NewParameterTypeResolver(managers.Parameters())},
		evaluate.WithOutputTypeResolver{OutputTypeResolver: resolve.NewOutputTypeResolver(managers.StepOutputs())},
	)

	complete := true
	evs := make(map[string]interface{})
	for name, value := range environment.Value {
		rv, rerr := model.EvaluateAll(ctx, eval, value)
		if rerr != nil {
			utilapi.WriteError(ctx, w, errors.NewExpressionEvaluationError(rerr.Error()))
			return
		}

		if !rv.Complete() {
			complete = false
		}

		evs[name] = rv.Value
	}

	env := &model.JSONResultEnvelope{
		Value:    transfer.JSONInterface{Data: evs},
		Complete: complete,
	}

	utilapi.WriteObjectOK(ctx, w, env)
}

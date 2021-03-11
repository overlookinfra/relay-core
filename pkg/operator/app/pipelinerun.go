package app

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	"github.com/puppetlabs/relay-core/pkg/model"
	"github.com/puppetlabs/relay-core/pkg/obj"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ConfigurePipelineRun(ctx context.Context, pr *obj.PipelineRun, pp *PipelineParts) error {
	if err := pp.Deps.WorkflowRun.Own(ctx, pr); err != nil {
		return err
	}

	lifecycle.Label(ctx, pr, model.RelayControllerWorkflowRunIDLabel, pp.Deps.WorkflowRun.Key.Name)

	sans := make([]tektonv1beta1.PipelineRunSpecServiceAccountName, len(pp.Pipeline.Object.Spec.Tasks))
	for i, pt := range pp.Pipeline.Object.Spec.Tasks {
		sans[i] = tektonv1beta1.PipelineRunSpecServiceAccountName{
			TaskName: pt.Name,
		}
	}

	pr.Object.Spec = tektonv1beta1.PipelineRunSpec{
		ServiceAccountNames: sans,
		PipelineRef: &tektonv1beta1.PipelineRef{
			Name: pp.Pipeline.Key.Name,
		},
	}

	if pp.Deps.WorkflowRun.IsCancelled() {
		pr.Object.Spec.Status = tektonv1beta1.PipelineRunSpecStatusCancelled
	}

	return nil
}

func ApplyPipelineRun(ctx context.Context, cl client.Client, pp *PipelineParts) (*obj.PipelineRun, error) {
	pr := obj.NewPipelineRun(pp.Pipeline.Key)

	if _, err := pr.Load(ctx, cl); err != nil {
		return nil, err
	}

	pr.LabelAnnotateFrom(ctx, pp.Pipeline.Object)

	if err := ConfigurePipelineRun(ctx, pr, pp); err != nil {
		return nil, err
	}

	if err := pr.Persist(ctx, cl); err != nil {
		return nil, err
	}

	return pr, nil
}

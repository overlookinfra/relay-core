//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APITokenSource) DeepCopyInto(out *APITokenSource) {
	*out = *in
	if in.SecretKeyRef != nil {
		in, out := &in.SecretKeyRef, &out.SecretKeyRef
		*out = new(SecretKeySelector)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APITokenSource.
func (in *APITokenSource) DeepCopy() *APITokenSource {
	if in == nil {
		return nil
	}
	out := new(APITokenSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APITriggerEventSink) DeepCopyInto(out *APITriggerEventSink) {
	*out = *in
	if in.TokenFrom != nil {
		in, out := &in.TokenFrom, &out.TokenFrom
		*out = new(APITokenSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APITriggerEventSink.
func (in *APITriggerEventSink) DeepCopy() *APITriggerEventSink {
	if in == nil {
		return nil
	}
	out := new(APITriggerEventSink)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIWorkflowExecutionSink) DeepCopyInto(out *APIWorkflowExecutionSink) {
	*out = *in
	if in.TokenFrom != nil {
		in, out := &in.TokenFrom, &out.TokenFrom
		*out = new(APITokenSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIWorkflowExecutionSink.
func (in *APIWorkflowExecutionSink) DeepCopy() *APIWorkflowExecutionSink {
	if in == nil {
		return nil
	}
	out := new(APIWorkflowExecutionSink)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Container) DeepCopyInto(out *Container) {
	*out = *in
	if in.Input != nil {
		in, out := &in.Input, &out.Input
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = make(UnstructuredObject, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make(UnstructuredObject, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Container.
func (in *Container) DeepCopy() *Container {
	if in == nil {
		return nil
	}
	out := new(Container)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Decorator) DeepCopyInto(out *Decorator) {
	*out = *in
	if in.Link != nil {
		in, out := &in.Link, &out.Link
		*out = new(DecoratorLink)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Decorator.
func (in *Decorator) DeepCopy() *Decorator {
	if in == nil {
		return nil
	}
	out := new(Decorator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DecoratorLink) DeepCopyInto(out *DecoratorLink) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DecoratorLink.
func (in *DecoratorLink) DeepCopy() *DecoratorLink {
	if in == nil {
		return nil
	}
	out := new(DecoratorLink)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Log) DeepCopyInto(out *Log) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Log.
func (in *Log) DeepCopy() *Log {
	if in == nil {
		return nil
	}
	out := new(Log)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogStepMessageSource) DeepCopyInto(out *LogStepMessageSource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogStepMessageSource.
func (in *LogStepMessageSource) DeepCopy() *LogStepMessageSource {
	if in == nil {
		return nil
	}
	out := new(LogStepMessageSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceTemplate) DeepCopyInto(out *NamespaceTemplate) {
	*out = *in
	in.Metadata.DeepCopyInto(&out.Metadata)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceTemplate.
func (in *NamespaceTemplate) DeepCopy() *NamespaceTemplate {
	if in == nil {
		return nil
	}
	out := new(NamespaceTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Parameter) DeepCopyInto(out *Parameter) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Parameter.
func (in *Parameter) DeepCopy() *Parameter {
	if in == nil {
		return nil
	}
	out := new(Parameter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Run) DeepCopyInto(out *Run) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Run.
func (in *Run) DeepCopy() *Run {
	if in == nil {
		return nil
	}
	out := new(Run)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Run) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RunCondition) DeepCopyInto(out *RunCondition) {
	*out = *in
	in.Condition.DeepCopyInto(&out.Condition)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RunCondition.
func (in *RunCondition) DeepCopy() *RunCondition {
	if in == nil {
		return nil
	}
	out := new(RunCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RunList) DeepCopyInto(out *RunList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Run, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RunList.
func (in *RunList) DeepCopy() *RunList {
	if in == nil {
		return nil
	}
	out := new(RunList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RunList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RunSpec) DeepCopyInto(out *RunSpec) {
	*out = *in
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make(UnstructuredObject, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	in.State.DeepCopyInto(&out.State)
	out.WorkflowRef = in.WorkflowRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RunSpec.
func (in *RunSpec) DeepCopy() *RunSpec {
	if in == nil {
		return nil
	}
	out := new(RunSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RunState) DeepCopyInto(out *RunState) {
	*out = *in
	if in.Workflow != nil {
		in, out := &in.Workflow, &out.Workflow
		*out = make(UnstructuredObject, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make(map[string]UnstructuredObject, len(*in))
		for key, val := range *in {
			var outVal map[string]Unstructured
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make(UnstructuredObject, len(*in))
				for key, val := range *in {
					(*out)[key] = *val.DeepCopy()
				}
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RunState.
func (in *RunState) DeepCopy() *RunState {
	if in == nil {
		return nil
	}
	out := new(RunState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RunStatus) DeepCopyInto(out *RunStatus) {
	*out = *in
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]*StepStatus, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(StepStatus)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		*out = (*in).DeepCopy()
	}
	if in.CompletionTime != nil {
		in, out := &in.CompletionTime, &out.CompletionTime
		*out = (*in).DeepCopy()
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]RunCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RunStatus.
func (in *RunStatus) DeepCopy() *RunStatus {
	if in == nil {
		return nil
	}
	out := new(RunStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretKeySelector) DeepCopyInto(out *SecretKeySelector) {
	*out = *in
	out.LocalObjectReference = in.LocalObjectReference
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretKeySelector.
func (in *SecretKeySelector) DeepCopy() *SecretKeySelector {
	if in == nil {
		return nil
	}
	out := new(SecretKeySelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SpecValidationStepMessageSource) DeepCopyInto(out *SpecValidationStepMessageSource) {
	*out = *in
	if in.Expression != nil {
		in, out := &in.Expression, &out.Expression
		*out = (*in).DeepCopy()
	}
	if in.Schema != nil {
		in, out := &in.Schema, &out.Schema
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SpecValidationStepMessageSource.
func (in *SpecValidationStepMessageSource) DeepCopy() *SpecValidationStepMessageSource {
	if in == nil {
		return nil
	}
	out := new(SpecValidationStepMessageSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Step) DeepCopyInto(out *Step) {
	*out = *in
	in.Container.DeepCopyInto(&out.Container)
	if in.When != nil {
		in, out := &in.When, &out.When
		*out = (*in).DeepCopy()
	}
	if in.DependsOn != nil {
		in, out := &in.DependsOn, &out.DependsOn
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Step.
func (in *Step) DeepCopy() *Step {
	if in == nil {
		return nil
	}
	out := new(Step)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StepCondition) DeepCopyInto(out *StepCondition) {
	*out = *in
	in.Condition.DeepCopyInto(&out.Condition)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StepCondition.
func (in *StepCondition) DeepCopy() *StepCondition {
	if in == nil {
		return nil
	}
	out := new(StepCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StepMessage) DeepCopyInto(out *StepMessage) {
	*out = *in
	in.Source.DeepCopyInto(&out.Source)
	in.ObservationTime.DeepCopyInto(&out.ObservationTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StepMessage.
func (in *StepMessage) DeepCopy() *StepMessage {
	if in == nil {
		return nil
	}
	out := new(StepMessage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StepMessageSource) DeepCopyInto(out *StepMessageSource) {
	*out = *in
	if in.WhenEvaluation != nil {
		in, out := &in.WhenEvaluation, &out.WhenEvaluation
		*out = new(WhenEvaluationStepMessageSource)
		(*in).DeepCopyInto(*out)
	}
	if in.SpecValidation != nil {
		in, out := &in.SpecValidation, &out.SpecValidation
		*out = new(SpecValidationStepMessageSource)
		(*in).DeepCopyInto(*out)
	}
	if in.Log != nil {
		in, out := &in.Log, &out.Log
		*out = new(LogStepMessageSource)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StepMessageSource.
func (in *StepMessageSource) DeepCopy() *StepMessageSource {
	if in == nil {
		return nil
	}
	out := new(StepMessageSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StepOutput) DeepCopyInto(out *StepOutput) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StepOutput.
func (in *StepOutput) DeepCopy() *StepOutput {
	if in == nil {
		return nil
	}
	out := new(StepOutput)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StepStatus) DeepCopyInto(out *StepStatus) {
	*out = *in
	if in.Outputs != nil {
		in, out := &in.Outputs, &out.Outputs
		*out = make([]*StepOutput, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(StepOutput)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Decorators != nil {
		in, out := &in.Decorators, &out.Decorators
		*out = make([]*Decorator, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Decorator)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Messages != nil {
		in, out := &in.Messages, &out.Messages
		*out = make([]*StepMessage, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(StepMessage)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		*out = (*in).DeepCopy()
	}
	if in.CompletionTime != nil {
		in, out := &in.CompletionTime, &out.CompletionTime
		*out = (*in).DeepCopy()
	}
	if in.InitializationTime != nil {
		in, out := &in.InitializationTime, &out.InitializationTime
		*out = (*in).DeepCopy()
	}
	if in.Logs != nil {
		in, out := &in.Logs, &out.Logs
		*out = make([]*Log, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Log)
				**out = **in
			}
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]StepCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StepStatus.
func (in *StepStatus) DeepCopy() *StepStatus {
	if in == nil {
		return nil
	}
	out := new(StepStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tenant) DeepCopyInto(out *Tenant) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tenant.
func (in *Tenant) DeepCopy() *Tenant {
	if in == nil {
		return nil
	}
	out := new(Tenant)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Tenant) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TenantCondition) DeepCopyInto(out *TenantCondition) {
	*out = *in
	in.Condition.DeepCopyInto(&out.Condition)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TenantCondition.
func (in *TenantCondition) DeepCopy() *TenantCondition {
	if in == nil {
		return nil
	}
	out := new(TenantCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TenantList) DeepCopyInto(out *TenantList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Tenant, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TenantList.
func (in *TenantList) DeepCopy() *TenantList {
	if in == nil {
		return nil
	}
	out := new(TenantList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TenantList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TenantSpec) DeepCopyInto(out *TenantSpec) {
	*out = *in
	in.NamespaceTemplate.DeepCopyInto(&out.NamespaceTemplate)
	in.ToolInjection.DeepCopyInto(&out.ToolInjection)
	in.TriggerEventSink.DeepCopyInto(&out.TriggerEventSink)
	in.WorkflowExecutionSink.DeepCopyInto(&out.WorkflowExecutionSink)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TenantSpec.
func (in *TenantSpec) DeepCopy() *TenantSpec {
	if in == nil {
		return nil
	}
	out := new(TenantSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TenantStatus) DeepCopyInto(out *TenantStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]TenantCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TenantStatus.
func (in *TenantStatus) DeepCopy() *TenantStatus {
	if in == nil {
		return nil
	}
	out := new(TenantStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ToolInjection) DeepCopyInto(out *ToolInjection) {
	*out = *in
	if in.VolumeClaimTemplate != nil {
		in, out := &in.VolumeClaimTemplate, &out.VolumeClaimTemplate
		*out = new(v1.PersistentVolumeClaim)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ToolInjection.
func (in *ToolInjection) DeepCopy() *ToolInjection {
	if in == nil {
		return nil
	}
	out := new(ToolInjection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TriggerEventSink) DeepCopyInto(out *TriggerEventSink) {
	*out = *in
	if in.API != nil {
		in, out := &in.API, &out.API
		*out = new(APITriggerEventSink)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TriggerEventSink.
func (in *TriggerEventSink) DeepCopy() *TriggerEventSink {
	if in == nil {
		return nil
	}
	out := new(TriggerEventSink)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Unstructured) DeepCopyInto(out *Unstructured) {
	clone := in.DeepCopy()
	*out = *clone
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in UnstructuredObject) DeepCopyInto(out *UnstructuredObject) {
	{
		in := &in
		*out = make(UnstructuredObject, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UnstructuredObject.
func (in UnstructuredObject) DeepCopy() UnstructuredObject {
	if in == nil {
		return nil
	}
	out := new(UnstructuredObject)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookTrigger) DeepCopyInto(out *WebhookTrigger) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookTrigger.
func (in *WebhookTrigger) DeepCopy() *WebhookTrigger {
	if in == nil {
		return nil
	}
	out := new(WebhookTrigger)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WebhookTrigger) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookTriggerCondition) DeepCopyInto(out *WebhookTriggerCondition) {
	*out = *in
	in.Condition.DeepCopyInto(&out.Condition)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookTriggerCondition.
func (in *WebhookTriggerCondition) DeepCopy() *WebhookTriggerCondition {
	if in == nil {
		return nil
	}
	out := new(WebhookTriggerCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookTriggerList) DeepCopyInto(out *WebhookTriggerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]WebhookTrigger, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookTriggerList.
func (in *WebhookTriggerList) DeepCopy() *WebhookTriggerList {
	if in == nil {
		return nil
	}
	out := new(WebhookTriggerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WebhookTriggerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookTriggerSpec) DeepCopyInto(out *WebhookTriggerSpec) {
	*out = *in
	out.TenantRef = in.TenantRef
	in.Container.DeepCopyInto(&out.Container)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookTriggerSpec.
func (in *WebhookTriggerSpec) DeepCopy() *WebhookTriggerSpec {
	if in == nil {
		return nil
	}
	out := new(WebhookTriggerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookTriggerStatus) DeepCopyInto(out *WebhookTriggerStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]WebhookTriggerCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookTriggerStatus.
func (in *WebhookTriggerStatus) DeepCopy() *WebhookTriggerStatus {
	if in == nil {
		return nil
	}
	out := new(WebhookTriggerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WhenEvaluationStepMessageSource) DeepCopyInto(out *WhenEvaluationStepMessageSource) {
	*out = *in
	if in.Expression != nil {
		in, out := &in.Expression, &out.Expression
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WhenEvaluationStepMessageSource.
func (in *WhenEvaluationStepMessageSource) DeepCopy() *WhenEvaluationStepMessageSource {
	if in == nil {
		return nil
	}
	out := new(WhenEvaluationStepMessageSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Workflow) DeepCopyInto(out *Workflow) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Workflow.
func (in *Workflow) DeepCopy() *Workflow {
	if in == nil {
		return nil
	}
	out := new(Workflow)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Workflow) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowExecutionSink) DeepCopyInto(out *WorkflowExecutionSink) {
	*out = *in
	if in.API != nil {
		in, out := &in.API, &out.API
		*out = new(APIWorkflowExecutionSink)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowExecutionSink.
func (in *WorkflowExecutionSink) DeepCopy() *WorkflowExecutionSink {
	if in == nil {
		return nil
	}
	out := new(WorkflowExecutionSink)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowList) DeepCopyInto(out *WorkflowList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Workflow, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowList.
func (in *WorkflowList) DeepCopy() *WorkflowList {
	if in == nil {
		return nil
	}
	out := new(WorkflowList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WorkflowList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowSpec) DeepCopyInto(out *WorkflowSpec) {
	*out = *in
	out.TenantRef = in.TenantRef
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make([]*Parameter, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Parameter)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]*Step, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Step)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowSpec.
func (in *WorkflowSpec) DeepCopy() *WorkflowSpec {
	if in == nil {
		return nil
	}
	out := new(WorkflowSpec)
	in.DeepCopyInto(out)
	return out
}

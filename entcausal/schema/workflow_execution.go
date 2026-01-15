// Package schema defines the Ent schema for causal provenance tracking.
//
// WorkflowExecution represents a workflow execution step.
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// WorkflowExecution holds the schema definition for the WorkflowExecution entity.
// Represents a step in a workflow execution.
type WorkflowExecution struct {
	ent.Schema
}

// Fields of the WorkflowExecution.
func (WorkflowExecution) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			Comment("Unique execution identifier"),

		field.Time("started_at").
			Default(time.Now).
			Comment("When execution started"),

		field.Time("completed_at").
			Optional().
			Nillable().
			Comment("When execution completed"),

		field.String("workflow_id").
			NotEmpty().
			Comment("ID of the workflow definition"),

		field.String("workflow_name").
			Optional().
			Comment("Name of the workflow"),

		field.String("step_id").
			Optional().
			Comment("Current step identifier"),

		field.String("step_name").
			Optional().
			Comment("Current step name"),

		field.Int("step_index").
			Default(0).
			Comment("Index of current step"),

		field.Enum("status").
			Values("pending", "running", "completed", "failed", "cancelled", "paused").
			Default("pending").
			Comment("Current execution status"),

		field.JSON("inputs", map[string]interface{}{}).
			Optional().
			Comment("Workflow inputs"),

		field.JSON("outputs", map[string]interface{}{}).
			Optional().
			Comment("Workflow outputs"),

		field.String("error").
			Optional().
			Comment("Error message if failed"),

		field.Float("duration_ms").
			Optional().
			Comment("Total duration in milliseconds"),

		field.String("parent_execution_id").
			Optional().
			Comment("Parent execution for nested workflows"),

		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Additional metadata"),
	}
}

// Edges of the WorkflowExecution.
func (WorkflowExecution) Edges() []ent.Edge {
	return []ent.Edge{
		// WorkflowExecution is executed by AgentActions
		edge.From("actions", AgentAction.Type).
			Ref("workflows").
			Comment("Agent actions that executed this workflow"),

		// WorkflowExecution produces ExternalOutputs
		edge.To("outputs", ExternalOutput.Type).
			Comment("External outputs produced by this workflow"),

		// Self-referential edge for nested workflows
		edge.To("child_executions", WorkflowExecution.Type).
			From("parent_execution").
			Unique().
			Comment("Child workflow executions"),
	}
}

// Indexes of the WorkflowExecution.
func (WorkflowExecution) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("workflow_id"),
		index.Fields("status"),
		index.Fields("started_at"),
		index.Fields("parent_execution_id"),
	}
}

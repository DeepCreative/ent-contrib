// Package schema defines the Ent schema for causal provenance tracking.
//
// AgentAction represents an action taken by an ARIA/PERSONA agent.
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// AgentAction holds the schema definition for the AgentAction entity.
// Represents an action taken by an AI agent (ARIA, PERSONA, etc.).
type AgentAction struct {
	ent.Schema
}

// Fields of the AgentAction.
func (AgentAction) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			Comment("Unique action identifier"),

		field.Time("timestamp").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when action was taken"),

		field.String("agent_id").
			NotEmpty().
			Comment("ID of the agent that took the action"),

		field.String("agent_type").
			NotEmpty().
			Comment("Type of agent (aria, persona, conductor, etc.)"),

		field.String("action_type").
			NotEmpty().
			Comment("Type of action taken"),

		field.String("action_name").
			Optional().
			Comment("Human-readable action name"),

		field.JSON("parameters", map[string]interface{}{}).
			Optional().
			Comment("Action parameters"),

		field.String("target_resource").
			Optional().
			Comment("Resource the action targets"),

		field.Enum("status").
			Values("pending", "executing", "completed", "failed", "cancelled").
			Default("pending").
			Comment("Current status of the action"),

		field.String("result").
			Optional().
			Comment("Result of the action"),

		field.String("error").
			Optional().
			Comment("Error message if failed"),

		field.Float("latency_ms").
			Optional().
			Comment("Execution latency in milliseconds"),

		field.String("session_id").
			Optional().
			Comment("Session identifier"),

		field.String("user_id").
			Optional().
			Comment("User who initiated the action"),

		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Additional metadata"),
	}
}

// Edges of the AgentAction.
func (AgentAction) Edges() []ent.Edge {
	return []ent.Edge{
		// AgentAction is triggered by RoutingDecisions
		edge.From("decisions", RoutingDecision.Type).
			Ref("actions").
			Comment("Routing decisions that triggered this action"),

		// AgentAction executes WorkflowExecutions
		edge.To("workflows", WorkflowExecution.Type).
			Comment("Workflow executions performed by this action"),
	}
}

// Indexes of the AgentAction.
func (AgentAction) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("agent_id"),
		index.Fields("agent_type"),
		index.Fields("action_type"),
		index.Fields("timestamp"),
		index.Fields("status"),
		index.Fields("session_id"),
		index.Fields("user_id"),
	}
}

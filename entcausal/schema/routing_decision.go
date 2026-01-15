// Package schema defines the Ent schema for causal provenance tracking.
//
// RoutingDecision represents a BMU routing decision.
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RoutingDecision holds the schema definition for the RoutingDecision entity.
// Represents a decision made by the Bellman Memory Unit.
type RoutingDecision struct {
	ent.Schema
}

// Fields of the RoutingDecision.
func (RoutingDecision) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			Comment("Unique decision identifier"),

		field.Time("timestamp").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when decision was made"),

		field.String("inference_id").
			NotEmpty().
			Comment("ID of the inference this decision belongs to"),

		field.Enum("decision_type").
			Values("exit", "skip", "route", "escalate", "iterate").
			Comment("Type of routing decision"),

		field.Int("layer_index").
			Default(0).
			Comment("Layer where decision was made"),

		field.Float("gate_probability").
			Default(0.0).
			Min(0.0).
			Max(1.0).
			Comment("Probability from gate network"),

		field.String("selected_model").
			Optional().
			Comment("Model selected for routing"),

		field.Int("iteration_count").
			Default(0).
			Comment("Current iteration number"),

		field.Float("confidence").
			Default(0.0).
			Min(0.0).
			Max(1.0).
			Comment("Confidence in the decision"),

		field.String("domain").
			Optional().
			Comment("Domain context for the decision"),

		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Additional metadata"),
	}
}

// Edges of the RoutingDecision.
func (RoutingDecision) Edges() []ent.Edge {
	return []ent.Edge{
		// RoutingDecision is caused by SpikeEvents
		edge.From("spike_events", SpikeEvent.Type).
			Ref("decisions").
			Comment("Spike events that caused this decision"),

		// RoutingDecision triggers AgentActions
		edge.To("actions", AgentAction.Type).
			Comment("Agent actions triggered by this decision"),
	}
}

// Indexes of the RoutingDecision.
func (RoutingDecision) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("inference_id"),
		index.Fields("decision_type"),
		index.Fields("selected_model"),
		index.Fields("timestamp"),
		index.Fields("domain"),
	}
}

// Package schema defines the Ent schema for causal provenance tracking.
//
// SpikeEvent represents a neural spike event from the D3N spiking fabric.
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SpikeEvent holds the schema definition for the SpikeEvent entity.
// Represents a spike event in the D3N spiking neural fabric.
type SpikeEvent struct {
	ent.Schema
}

// Fields of the SpikeEvent.
func (SpikeEvent) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			Comment("Unique spike event identifier"),

		field.Time("timestamp").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when spike occurred"),

		field.Int64("timestamp_ns").
			Optional().
			Comment("Nanosecond precision timestamp"),

		field.String("population_id").
			NotEmpty().
			Comment("ID of the neuron population that fired"),

		field.Int("layer_index").
			Default(0).
			Comment("Layer index in the neural network"),

		field.JSON("neuron_indices", []int{}).
			Comment("Indices of neurons that fired"),

		field.JSON("spike_counts", []int{}).
			Optional().
			Comment("Number of spikes per neuron"),

		field.JSON("membrane_potentials", []float64{}).
			Optional().
			Comment("Membrane potentials at spike time"),

		field.JSON("input_currents", []float64{}).
			Optional().
			Comment("Input currents at spike time"),

		field.String("pattern_hash").
			NotEmpty().
			Comment("Deterministic hash of the firing pattern"),

		field.String("inference_id").
			Optional().
			Comment("Links to the neural audit entry"),

		field.Bool("is_emergent").
			Default(false).
			Comment("Whether this is an emergent pattern"),

		field.Float("entropy").
			Default(0.0).
			Comment("Entropy of the spike pattern"),

		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Additional metadata"),
	}
}

// Edges of the SpikeEvent.
func (SpikeEvent) Edges() []ent.Edge {
	return []ent.Edge{
		// SpikeEvent causes RoutingDecision
		edge.To("decisions", RoutingDecision.Type).
			Comment("Routing decisions caused by this spike event"),
	}
}

// Indexes of the SpikeEvent.
func (SpikeEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("pattern_hash"),
		index.Fields("inference_id"),
		index.Fields("population_id"),
		index.Fields("timestamp"),
		index.Fields("is_emergent"),
	}
}

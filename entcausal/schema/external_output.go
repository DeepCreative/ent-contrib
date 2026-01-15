// Package schema defines the Ent schema for causal provenance tracking.
//
// ExternalOutput represents an external output (Foundation TX, AWS record, trade, etc.).
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ExternalOutput holds the schema definition for the ExternalOutput entity.
// Represents an external output produced by the system.
type ExternalOutput struct {
	ent.Schema
}

// Fields of the ExternalOutput.
func (ExternalOutput) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			Comment("Unique output identifier"),

		field.Time("timestamp").
			Default(time.Now).
			Immutable().
			Comment("When output was produced"),

		field.Enum("output_type").
			Values(
				"foundation_tx",   // Hyperledger Fabric transaction
				"aws_record",      // AWS database record
				"trade_execution", // Trading execution
				"document",        // Generated document
				"api_response",    // External API response
				"notification",    // Alert/notification
				"file",            // File output
				"other",           // Other output type
			).
			Comment("Type of external output"),

		field.String("destination").
			Optional().
			Comment("Destination system/service"),

		field.String("destination_id").
			Optional().
			Comment("ID in the destination system"),

		field.String("transaction_id").
			Optional().
			Comment("Transaction ID for blockchain outputs"),

		field.String("block_hash").
			Optional().
			Comment("Block hash for blockchain outputs"),

		field.Int64("block_number").
			Optional().
			Comment("Block number for blockchain outputs"),

		field.String("content_hash").
			NotEmpty().
			Comment("Hash of the output content"),

		field.Int64("content_size").
			Optional().
			Comment("Size of the output in bytes"),

		field.Enum("status").
			Values("pending", "confirmed", "failed", "reverted").
			Default("pending").
			Comment("Status of the output"),

		field.String("domain").
			Optional().
			Comment("Domain context (trading, medical, etc.)"),

		field.JSON("compliance", map[string]interface{}{}).
			Optional().
			Comment("Compliance metadata (SEC, HIPAA, etc.)"),

		field.Int("retention_years").
			Default(7).
			Comment("Retention period in years"),

		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Additional metadata"),
	}
}

// Edges of the ExternalOutput.
func (ExternalOutput) Edges() []ent.Edge {
	return []ent.Edge{
		// ExternalOutput is produced by WorkflowExecutions
		edge.From("workflows", WorkflowExecution.Type).
			Ref("outputs").
			Comment("Workflow executions that produced this output"),
	}
}

// Indexes of the ExternalOutput.
func (ExternalOutput) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("output_type"),
		index.Fields("destination"),
		index.Fields("transaction_id"),
		index.Fields("content_hash"),
		index.Fields("status"),
		index.Fields("domain"),
		index.Fields("timestamp"),
	}
}

// Package queries provides causal provenance graph queries.
//
// Implements efficient traversal of the causal graph from outputs
// back to neural spike events.
package queries

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
)

// CausalPath represents a path through the causal graph.
type CausalPath struct {
	OutputID       string           `json:"output_id"`
	Nodes          []CausalNode     `json:"nodes"`
	Edges          []CausalEdge     `json:"edges"`
	Depth          int              `json:"depth"`
	TotalLatencyMs float64          `json:"total_latency_ms"`
	TracedAt       time.Time        `json:"traced_at"`
}

// CausalNode represents a node in the causal path.
type CausalNode struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Depth     int                    `json:"depth"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// CausalEdge represents an edge in the causal path.
type CausalEdge struct {
	SourceID   string  `json:"source_id"`
	SourceType string  `json:"source_type"`
	TargetID   string  `json:"target_id"`
	TargetType string  `json:"target_type"`
	EdgeType   string  `json:"edge_type"`
	Confidence float64 `json:"confidence"`
}

// EmergentPatternResult represents a detected emergent pattern.
type EmergentPatternResult struct {
	PatternHash     string    `json:"pattern_hash"`
	OccurrenceCount int       `json:"occurrence_count"`
	NeuronIndices   []int     `json:"neuron_indices"`
	PopulationID    string    `json:"population_id"`
	FirstSeen       time.Time `json:"first_seen"`
	LastSeen        time.Time `json:"last_seen"`
	Significance    float64   `json:"significance"`
}

// AgentDecisionPath represents the full path of an agent's decision.
type AgentDecisionPath struct {
	AgentID      string       `json:"agent_id"`
	ActionID     string       `json:"action_id"`
	SpikeEvents  []CausalNode `json:"spike_events"`
	Decisions    []CausalNode `json:"decisions"`
	Workflows    []CausalNode `json:"workflows"`
	Outputs      []CausalNode `json:"outputs"`
	TotalDepth   int          `json:"total_depth"`
}

// CausalQueryService provides causal graph query operations.
type CausalQueryService struct {
	// client would be the generated ent.Client
	// For this example, we use interface{} as placeholder
	client interface{}
}

// NewCausalQueryService creates a new causal query service.
func NewCausalQueryService(client interface{}) *CausalQueryService {
	return &CausalQueryService{client: client}
}

// TraceCausality traces the causal chain from an output back to spike events.
//
// This performs a breadth-first traversal of the causal graph, following
// edges backwards from the output to find all contributing spike events.
//
// Example:
//
//	path, err := service.TraceCausality(ctx, "output-123", 100)
//	spikeEvents := path.GetSpikeEvents()
func (s *CausalQueryService) TraceCausality(
	ctx context.Context,
	outputID string,
	maxDepth int,
) (*CausalPath, error) {
	if maxDepth <= 0 {
		maxDepth = 100
	}

	path := &CausalPath{
		OutputID: outputID,
		Nodes:    make([]CausalNode, 0),
		Edges:    make([]CausalEdge, 0),
		TracedAt: time.Now(),
	}

	// Start with the output node
	visited := make(map[string]bool)
	queue := []struct {
		id       string
		nodeType string
		depth    int
	}{{outputID, "external_output", 0}}

	for len(queue) > 0 && path.Depth <= maxDepth {
		current := queue[0]
		queue = queue[1:]

		if visited[current.id] {
			continue
		}
		visited[current.id] = true

		// Add node to path
		node := CausalNode{
			ID:        current.id,
			Type:      current.nodeType,
			Timestamp: time.Now(), // Would be fetched from DB
			Depth:     current.depth,
		}
		path.Nodes = append(path.Nodes, node)
		path.Depth = max(path.Depth, current.depth)

		// Get parent nodes based on type
		parents, edges := s.getParentNodes(ctx, current.id, current.nodeType)
		for _, edge := range edges {
			path.Edges = append(path.Edges, edge)
		}
		for _, parent := range parents {
			if !visited[parent.ID] {
				queue = append(queue, struct {
					id       string
					nodeType string
					depth    int
				}{parent.ID, parent.Type, current.depth + 1})
			}
		}
	}

	return path, nil
}

// getParentNodes returns parent nodes for a given node.
// This is a placeholder - actual implementation would query the database.
func (s *CausalQueryService) getParentNodes(
	ctx context.Context,
	nodeID string,
	nodeType string,
) ([]CausalNode, []CausalEdge) {
	// In actual implementation, this would query the ent client
	// based on the node type to find parent edges
	//
	// For example:
	// switch nodeType {
	// case "external_output":
	//     workflows := client.ExternalOutput.Query().
	//         Where(externaloutput.ID(nodeID)).
	//         QueryWorkflows().
	//         AllX(ctx)
	// case "workflow":
	//     actions := client.WorkflowExecution.Query().
	//         Where(workflowexecution.ID(nodeID)).
	//         QueryActions().
	//         AllX(ctx)
	// ...
	// }

	return nil, nil
}

// FindEmergentPatterns finds emergent spike patterns in a time range.
//
// Emergent patterns are groups of neurons that fire together more
// frequently than expected by chance.
func (s *CausalQueryService) FindEmergentPatterns(
	ctx context.Context,
	startTime time.Time,
	endTime time.Time,
	minOccurrences int,
) ([]EmergentPatternResult, error) {
	if minOccurrences <= 0 {
		minOccurrences = 5
	}

	// In actual implementation, this would:
	// 1. Query spike events in the time range
	// 2. Group by pattern_hash
	// 3. Filter by occurrence count
	// 4. Calculate significance scores
	//
	// Example query:
	// SELECT pattern_hash, COUNT(*) as count, 
	//        MIN(timestamp) as first_seen, MAX(timestamp) as last_seen,
	//        population_id, neuron_indices
	// FROM spike_events
	// WHERE timestamp BETWEEN ? AND ?
	// GROUP BY pattern_hash
	// HAVING COUNT(*) >= ?
	// ORDER BY COUNT(*) DESC

	results := make([]EmergentPatternResult, 0)
	return results, nil
}

// GetAgentDecisionPath gets the full provenance path for an agent action.
//
// Returns all spike events, decisions, workflows, and outputs
// that are causally related to the given action.
func (s *CausalQueryService) GetAgentDecisionPath(
	ctx context.Context,
	agentID string,
	actionID string,
) (*AgentDecisionPath, error) {
	path := &AgentDecisionPath{
		AgentID:     agentID,
		ActionID:    actionID,
		SpikeEvents: make([]CausalNode, 0),
		Decisions:   make([]CausalNode, 0),
		Workflows:   make([]CausalNode, 0),
		Outputs:     make([]CausalNode, 0),
	}

	// In actual implementation:
	// 1. Get the agent action
	// 2. Traverse backwards to decisions and spike events
	// 3. Traverse forwards to workflows and outputs
	//
	// action := client.AgentAction.Query().
	//     Where(agentaction.ID(actionID)).
	//     WithDecisions(func(q *ent.RoutingDecisionQuery) {
	//         q.WithSpikeEvents()
	//     }).
	//     WithWorkflows(func(q *ent.WorkflowExecutionQuery) {
	//         q.WithOutputs()
	//     }).
	//     OnlyX(ctx)

	return path, nil
}

// QueryByPatternHash finds all entries related to a specific spike pattern.
func (s *CausalQueryService) QueryByPatternHash(
	ctx context.Context,
	patternHash string,
	limit int,
) ([]CausalNode, error) {
	if limit <= 0 {
		limit = 100
	}

	// In actual implementation:
	// return client.SpikeEvent.Query().
	//     Where(spikeevent.PatternHash(patternHash)).
	//     Limit(limit).
	//     AllX(ctx)

	return nil, nil
}

// QueryByInferenceID finds all causal nodes related to an inference.
func (s *CausalQueryService) QueryByInferenceID(
	ctx context.Context,
	inferenceID string,
) (*CausalPath, error) {
	// Get spike events for inference
	// Get decisions for inference
	// Get actions triggered by those decisions
	// Get workflows executed by those actions
	// Get outputs produced by those workflows

	return s.buildInferencePath(ctx, inferenceID)
}

// buildInferencePath builds a causal path for an inference.
func (s *CausalQueryService) buildInferencePath(
	ctx context.Context,
	inferenceID string,
) (*CausalPath, error) {
	path := &CausalPath{
		OutputID: inferenceID,
		Nodes:    make([]CausalNode, 0),
		Edges:    make([]CausalEdge, 0),
		TracedAt: time.Now(),
	}

	// Placeholder implementation
	return path, nil
}

// GetSpikeEvents returns all spike event nodes from a CausalPath.
func (p *CausalPath) GetSpikeEvents() []CausalNode {
	events := make([]CausalNode, 0)
	for _, node := range p.Nodes {
		if node.Type == "spike_event" {
			events = append(events, node)
		}
	}
	return events
}

// GetDecisions returns all routing decision nodes from a CausalPath.
func (p *CausalPath) GetDecisions() []CausalNode {
	decisions := make([]CausalNode, 0)
	for _, node := range p.Nodes {
		if node.Type == "routing_decision" {
			decisions = append(decisions, node)
		}
	}
	return decisions
}

// CountByType returns counts of nodes by type.
func (p *CausalPath) CountByType() map[string]int {
	counts := make(map[string]int)
	for _, node := range p.Nodes {
		counts[node.Type]++
	}
	return counts
}

// Helper functions

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SQL query helpers for raw queries if needed

// TraceCausalitySQL returns the SQL for tracing causality.
// This can be used with raw SQL queries for performance.
func TraceCausalitySQL(outputID string, maxDepth int) string {
	return fmt.Sprintf(`
		WITH RECURSIVE causal_chain AS (
			-- Base case: start with the output
			SELECT 
				id, 
				'external_output' as node_type,
				0 as depth
			FROM external_outputs
			WHERE id = '%s'
			
			UNION ALL
			
			-- Recursive case: follow edges backwards
			SELECT 
				CASE 
					WHEN cc.node_type = 'external_output' THEN we.id
					WHEN cc.node_type = 'workflow_execution' THEN aa.id
					WHEN cc.node_type = 'agent_action' THEN rd.id
					WHEN cc.node_type = 'routing_decision' THEN se.id
				END as id,
				CASE 
					WHEN cc.node_type = 'external_output' THEN 'workflow_execution'
					WHEN cc.node_type = 'workflow_execution' THEN 'agent_action'
					WHEN cc.node_type = 'agent_action' THEN 'routing_decision'
					WHEN cc.node_type = 'routing_decision' THEN 'spike_event'
				END as node_type,
				cc.depth + 1 as depth
			FROM causal_chain cc
			LEFT JOIN workflow_execution_outputs weo ON cc.id = weo.external_output_id
			LEFT JOIN workflow_executions we ON weo.workflow_execution_id = we.id
			LEFT JOIN agent_action_workflows aaw ON we.id = aaw.workflow_execution_id
			LEFT JOIN agent_actions aa ON aaw.agent_action_id = aa.id
			LEFT JOIN routing_decision_actions rda ON aa.id = rda.agent_action_id
			LEFT JOIN routing_decisions rd ON rda.routing_decision_id = rd.id
			LEFT JOIN spike_event_decisions sed ON rd.id = sed.routing_decision_id
			LEFT JOIN spike_events se ON sed.spike_event_id = se.id
			WHERE cc.depth < %d
		)
		SELECT DISTINCT id, node_type, depth
		FROM causal_chain
		WHERE id IS NOT NULL
		ORDER BY depth, node_type
	`, outputID, maxDepth)
}

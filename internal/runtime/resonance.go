package runtime

// resonance.go provides resonance graph and field coupling for String Theory semantics (M2.4).
// Resonance describes how changes in one string field propagate to other coupled fields.

// ResonanceGraph represents the coupling relationships between string fields.
// When a field value changes, the graph determines which other fields are affected
// and by how much (coupling strength).
type ResonanceGraph struct {
	links map[string][]*ResonanceLink // Source field -> list of coupled fields
}

// ResonanceLink represents a directed coupling from one field to another.
type ResonanceLink struct {
	SourceString string  // Source string name
	SourceMode   string  // Source mode name
	TargetString string  // Target string name
	TargetMode   string  // Target mode name
	Strength     float64 // Coupling strength (0.0 to 1.0)
	CouplingType CouplingType
}

// CouplingType represents the type of coupling between fields.
type CouplingType int

const (
	CouplingLinear CouplingType = iota // Linear coupling: Δtarget = strength * Δsource
	CouplingNonlinear                  // Nonlinear coupling (future)
	CouplingResonant                   // Resonant coupling (frequency-dependent)
)

// NewResonanceGraph creates a new empty resonance graph.
func NewResonanceGraph() *ResonanceGraph {
	return &ResonanceGraph{
		links: make(map[string][]*ResonanceLink),
	}
}

// AddLink adds a resonance link to the graph.
func (rg *ResonanceGraph) AddLink(link *ResonanceLink) {
	key := makeFieldKey(link.SourceString, link.SourceMode)
	rg.links[key] = append(rg.links[key], link)
}

// GetLinks returns all resonance links from a given source field.
func (rg *ResonanceGraph) GetLinks(stringName, modeName string) []*ResonanceLink {
	key := makeFieldKey(stringName, modeName)
	return rg.links[key]
}

// PropagateChange propagates a value change through the resonance graph.
// Returns a map of affected fields and their delta values.
func (rg *ResonanceGraph) PropagateChange(stringName, modeName string, delta float64) map[string]float64 {
	affected := make(map[string]float64)
	
	// Get all links from this source field
	links := rg.GetLinks(stringName, modeName)
	
	for _, link := range links {
		targetKey := makeFieldKey(link.TargetString, link.TargetMode)
		
		// Calculate the propagated change based on coupling type
		var propagatedDelta float64
		switch link.CouplingType {
		case CouplingLinear:
			propagatedDelta = delta * link.Strength
		case CouplingNonlinear:
			// TODO(M2.4): Implement nonlinear coupling
			propagatedDelta = delta * link.Strength
		case CouplingResonant:
			// TODO(M2.4): Implement frequency-dependent resonance
			propagatedDelta = delta * link.Strength
		}
		
		affected[targetKey] = propagatedDelta
	}
	
	return affected
}

// makeFieldKey creates a unique key for a string field.
func makeFieldKey(stringName, modeName string) string {
	return stringName + "." + modeName
}

// ==================== Resonance Conditions ====================

// ResonanceCondition represents a condition that triggers when resonance occurs.
// This is used for the `resonance when ...` syntax in M2.4.
type ResonanceCondition struct {
	Predicate func(state *StringState) bool // Condition to check
	Action    func(state *StringState)      // Action to execute when condition is true
}

// CheckCondition evaluates the resonance condition.
func (rc *ResonanceCondition) CheckCondition(state *StringState) bool {
	return rc.Predicate(state)
}

// ExecuteAction executes the resonance action.
func (rc *ResonanceCondition) ExecuteAction(state *StringState) {
	rc.Action(state)
}

// ==================== String State (M2.4) ====================

// StringState represents the runtime state of a string field.
// Each string has multiple modes (vibrational degrees of freedom).
type StringState struct {
	Name  string             // String name
	Modes map[string]float64 // Mode name -> current value
}

// NewStringState creates a new string state.
func NewStringState(name string) *StringState {
	return &StringState{
		Name:  name,
		Modes: make(map[string]float64),
	}
}

// SetMode sets the value of a mode.
func (ss *StringState) SetMode(modeName string, value float64) {
	ss.Modes[modeName] = value
}

// GetMode gets the value of a mode.
func (ss *StringState) GetMode(modeName string) (float64, bool) {
	value, exists := ss.Modes[modeName]
	return value, exists
}

// UpdateMode updates a mode value and returns the delta.
func (ss *StringState) UpdateMode(modeName string, delta float64) float64 {
	oldValue := ss.Modes[modeName]
	newValue := oldValue + delta
	ss.Modes[modeName] = newValue
	return delta
}

// ==================== Brane Context (M2.4) ====================

// BraneContext represents the dimensional context for strings.
// Branes provide the background space in which strings exist.
type BraneContext struct {
	Name       string   // Brane name
	Dimensions []string // Dimension names
	Strings    []string // Strings attached to this brane
}

// NewBraneContext creates a new brane context.
func NewBraneContext(name string, dimensions []string) *BraneContext {
	return &BraneContext{
		Name:       name,
		Dimensions: dimensions,
		Strings:    []string{},
	}
}

// AttachString attaches a string to this brane.
func (bc *BraneContext) AttachString(stringName string) {
	bc.Strings = append(bc.Strings, stringName)
}

// IsStringAttached checks if a string is attached to this brane.
func (bc *BraneContext) IsStringAttached(stringName string) bool {
	for _, s := range bc.Strings {
		if s == stringName {
			return true
		}
	}
	return false
}

// GetDimensionCount returns the number of dimensions in this brane.
func (bc *BraneContext) GetDimensionCount() int {
	return len(bc.Dimensions)
}

// ==================== Future Extensions ====================

// TODO(M2.4): Full String Theory Implementation
// - Frequency analysis and matching
// - Amplitude thresholds for resonance triggering
// - Phase relationships between modes
// - Damping and feedback mechanisms
// - Multi-string resonance (more than two strings)
// - Brane geometry and metric

// TODO(M2.4.3): Integration with Runtime
// - Hook resonance updates into variable assignment
// - Automatic propagation when string modes change
// - Resonance condition checking in event loop
// - Visualization of resonance graph


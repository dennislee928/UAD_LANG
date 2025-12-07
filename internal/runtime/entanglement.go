package runtime

import "fmt"

// entanglement.go provides quantum entanglement semantics for variable synchronization (M2.5).
// Entangled variables share the same underlying value: when one changes, all change.

// EntanglementGroup represents a group of entangled variables.
// All variables in a group share the same backing value.
type EntanglementGroup struct {
	ID      string   // Unique group identifier
	Members []string // Variable names in this group
	Value   Value    // Shared backing value
}

// NewEntanglementGroup creates a new entanglement group.
func NewEntanglementGroup(id string) *EntanglementGroup {
	return &EntanglementGroup{
		ID:      id,
		Members: []string{},
		Value:   nil,
	}
}

// AddMember adds a variable to the entanglement group.
func (eg *EntanglementGroup) AddMember(varName string) {
	eg.Members = append(eg.Members, varName)
}

// SetValue sets the shared value for all members.
func (eg *EntanglementGroup) SetValue(value Value) {
	eg.Value = value
}

// GetValue returns the shared value.
func (eg *EntanglementGroup) GetValue() Value {
	return eg.Value
}

// IsMember checks if a variable is a member of this group.
func (eg *EntanglementGroup) IsMember(varName string) bool {
	for _, member := range eg.Members {
		if member == varName {
			return true
		}
	}
	return false
}

// ==================== Entanglement Manager ====================

// EntanglementManager manages all entanglement groups in a program.
type EntanglementManager struct {
	groups       map[string]*EntanglementGroup // Group ID -> Group
	varToGroup   map[string]string             // Variable name -> Group ID
	nextGroupID  int                           // Counter for generating unique group IDs
}

// NewEntanglementManager creates a new entanglement manager.
func NewEntanglementManager() *EntanglementManager {
	return &EntanglementManager{
		groups:      make(map[string]*EntanglementGroup),
		varToGroup:  make(map[string]string),
		nextGroupID: 0,
	}
}

// CreateGroup creates a new entanglement group.
func (em *EntanglementManager) CreateGroup() *EntanglementGroup {
	id := fmt.Sprintf("entangle_%d", em.nextGroupID)
	em.nextGroupID++
	
	group := NewEntanglementGroup(id)
	em.groups[id] = group
	
	return group
}

// EntangleVariables creates an entanglement group for the given variables.
func (em *EntanglementManager) EntangleVariables(varNames []string) (*EntanglementGroup, error) {
	// Check if any variable is already entangled
	for _, varName := range varNames {
		if _, exists := em.varToGroup[varName]; exists {
			return nil, fmt.Errorf("variable '%s' is already entangled", varName)
		}
	}
	
	// Create new group
	group := em.CreateGroup()
	
	// Add all variables to the group
	for _, varName := range varNames {
		group.AddMember(varName)
		em.varToGroup[varName] = group.ID
	}
	
	return group, nil
}

// GetGroup returns the entanglement group for a variable.
func (em *EntanglementManager) GetGroup(varName string) (*EntanglementGroup, bool) {
	groupID, exists := em.varToGroup[varName]
	if !exists {
		return nil, false
	}
	
	group, exists := em.groups[groupID]
	return group, exists
}

// IsEntangled checks if a variable is entangled.
func (em *EntanglementManager) IsEntangled(varName string) bool {
	_, exists := em.varToGroup[varName]
	return exists
}

// SetEntangledValue sets the value for all variables in an entanglement group.
// This is called when any entangled variable is assigned a new value.
func (em *EntanglementManager) SetEntangledValue(varName string, value Value, env *Environment) error {
	group, exists := em.GetGroup(varName)
	if !exists {
		return fmt.Errorf("variable '%s' is not entangled", varName)
	}
	
	// Update the shared value
	group.SetValue(value)
	
	// Update all member variables in the environment
	for _, member := range group.Members {
		if err := env.Set(member, value); err != nil {
			return err
		}
	}
	
	return nil
}

// GetEntangledValue gets the shared value for an entangled variable.
func (em *EntanglementManager) GetEntangledValue(varName string) (Value, bool) {
	group, exists := em.GetGroup(varName)
	if !exists {
		return nil, false
	}
	
	return group.GetValue(), true
}

// DisentangleVariable removes a variable from its entanglement group.
func (em *EntanglementManager) DisentangleVariable(varName string) error {
	groupID, exists := em.varToGroup[varName]
	if !exists {
		return fmt.Errorf("variable '%s' is not entangled", varName)
	}
	
	group := em.groups[groupID]
	
	// Remove from group members
	newMembers := []string{}
	for _, member := range group.Members {
		if member != varName {
			newMembers = append(newMembers, member)
		}
	}
	group.Members = newMembers
	
	// Remove from mapping
	delete(em.varToGroup, varName)
	
	// If group is now empty, remove it
	if len(group.Members) == 0 {
		delete(em.groups, groupID)
	}
	
	return nil
}

// ListGroups returns all entanglement groups.
func (em *EntanglementManager) ListGroups() []*EntanglementGroup {
	groups := make([]*EntanglementGroup, 0, len(em.groups))
	for _, group := range em.groups {
		groups = append(groups, group)
	}
	return groups
}

// ==================== Type Compatibility Checking ====================

// CheckTypeCompatibility checks if all variables in an entanglement group have compatible types.
// This should be called during semantic analysis (M2.5.2).
func CheckTypeCompatibility(varTypes map[string]string) error {
	if len(varTypes) == 0 {
		return nil
	}
	
	// Get the first type as reference
	var referenceType string
	for _, typ := range varTypes {
		referenceType = typ
		break
	}
	
	// Check all other types match
	for varName, typ := range varTypes {
		if typ != referenceType {
			return fmt.Errorf("entangled variable '%s' has incompatible type '%s' (expected '%s')", 
				varName, typ, referenceType)
		}
	}
	
	return nil
}

// ==================== Future Extensions ====================

// TODO(M2.5.1): AST Integration
// - Parse `entangle x, y, z` statements
// - Support entangle expressions in various contexts

// TODO(M2.5.2): Semantic Analysis
// - Type checking for entangled variables
// - Scope analysis (ensure all variables are in scope)
// - Lifetime analysis (entanglement duration)

// TODO(M2.5.3): Runtime Integration
// - Hook into variable assignment in interpreter
// - Automatic synchronization on write
// - Support for entanglement in VM bytecode
// - Visualization of entanglement relationships

// TODO(Advanced): Partial Entanglement
// - Entangle only specific fields of structs
// - Conditional entanglement (entangle when condition is true)
// - Temporal entanglement (entangle for a specific duration)


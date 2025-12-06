package runtime

// temporal.go provides temporal scheduling and time-based execution concepts.
// This module supports the Musical DSL (M2.3) by providing a temporal grid
// that allows precise control over when events occur in time.

// TemporalGrid represents a time-structured execution environment.
// It divides time into discrete units (beats, bars) and schedules events accordingly.
// This is the foundation for Musical DSL execution.
type TemporalGrid struct {
	CurrentTick int // Current time tick (finest granularity)
	CurrentBeat int // Current beat number
	CurrentBar  int // Current bar number

	TicksPerBeat int // How many ticks in one beat
	BeatsPerBar  int // How many beats in one bar

	Tempo int // Beats per minute (BPM)

	// Event queue: events scheduled for future execution
	EventQueue []*ScheduledEvent
}

// ScheduledEvent represents an event scheduled to occur at a specific time.
type ScheduledEvent struct {
	Tick     int         // When this event should occur (in ticks)
	Beat     int         // Beat number (for convenience)
	Bar      int         // Bar number (for convenience)
	EventType EventType  // Type of event
	Payload  interface{} // Event-specific data
}

// EventType represents the type of scheduled event.
type EventType int

const (
	EventTypeAction EventType = iota // Execute an action
	EventTypeMotif                   // Instantiate a motif
	EventTypeSync                    // Synchronization point
	EventTypeMarker                  // Marker for debugging/logging
)

// NewTemporalGrid creates a new temporal grid with default settings.
func NewTemporalGrid() *TemporalGrid {
	return &TemporalGrid{
		CurrentTick:  0,
		CurrentBeat:  0,
		CurrentBar:   0,
		TicksPerBeat: 480,  // Standard MIDI resolution
		BeatsPerBar:  4,    // 4/4 time signature
		Tempo:        120,  // 120 BPM
		EventQueue:   []*ScheduledEvent{},
	}
}

// Tick advances the temporal grid by one tick.
// Processes all events scheduled for the current tick.
func (tg *TemporalGrid) Tick() []*ScheduledEvent {
	tg.CurrentTick++

	// Update beat and bar counters
	if tg.CurrentTick%tg.TicksPerBeat == 0 {
		tg.CurrentBeat++
		if tg.CurrentBeat%tg.BeatsPerBar == 0 {
			tg.CurrentBar++
			tg.CurrentBeat = 0
		}
	}

	// Process events scheduled for this tick
	var readyEvents []*ScheduledEvent
	var remainingEvents []*ScheduledEvent

	for _, event := range tg.EventQueue {
		if event.Tick <= tg.CurrentTick {
			readyEvents = append(readyEvents, event)
		} else {
			remainingEvents = append(remainingEvents, event)
		}
	}

	tg.EventQueue = remainingEvents
	return readyEvents
}

// AdvanceBeat advances the temporal grid by one beat.
// This is a convenience method that ticks multiple times.
func (tg *TemporalGrid) AdvanceBeat() []*ScheduledEvent {
	targetTick := tg.CurrentTick + tg.TicksPerBeat
	var allEvents []*ScheduledEvent

	for tg.CurrentTick < targetTick {
		events := tg.Tick()
		allEvents = append(allEvents, events...)
	}

	return allEvents
}

// AdvanceBar advances the temporal grid by one bar.
func (tg *TemporalGrid) AdvanceBar() []*ScheduledEvent {
	targetTick := tg.CurrentTick + (tg.TicksPerBeat * tg.BeatsPerBar)
	var allEvents []*ScheduledEvent

	for tg.CurrentTick < targetTick {
		events := tg.Tick()
		allEvents = append(allEvents, events...)
	}

	return allEvents
}

// ScheduleEvent adds an event to the event queue.
func (tg *TemporalGrid) ScheduleEvent(event *ScheduledEvent) {
	// Insert event in sorted order (by tick)
	inserted := false
	for i, existing := range tg.EventQueue {
		if event.Tick < existing.Tick {
			// Insert before this event
			tg.EventQueue = append(tg.EventQueue[:i], append([]*ScheduledEvent{event}, tg.EventQueue[i:]...)...)
			inserted = true
			break
		}
	}
	if !inserted {
		tg.EventQueue = append(tg.EventQueue, event)
	}
}

// ScheduleAtBeat schedules an event at a specific beat number.
func (tg *TemporalGrid) ScheduleAtBeat(beat int, eventType EventType, payload interface{}) {
	tick := beat * tg.TicksPerBeat
	bar := beat / tg.BeatsPerBar
	event := &ScheduledEvent{
		Tick:      tick,
		Beat:      beat,
		Bar:       bar,
		EventType: eventType,
		Payload:   payload,
	}
	tg.ScheduleEvent(event)
}

// ScheduleAtBar schedules an event at the start of a specific bar.
func (tg *TemporalGrid) ScheduleAtBar(bar int, eventType EventType, payload interface{}) {
	beat := bar * tg.BeatsPerBar
	tg.ScheduleAtBeat(beat, eventType, payload)
}

// GetCurrentTime returns the current time in various units.
func (tg *TemporalGrid) GetCurrentTime() (tick, beat, bar int) {
	return tg.CurrentTick, tg.CurrentBeat, tg.CurrentBar
}

// SetTempo changes the tempo (BPM).
func (tg *TemporalGrid) SetTempo(bpm int) {
	tg.Tempo = bpm
}

// SetTimeSignature changes the time signature.
func (tg *TemporalGrid) SetTimeSignature(beatsPerBar, ticksPerBeat int) {
	tg.BeatsPerBar = beatsPerBar
	tg.TicksPerBeat = ticksPerBeat
}

// Reset resets the temporal grid to the initial state.
func (tg *TemporalGrid) Reset() {
	tg.CurrentTick = 0
	tg.CurrentBeat = 0
	tg.CurrentBar = 0
	tg.EventQueue = []*ScheduledEvent{}
}

// ==================== Motif Registry (M2.3) ====================

// MotifRegistry stores and manages motif definitions.
// Motifs are reusable temporal patterns that can be instantiated at different times.
type MotifRegistry struct {
	motifs map[string]*MotifDefinition
}

// MotifDefinition represents a defined motif.
type MotifDefinition struct {
	Name       string
	Parameters []string      // Parameter names
	Duration   int           // Duration in beats
	Body       interface{}   // AST node representing the motif body
	// TODO(M2.3): Add more motif metadata (key signature, etc.)
}

// NewMotifRegistry creates a new motif registry.
func NewMotifRegistry() *MotifRegistry {
	return &MotifRegistry{
		motifs: make(map[string]*MotifDefinition),
	}
}

// RegisterMotif adds a motif definition to the registry.
func (mr *MotifRegistry) RegisterMotif(motif *MotifDefinition) error {
	if _, exists := mr.motifs[motif.Name]; exists {
		return fmt.Errorf("motif '%s' already registered", motif.Name)
	}
	mr.motifs[motif.Name] = motif
	return nil
}

// GetMotif retrieves a motif definition by name.
func (mr *MotifRegistry) GetMotif(name string) (*MotifDefinition, bool) {
	motif, exists := mr.motifs[name]
	return motif, exists
}

// InstantiateMotif creates a scheduled instance of a motif.
// This is called when a motif is used in a score.
// Returns a list of events to be scheduled on the temporal grid.
func (mr *MotifRegistry) InstantiateMotif(name string, startBeat int, args map[string]Value) ([]*ScheduledEvent, error) {
	motif, exists := mr.motifs[name]
	if !exists {
		return nil, fmt.Errorf("motif '%s' not found", name)
	}

	// TODO(M2.3): Implement motif instantiation logic
	// This will involve:
	// 1. Binding parameters to arguments
	// 2. Evaluating the motif body
	// 3. Generating events based on the motif structure
	// 4. Applying transformations (transposition, etc.)

	// Placeholder: return empty event list
	return []*ScheduledEvent{}, nil
}

// ListMotifs returns all registered motif names.
func (mr *MotifRegistry) ListMotifs() []string {
	names := make([]string, 0, len(mr.motifs))
	for name := range mr.motifs {
		names = append(names, name)
	}
	return names
}

// ==================== Bar Range ====================

// BarRange represents a range of bars (e.g., bars 1..4).
type BarRange struct {
	Start int // Start bar (inclusive)
	End   int // End bar (inclusive)
}

// NewBarRange creates a new bar range.
func NewBarRange(start, end int) *BarRange {
	return &BarRange{
		Start: start,
		End:   end,
	}
}

// Contains checks if a bar number is within this range.
func (br *BarRange) Contains(bar int) bool {
	return bar >= br.Start && bar <= br.End
}

// Duration returns the duration of this range in bars.
func (br *BarRange) Duration() int {
	return br.End - br.Start + 1
}

// ToBeats converts the bar range to a beat range.
func (br *BarRange) ToBeats(beatsPerBar int) (startBeat, endBeat int) {
	startBeat = br.Start * beatsPerBar
	endBeat = (br.End + 1) * beatsPerBar - 1
	return
}

// ==================== Future Extensions ====================

// TODO(M2.3): Implement full Musical DSL support
// - Track management (multiple concurrent tracks)
// - Motif variations (transposition, augmentation, diminution)
// - Harmony analysis (chord progressions, voice leading)
// - Synchronization between tracks

// TODO(M2.4): Integrate with Resonance Graph
// - Temporal events can trigger resonance updates
// - Resonance can modulate temporal parameters (tempo, duration)

// TODO(M2.5): Integrate with Entanglement
// - Entangled variables can be synchronized at specific time points
// - Temporal constraints on entanglement lifetime

import "fmt"


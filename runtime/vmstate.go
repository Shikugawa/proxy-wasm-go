package runtime

var currentState = &state{
	rootContexts:   make(map[uint32]RootContext),
	httpContexts:   make(map[uint32]HttpContext),
	streamContexts: make(map[uint32]StreamContext),
	callOuts:       make(map[uint32]uint32),
}

type state struct {
	newRootContext   func(contextID uint32) RootContext
	newStreamContext func(contextID uint32) StreamContext
	newHttpContext   func(contextID uint32) HttpContext
	rootContexts     map[uint32]RootContext
	httpContexts     map[uint32]HttpContext
	streamContexts   map[uint32]StreamContext
	activeContextID  uint32
	callOuts         map[uint32]uint32
}

func SetNewRootContext(f func(contextID uint32) RootContext) {
	currentState.newRootContext = f
}

func SetNewHttpContext(f func(contextID uint32) HttpContext) {
	currentState.newHttpContext = f
}

func SetNewStreamContext(f func(contextID uint32) StreamContext) {
	currentState.newStreamContext = f
}

func (s *state) createRootContext(contextID uint32) {
	var ctx RootContext
	if s.newRootContext == nil {
		ctx = &DefaultContext{}
	} else {
		ctx = s.newRootContext(contextID)
	}

	s.rootContexts[contextID] = ctx
}

func (s *state) createStreamContext(contextID uint32, rootContextID uint32) {
	if _, ok := currentState.rootContexts[rootContextID]; !ok {
		panic("invalid root context id")
	}

	if _, ok := currentState.streamContexts[contextID]; ok {
		panic("context id duplicated")
	}

	currentState.streamContexts[contextID] = currentState.newStreamContext(contextID)
}

func (s *state) createHttpContext(contextID uint32, rootContextID uint32) {
	if _, ok := currentState.rootContexts[rootContextID]; !ok {
		panic("invalid root context id")
	}

	if _, ok := currentState.httpContexts[contextID]; ok {
		panic("context id duplicated")
	}

	currentState.httpContexts[contextID] = currentState.newHttpContext(contextID)
}

func (s *state) registerCallout(calloutID uint32) {
	if _, ok := s.callOuts[calloutID]; ok {
		panic("duplicated calloutID")
	}

	s.callOuts[calloutID] = s.activeContextID
}

func (s *state) setActiveContextID(contextID uint32) {
	s.activeContextID = contextID
}

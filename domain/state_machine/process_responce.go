package state_machine

func NewProcessResponse(addInitStates bool, states ...State) ProcessResponse {
	return &processResponse{states: states, addInitStates: addInitStates}
}

type processResponse struct {
	addInitStates bool
	states        []State
}

func (r *processResponse) States() []State {
	return r.states
}

func (r *processResponse) NeedAddInitStates() bool {
	return r.addInitStates
}

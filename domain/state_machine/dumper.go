package state_machine

import "runtime"

type Dumper struct {
}

func (s *Dumper) Dump() string {
	stack := []byte{}
	runtime.Stack(stack, true)
	return string(stack)
}

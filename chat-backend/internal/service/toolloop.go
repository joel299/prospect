package service

import (
	"context"
	"fmt"
	"github.com/iainfinito/chat-backend/internal/domain"
	"github.com/iainfinito/chat-backend/internal/providers/composio"
)

type Generator interface {
	Generate(context.Context, domain.AgentRequest) (domain.AgentResponse, error)
}
type ToolLoop struct {
	Generator Generator
	Executor  composio.ToolExecutor
	MaxRounds int
}

func (l ToolLoop) Run(ctx context.Context, r domain.AgentRequest) (domain.AgentResponse, error) {
	if l.MaxRounds <= 0 {
		l.MaxRounds = 4
	}
	for i := 0; i < l.MaxRounds; i++ {
		out, e := l.Generator.Generate(ctx, r)
		if e != nil {
			return out, e
		}
		if len(out.ToolCalls) == 0 {
			return out, nil
		}
		for _, call := range out.ToolCalls {
			if l.Executor == nil {
				return out, fmt.Errorf("tool executor unavailable")
			}
			result, e := l.Executor.Execute(ctx, call.Name, call.Arguments)
			if e != nil {
				return out, e
			}
			r.User = fmt.Sprintf("%s\nTOOL_RESULT %s: %v", r.User, call.Name, result)
		}
	}
	return domain.AgentResponse{}, fmt.Errorf("tool loop exceeded %d rounds", l.MaxRounds)
}

package internals

// Internal Processor interface
type ZProcessor[T any] interface {
	ZProcess(valPtr T, ctx Ctx)
}

type TransformProcessor[T any] struct {
	Transform Transform[T]
}

func (p *TransformProcessor[T]) ZProcess(valPtr T, ctx Ctx) {
	err := p.Transform(valPtr, ctx)
	if err != nil {
		s := ctx.(*SchemaCtx)
		s.AddIssue(s.IssueFromUnknownError(err))
		s.Exit = true
	}
}

func (p *TransformProcessor[T]) GetName() string {
	// Temporary. TODO actually implement this
	return ""
}

type TransformerInterface interface {
	GetName() string
}

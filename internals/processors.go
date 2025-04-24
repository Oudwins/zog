package internals

// Internal Processor interface
type ZProcessor interface {
	ZProcess(valPtr any, ctx Ctx)
}

type TransformProcessor struct {
	Transform Transform
}

func (p *TransformProcessor) ZProcess(valPtr any, ctx Ctx) {
	err := p.Transform(valPtr, ctx)
	if err != nil {
		s := ctx.(*SchemaCtx)
		s.AddIssue(s.IssueFromUnknownError(err))
		s.Exit = true
	}
}

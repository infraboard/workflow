package pipeline_test

import (
	"testing"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

func TestPipelineNextStep(t *testing.T) {
	sample := SamplePipeline()
	t.Log(sample)
	t.Log(sample.NextStep())
	sample.Stages[0].Steps[0].Success(map[string]string{"status": "ok"})
	t.Log(sample.NextStep())
	sample.Stages[0].Steps[1].Success(map[string]string{"status": "ok"})
	t.Log(sample.Stages[0].IsComplete())
	t.Log(sample.NextStep())

	sample.Stages[1].Steps[0].Success(map[string]string{"status": "ok"})
	t.Log(sample.NextStep())
	sample.Stages[1].Steps[1].Success(map[string]string{"status": "ok"})
	t.Log(sample.Stages[1].IsComplete())
	t.Log(sample.NextStep())

	t.Log(sample.IsComplete())
}

func TestStageNextStep(t *testing.T) {
	sample := SampleStage("stage01")
	t.Log(sample)
	t.Log(sample.NextStep())
	sample.Steps[0].Success(map[string]string{"status": "ok"})
	t.Log(sample.NextStep())
	sample.Steps[1].Success(map[string]string{"status": "ok"})
	t.Log(sample)
	t.Log(sample.NextStep())
	t.Log(sample.IsComplete())
}

func SamplePipeline() *pipeline.Pipeline {
	p := pipeline.NewDefaultPipeline()
	p.AddStage(SampleStage("stage01"))
	p.AddStage(SampleStage("stage02"))
	return p
}

func SampleStage(name string) *pipeline.Stage {
	stage := pipeline.NewDefaultStage()
	stage.Name = name

	s1 := pipeline.NewDefaultStep()
	s1.Action = name + ".action01"
	s2 := pipeline.NewDefaultStep()
	s2.Action = name + ".action02"
	stage.AddStep(s1)
	stage.AddStep(s2)
	return stage
}

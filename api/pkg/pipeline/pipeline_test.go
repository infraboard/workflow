package pipeline_test

import (
	"encoding/json"
	"testing"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/stretchr/testify/assert"
)

var testP = `{
	"id": "c3r8psea0brhlbg0qdd0",
	"resource_version": 831313,
	"domain": "admin-domain",
	"namespace": "c16mhsddrei91m4ri0jg",
	"create_at": 1626770673799,
	"create_by": "admin",
	"name": "test33",
	"with": null,
	"mount": null,
	"tags": null,
	"description": "",
	"on": null,
	"status": {
		"current_flow": 0,
		"start_at": 1626770673846,
		"end_at": 0,
		"status": "EXECUTING",
		"scheduler_node": "MacBook-Pro-9",
		"message": ""
	},
	"stages": [
		{
			"id": 1,
			"name": "stage1",
			"needs": null,
			"steps": [
				{
					"key": "c16mhsddrei91m4ri0jg.c3r8psea0brhlbg0qdd0.1.1",
					"id": 1,
					"name": "step1.1",
					"action": "action01",
					"with": {
						"ENV1": "env1",
						"ENV2": "env2"
					},
					"is_parallel": false,
					"ignore_failed": false,
					"with_audit": false,
					"audit_params": null,
					"with_notify": false,
					"notify_params": null,
					"webhook": null,
					"node_selector": null,
					"status": {
						"flow_number": 0,
						"start_at": 0,
						"end_at": 0,
						"status": "PENDDING",
						"scheduled_node": "",
						"audit_response": "",
						"notify_at": 0,
						"notify_error": "",
						"message": "",
						"response": null
					}
				}
			]
		},
		{
			"id": 2,
			"name": "stage2",
			"needs": null,
			"steps": [
				{
					"id": 1,
					"name": "step2.1",
					"action": "action01",
					"with": {
						"ENV1": "env3",
						"ENV2": "env4"
					},
					"is_parallel": false,
					"ignore_failed": false,
					"with_audit": false,
					"audit_params": null,
					"with_notify": false,
					"notify_params": null,
					"webhook": null,
					"node_selector": null,
					"status": {
						"flow_number": 0,
						"start_at": 0,
						"end_at": 0,
						"status": "PENDDING",
						"scheduled_node": "",
						"audit_response": "",
						"notify_at": 0,
						"notify_error": "",
						"message": "",
						"response": null
					}
				}
			]
		}
	]
}`

func TestPipelineFlow(t *testing.T) {
	should := assert.New(t)
	sample := pipeline.NewDefaultPipeline()
	err := json.Unmarshal([]byte(testP), sample)
	should.NoError(err)

	t.Log(sample)
	steps, ok := sample.NextStep()
	t.Log(sample.GetCurrentFlow())
	t.Log("is complete: ", ok, "start steps: ", steps)
}

func TestPipelineNextStepOK(t *testing.T) {
	sample := SamplePipeline()
	steps, ok := sample.NextStep()
	t.Log("is complete: ", ok, "start steps: ", steps)
	sample.Stages[0].Steps[0].Success(map[string]string{"status": "ok"})
	t.Log(sample.NextStep())
	sample.Stages[0].Steps[1].Success(map[string]string{"status": "ok"})
	t.Log(sample.Stages[0].IsPassed())
	t.Log(sample.NextStep())

	sample.Stages[1].Steps[0].Success(map[string]string{"status": "ok"})
	t.Log(sample.NextStep())
	sample.Stages[1].Steps[1].Success(map[string]string{"status": "ok"})
	t.Log(sample.Stages[1].IsPassed())

	steps, ok = sample.NextStep()
	t.Log("is complete: ", ok, "steps: ", steps)
}

func TestPipelineNextStepBreak(t *testing.T) {
	sample := SamplePipeline()
	t.Log(sample)
	t.Log(sample.NextStep())
	sample.Stages[0].Steps[0].Success(map[string]string{"status": "ok"})
	t.Log(sample.NextStep())
	sample.Stages[0].Steps[1].Failed("step failed")
	t.Log("step is passed: ", sample.Stages[0].IsPassed())
	steps, ok := sample.NextStep()
	t.Log("is complete: ", ok, "steps: ", steps)
}

func TestPipelineCurrentFlowOK(t *testing.T) {
	sample := SamplePipeline()
	t.Log(sample)
	t.Log(sample.NextStep())
	sample.Stages[0].Steps[0].Success(map[string]string{"status": "ok"})
	t.Log("current flow: ", sample.GetCurrentFlow())

	steps, ok := sample.NextStep()
	t.Log("is complete: ", ok, "steps: ", steps)

	sample.Stages[0].Steps[1].Run()
	t.Log("current flow: ", sample.GetCurrentFlow())
	steps, ok = sample.NextStep()
	t.Log("is complete: ", ok, "steps: ", steps)

	sample.Stages[0].Steps[1].Success(map[string]string{"status": "ok"})
	steps, ok = sample.NextStep()
	t.Log("is complete: ", ok, "steps: ", steps)
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
	t.Log(sample.IsPassed())
}

func TestUpdateStep(t *testing.T) {
	should := assert.New(t)
	sample := SamplePipeline()
	s1 := pipeline.NewDefaultStep()
	s1.Key = "..1.1"
	s1.Action = "update01"
	err := sample.UpdateStep(s1)

	if should.NoError(err) {
		t.Log(sample)
	}
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
	s1.Key = "..1.1"
	s2 := pipeline.NewDefaultStep()
	s2.Action = name + ".action02"
	stage.AddStep(s1)
	stage.AddStep(s2)
	return stage
}

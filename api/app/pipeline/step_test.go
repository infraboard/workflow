package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/workflow/api/app/pipeline"
)

func TestStepNamespace(t *testing.T) {
	should := assert.New(t)

	s := pipeline.NewDefaultStep()
	s.Key = "c16mhsddrei91m4ri0jg.c3iqcama0brimaq08e40.2.1"
	should.Equal(s.GetNamespace(), "c16mhsddrei91m4ri0jg")
}

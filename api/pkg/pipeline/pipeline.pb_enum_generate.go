// Code generated by github.com/infraboard/mcube
// DO NOT EDIT

package pipeline

import (
	"bytes"
	"fmt"
	"strings"
)

// ParseSTEP_STATUSFromString Parse STEP_STATUS from string
func ParseSTEP_STATUSFromString(str string) (STEP_STATUS, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := STEP_STATUS_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown STEP_STATUS: %s", str)
	}

	return STEP_STATUS(v), nil
}

// Equal type compare
func (t STEP_STATUS) Equal(target STEP_STATUS) bool {
	return t == target
}

// IsIn todo
func (t STEP_STATUS) IsIn(targets ...STEP_STATUS) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t STEP_STATUS) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *STEP_STATUS) UnmarshalJSON(b []byte) error {
	ins, err := ParseSTEP_STATUSFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}

// ParsePIPELINE_STATUSFromString Parse PIPELINE_STATUS from string
func ParsePIPELINE_STATUSFromString(str string) (PIPELINE_STATUS, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := PIPELINE_STATUS_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown PIPELINE_STATUS: %s", str)
	}

	return PIPELINE_STATUS(v), nil
}

// Equal type compare
func (t PIPELINE_STATUS) Equal(target PIPELINE_STATUS) bool {
	return t == target
}

// IsIn todo
func (t PIPELINE_STATUS) IsIn(targets ...PIPELINE_STATUS) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t PIPELINE_STATUS) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *PIPELINE_STATUS) UnmarshalJSON(b []byte) error {
	ins, err := ParsePIPELINE_STATUSFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}

// ParseVALUE_TYPEFromString Parse VALUE_TYPE from string
func ParseVALUE_TYPEFromString(str string) (VALUE_TYPE, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := VALUE_TYPE_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown VALUE_TYPE: %s", str)
	}

	return VALUE_TYPE(v), nil
}

// Equal type compare
func (t VALUE_TYPE) Equal(target VALUE_TYPE) bool {
	return t == target
}

// IsIn todo
func (t VALUE_TYPE) IsIn(targets ...VALUE_TYPE) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t VALUE_TYPE) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *VALUE_TYPE) UnmarshalJSON(b []byte) error {
	ins, err := ParseVALUE_TYPEFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}

// ParseAUDIT_RESPONSEFromString Parse AUDIT_RESPONSE from string
func ParseAUDIT_RESPONSEFromString(str string) (AUDIT_RESPONSE, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := AUDIT_RESPONSE_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown AUDIT_RESPONSE: %s", str)
	}

	return AUDIT_RESPONSE(v), nil
}

// Equal type compare
func (t AUDIT_RESPONSE) Equal(target AUDIT_RESPONSE) bool {
	return t == target
}

// IsIn todo
func (t AUDIT_RESPONSE) IsIn(targets ...AUDIT_RESPONSE) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t AUDIT_RESPONSE) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *AUDIT_RESPONSE) UnmarshalJSON(b []byte) error {
	ins, err := ParseAUDIT_RESPONSEFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}

// ParseSTEP_CREATE_BYFromString Parse STEP_CREATE_BY from string
func ParseSTEP_CREATE_BYFromString(str string) (STEP_CREATE_BY, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := STEP_CREATE_BY_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown STEP_CREATE_BY: %s", str)
	}

	return STEP_CREATE_BY(v), nil
}

// Equal type compare
func (t STEP_CREATE_BY) Equal(target STEP_CREATE_BY) bool {
	return t == target
}

// IsIn todo
func (t STEP_CREATE_BY) IsIn(targets ...STEP_CREATE_BY) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t STEP_CREATE_BY) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *STEP_CREATE_BY) UnmarshalJSON(b []byte) error {
	ins, err := ParseSTEP_CREATE_BYFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}

// ParsePIPELINE_WATCH_MODFromString Parse PIPELINE_WATCH_MOD from string
func ParsePIPELINE_WATCH_MODFromString(str string) (PIPELINE_WATCH_MOD, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := PIPELINE_WATCH_MOD_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown PIPELINE_WATCH_MOD: %s", str)
	}

	return PIPELINE_WATCH_MOD(v), nil
}

// Equal type compare
func (t PIPELINE_WATCH_MOD) Equal(target PIPELINE_WATCH_MOD) bool {
	return t == target
}

// IsIn todo
func (t PIPELINE_WATCH_MOD) IsIn(targets ...PIPELINE_WATCH_MOD) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t PIPELINE_WATCH_MOD) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *PIPELINE_WATCH_MOD) UnmarshalJSON(b []byte) error {
	ins, err := ParsePIPELINE_WATCH_MODFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}

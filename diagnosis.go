package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// DiagnosisReq is a struct to request diagnosis
type DiagnosisReq struct {
	Sex       Sex               `json:"sex"`
	Age       Age               `json:"age"`
	Evidences []Evidence        `json:"evidence"`
	Extras    DiagnosisReqExtras `json:"extras"`
}

// DiagnosisReqExtras contains extra params for DiagnosisReq
type DiagnosisReqExtras struct {
	DisableGroups bool `json:"disable_groups"`
}

// DiagnosisRes is a response struct for diagnosis
type DiagnosisRes struct {
	Question   Question                `json:"question"`
	Conditions []DiagnosisConditionRes `json:"conditions"`
	ShouldStop bool                    `json:"should_stop"`
	Extras     interface{}             `json:"extras"`
}

// QuestionType is a list of question types
type QuestionType string

const (
	// QuestionTypeSingle single question
	QuestionTypeSingle QuestionType = "single"
	// QuestionTypeGroupSingle question group
	QuestionTypeGroupSingle QuestionType = "group_single"
	// QuestionTypeGroupMultiple multiple question groups
	QuestionTypeGroupMultiple QuestionType = "group_multiple"
)

func (qt QuestionType) Ptr() *QuestionType { return &qt }
func (qt QuestionType) String() string     { return string(qt) }

func (qt *QuestionType) IsValid() error {
	_, err := QuestionTypeFromString(qt.String())
	if err != nil {
		return err
	}
	return nil
}

func QuestionTypeFromString(x string) (QuestionType, error) {
	switch strings.ToLower(x) {
	case "single":
		return QuestionTypeSingle, nil
	case "group_single":
		return QuestionTypeGroupSingle, nil
	case "group_multiple":
		return QuestionTypeGroupMultiple, nil
	default:
		return "", fmt.Errorf("unexpected value for Question Type: %q", x)
	}
}

// Question struct
type Question struct {
	Type  QuestionType   `json:"type"`
	Text  string         `json:"text"`
	Items []QuestionItem `json:"items"`
}

// QuestionItem question item struct
type QuestionItem struct {
	ID      string               `json:"id"`
	Name    string               `json:"name"`
	Choices []QuestionItemChoice `json:"choices"`
}

// QuestionItemChoice question item choice struct
type QuestionItemChoice struct {
	ID    EvidenceChoiceID `json:"id"`
	Label string           `json:"label"`
}

// DiagnosisConditionRes is a response struct for condition + probability
type DiagnosisConditionRes struct {
	Condition struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		CommonName string `json:"common_name"`
		ICD10Code  string `json:"icd10_code"`
	} `json:"condition"`
	Probability float64 `json:"probability"`
}

// Diagnosis is a func to request diagnosis for given data
func (a *App) Diagnosis(dr DiagnosisReq) (*DiagnosisRes, error) {
	if dr.Sex.IsValid() != nil {
		return nil, errors.New("Unexpected value for Sex")
	}
	req, err := a.prepareRequest("POST", "diagnosis", dr)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	err = checkResponse(res)

	if err != nil{
		return nil, err
	}

	err = checkResponse(res)

	if err != nil{
		return nil, err
	}

	r := DiagnosisRes{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

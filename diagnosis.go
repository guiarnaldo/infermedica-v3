package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type DiagnosisReq struct {
	Sex         Sex                 `json:"sex"`
	Age         Age                 `json:"age"`
	EvaluatedAt string              `json:"evaluated_at,omitempty"`
	Evidences   []Evidence          `json:"evidence,omitempty"`
	Extras      *DiagnosisReqExtras `json:"extras,omitempty"`
}

type DiagnosisReqExtras struct {
	EnableTriage3              bool          `json:"enable_triage_3,omitempty"`
	InterviewMode              InterviewMode `json:"interview_mode,omitempty"`
	DisableAdaptiveRanking     bool          `json:"disable_adaptive_ranking,omitempty"`
	EnableExplanations         bool          `json:"enable_explanations,omitempty"`
	EnableThirdPersonQuestions bool          `json:"enable_third_person_questions,omitempty"`
	IncludeConditionDetails    bool          `json:"include_condition_details,omitempty"` // When included in a request, each condition in the output gains an additional section - ConditionDetails
	DisableIntimateContent     bool          `json:"disable_intimate_content,omitempty"`
	EnableSymptomDuration      bool          `json:"enable_symptom_duration,omitempty"` // This flag enables questions of the type duration which contain a new field EvidenceID
}

type DiagnosisRes struct {
	Question         Question         `json:"question"`
	Conditions       []Conditions     `json:"conditions"`
	ShouldStop       bool             `json:"should_stop"`
	ConditionDetails ConditionDetails `json:"condition_details"` // Only enabled when IncludeConditionDetails is true
	Extras           struct {
	} `json:"extras"`
	HasEmergencyEvidence bool `json:"has_emergency_evidence"`
}

type Question struct {
	Type       QuestionType   `json:"type"`
	Text       string         `json:"text"`
	EvidenceID string         `json:"evidence_id"` // Only when questions of the type duration is enabled
	Items      []QuestionItem `json:"items"`
	Extras     struct {
	} `json:"extras"`
}

type QuestionItem struct {
	ID      string               `json:"id"`
	Name    string               `json:"name"`
	Choices []QuestionItemChoice `json:"choices"`
}

type QuestionItemChoice struct {
	ID    EvidenceChoiceID `json:"id"`
	Label string           `json:"label"`
}

type Conditions struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	CommonName  string  `json:"common_name"`
	Probability float64 `json:"probability"`
}

type ConditionDetails struct {
	Icd10Code           string                   `json:"icd10_code"`
	Category            ConditionDetailsCategory `json:"category"`
	Prevalence          Prevalence               `json:"prevalence"`
	Severity            Severity                 `json:"severity"`
	Acuteness           Acuteness                `json:"acuteness"`
	TriageLevel         TriageLevel              `json:"triage_level"`
	Hint                string                   `json:"hint"`
	HasPatientEducation bool                     `json:"has_patient_education"`
}

type ConditionDetailsCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// QuestionType is a list of question types
type QuestionType string

const (
	QuestionTypeSingle        QuestionType = "single"         // QuestionTypeSingle single question
	QuestionTypeGroupSingle   QuestionType = "group_single"   // QuestionTypeGroupSingle question group
	QuestionTypeGroupMultiple QuestionType = "group_multiple" // QuestionTypeGroupMultiple multiple question groups
	QuestionTypeDuration      QuestionType = "duration"       // QuestionTypeDuration only available when EnableSymptomDuration is true
)

type InterviewMode string

const (
	InterviewModeDefault InterviewMode = "default" // suitable for symptom checking applications, providing the right balance between duration of interview and accuracy of the presented results
	InterviewModeTriage  InterviewMode = "triage"  // suitable for triage applications where duration of the interview is shorter and optimized for the assessment of the correct triage level rather than accuracy of the final list of most probable conditions
)

func (qt *QuestionType) IsValid() error {
	_, err := QuestionTypeFromString(string(*qt))
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
	case "duration":
		return QuestionTypeGroupMultiple, nil
	default:
		return "", fmt.Errorf("infermedica: unexpected value for Question Type: %q", x)
	}
}

// Diagnosis is a func to request diagnosis for given data
func (a *App) Diagnosis(dr DiagnosisReq) (*DiagnosisRes, error) {
	if dr.Sex.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: Unexpected value for Sex")
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
	if err != nil {
		return nil, err
	}

	var r DiagnosisRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

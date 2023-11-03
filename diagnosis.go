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
	Sex         Sex                `json:"sex"`
	Age         Age                `json:"age"`
	EvaluatedAt string             `json:"evaluated_at,omitempty"`
	Evidences   []Evidence         `json:"evidence"`
	Extras      DiagnosisReqExtras `json:"extras"`
}

// DiagnosisReqExtras contains extra params for DiagnosisReq
type DiagnosisReqExtras struct {
	DisableGroups              bool          `json:"disable_groups"`                // Using this option forces diagnosis to return only questions of the single type, disabling those of the group_single and group_multiple types
	EnableTriage3              bool          `json:"enable_triage_3"`               // Using this option disables the 5-level triage mode that is recommended for all applications
	InterviewMode              InterviewMode `json:"interview_mode"`                // This option allows you to control the behavior of the question selection algorithm. The interview mode may have an influence on the duration of the interview as well as the sequencing of questions
	DisableAdaptiveRanking     bool          `json:"disable_adaptive_ranking"`      // When adaptive ranking is enabled, only conditions having sufficient probability will be returned. Additionally, ranking will be limited to 8 conditions. We strongly recommend not disabling this option.
	EnableExplanations         bool          `json:"enable_explanations"`           // Explanation is optional and not every question/question item will have it
	EnableThirdPersonQuestions bool          `json:"enable_third_person_questions"` // When this parameter is set to true, each question from diagnosis is returned in third person form
	IncludeConditionDetails    bool          `json:"include_condition_details"`     // When included in a request, each condition in the output gains an additional section - ConditionDetails
	DisableIntimateContent     bool          `json:"disable_intimate_content"`      // Gives the possibility of excluding intimate concepts from the response e.g concepts related to sexual activity.
	EnableSymptomDuration      bool          `json:"enable_symptom_duration"`       // This flag enables questions of the type duration which contain a new field evidence_id
}

// DiagnosisRes is a response struct for diagnosis
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
	QuestionTypeDuration      QuestionType = "duration"       // QuestionTypeDuration only avaliable when EnableSymptomDuration is true
)

type InterviewMode string

const (
	InterviewModeDefault InterviewMode = "default" // suitable for symptom checking applications, providing the right balance between duration of interview and accuracy of the presented results
	InterviewModeTriage  InterviewMode = "triage"  // suitable for triage applications where duration of the interview is shorter and optimized for the assessment of the correct triage level rather than accuracy of the final list of most probable conditions
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
	case "duration":
		return QuestionTypeGroupMultiple, nil
	default:
		return "", fmt.Errorf("unexpected value for Question Type: %q", x)
	}
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
	if err != nil {
		return nil, err
	}

	r := DiagnosisRes{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

package infermedica

import (
	"fmt"
	"strings"
)

type App struct {
	baseURL     string
	appID       string
	appKey      string
	model       string
	interviewID string
}

func NewApp(id, key, model, interviewID string) App {
	return App{
		baseURL:     "https://api.infermedica.com/v3/",
		appID:       id,
		appKey:      key,
		model:       model,
		interviewID: interviewID,
	}
}

type Sex string

const (
	SexMale   Sex = "male"
	SexFemale Sex = "female"
)

func (s Sex) Ptr() *Sex      { return &s }
func (s Sex) String() string { return string(s) }

func (s *Sex) IsValid() error {
	_, err := SexFromString(s.String())
	if err != nil {
		return err
	}
	return nil
}

func SexFromString(x string) (Sex, error) {
	switch strings.ToLower(x) {
	case "male":
		return SexMale, nil
	case "female":
		return SexFemale, nil
	default:
		return "", fmt.Errorf("infermedica: unexpected value for Sex: %q", x)
	}
}

type SexFilter string

const (
	SexFilterBoth   SexFilter = "both"
	SexFilterMale   SexFilter = "male"
	SexFilterFemale SexFilter = "female"
)

func (s SexFilter) Ptr() *SexFilter { return &s }
func (s SexFilter) String() string  { return string(s) }

func (s *SexFilter) IsValid() error {
	_, err := SexFilterFromString(s.String())
	if err != nil {
		return err
	}
	return nil
}

func SexFilterFromString(x string) (SexFilter, error) {
	switch strings.ToLower(x) {
	case "both":
		return SexFilterBoth, nil
	case "male":
		return SexFilterMale, nil
	case "female":
		return SexFilterFemale, nil
	default:
		return "", fmt.Errorf("infermedica: unexpected value for SexFilter: %q", x)
	}
}

type EvidenceChoiceID string

const (
	EvidenceChoiceIDPresent EvidenceChoiceID = "present"
	EvidenceChoiceIDAbsent  EvidenceChoiceID = "absent"
	EvidenceChoiceIDUnknown EvidenceChoiceID = "unknown"
)

func (ecID EvidenceChoiceID) Ptr() *EvidenceChoiceID { return &ecID }
func (ecID EvidenceChoiceID) String() string         { return string(ecID) }

func (ecID EvidenceChoiceID) IsValid() error {
	_, err := EvidenceChoiceIDFromString(ecID.String())
	if err != nil {
		return err
	}
	return nil
}

func EvidenceChoiceIDFromString(x string) (EvidenceChoiceID, error) {
	switch strings.ToLower(x) {
	case "present":
		return EvidenceChoiceIDPresent, nil
	case "absent":
		return EvidenceChoiceIDAbsent, nil
	case "unknown":
		return EvidenceChoiceIDUnknown, nil
	default:
		return "", fmt.Errorf("infermedica: unexpected value for evidence choice id: %q", x)
	}
}

type DurationUnit string

const (
	DurationUnitWeek   DurationUnit = "week"
	DurationUnitDay    DurationUnit = "day"
	DurationUnitHour   DurationUnit = "hour"
	DurationUnitMinute DurationUnit = "minute"
)

// Contains source valid types
type EvidenceSource string

const (
	EvidenceSourceInitial    EvidenceSource = "initial"
	EvidenceSourceSuggest    EvidenceSource = "suggest"
	EvidenceSourcePredefined EvidenceSource = "predefined"
	EvidenceSourceRedFlags   EvidenceSource = "red_flags"
)

type Evidence struct {
	ID         string           `json:"id"`        // Required
	ChoiceID   EvidenceChoiceID `json:"choice_id"` // Required
	ObservedAt string           `json:"observed_at,omitempty"`
	Source     EvidenceSource   `json:"source,omitempty"`
	Duration   Duration         `json:"duration,omitempty"` // Required only when EnableSymptomDuration is true
}

// Required only when EnableSymptomDuration is true
type Duration struct {
	Value int          `json:"value,omitempty"`
	Unit  DurationUnit `json:"unit,omitempty"`
}

type AgeUnit string

const (
	AgeUnitYear  AgeUnit = "year"  // Age in years (Default)
	AgeUnitMonth AgeUnit = "month" // Age in months
)

type Age struct {
	Value int     `json:"value"` // Numeric value, this attribute is required
	Unit  AgeUnit `json:"unit"`  // This attribute is optional and the default value is year
}

// Base struct for diagnosis, triage and recommend specialist
type ObservationReq struct {
	Sex         Sex                  `json:"sex"`
	Age         Age                  `json:"age"`
	EvaluatedAt string               `json:"evaluated_at,omitempty"`
	Evidences   []Evidence           `json:"evidence"`
	Extras      ObservationReqExtras `json:"extras"`
}

// Contains extra params for ObservationReq
type ObservationReqExtras struct {
	DisableGroups              bool              `json:"disable_groups,omitempty"`      // Using this option forces diagnosis to return only questions of the single type, disabling those of the group_single and group_multiple types
	EnableTriage3              bool              `json:"enable_triage_3"`               // Using this option disables the 5-level triage mode that is recommended for all applications
	InterviewMode              InterviewMode     `json:"interview_mode"`                // This option allows you to control the behavior of the question selection algorithm. The interview mode may have an influence on the duration of the interview as well as the sequencing of questions
	DisableAdaptiveRanking     bool              `json:"disable_adaptive_ranking"`      // When adaptive ranking is enabled, only conditions having sufficient probability will be returned. Additionally, ranking will be limited to 8 conditions. We strongly recommend not disabling this option.
	EnableExplanations         bool              `json:"enable_explanations"`           // Explanation is optional and not every question/question item will have it
	EnableThirdPersonQuestions bool              `json:"enable_third_person_questions"` // When this parameter is set to true, each question from diagnosis is returned in third person form
	IncludeConditionDetails    bool              `json:"include_condition_details"`     // When included in a request, each condition in the output gains an additional section - ConditionDetails
	DisableIntimateContent     bool              `json:"disable_intimate_content"`      // Gives the possibility of excluding intimate concepts from the response e.g concepts related to sexual activity.
	EnableSymptomDuration      bool              `json:"enable_symptom_duration"`       // This flag enables questions of the type duration which contain a new field evidence_id
	SpecialistMapping          map[string]string `json:"specialist_mapping"`            // The recommend_specialist endpoint allows for the remapping of specified specialties in a many-to-one fashion. This is useful when some specialties are not appropriate for the regional, regulatory, or clinical setting (or unwanted for any other reason)
}

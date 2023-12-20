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
	devMode     bool
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

// EnableDevMode
func (a *App) EnableDevMode() {
	a.devMode = true
}

func (a *App) DisableDevMode() {
	a.devMode = false
}

type Sex string

const (
	SexMale   Sex = "male"
	SexFemale Sex = "female"
)

func (s *Sex) IsValid() error {
	_, err := SexFromString(string(*s))
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

func (s *SexFilter) IsValid() error {
	_, err := SexFilterFromString(string(*s))
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

func (ecID EvidenceChoiceID) IsValid() error {
	_, err := EvidenceChoiceIDFromString(string(ecID))
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
	ID         string           `json:"id,omitempty"`
	ChoiceID   EvidenceChoiceID `json:"choice_id,omitempty"`
	ObservedAt string           `json:"observed_at,omitempty"`
	Source     EvidenceSource   `json:"source,omitempty"`
	Duration   *Duration        `json:"duration,omitempty"` // Required only when EnableSymptomDuration is true
}

// Duration is required only when EnableSymptomDuration is true
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
	Value int     `json:"value"`          // Numeric value, this attribute is required
	Unit  AgeUnit `json:"unit,omitempty"` // This attribute is optional and the default value is year
}

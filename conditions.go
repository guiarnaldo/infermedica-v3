package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Prevalence string

const (
	PrevalenceVeryRare Prevalence = "very_rare"
	PrevalenceRare     Prevalence = "rare"
	PrevalenceModerate Prevalence = "moderate"
	PrevalenceCommon   Prevalence = "common"
)

func (p Prevalence) Ptr() *Prevalence { return &p }
func (p Prevalence) String() string   { return string(p) }

func (p *Prevalence) IsValid() error {
	_, err := PrevalenceFromString(p.String())
	if err != nil {
		return err
	}
	return nil
}

func PrevalenceFromString(x string) (Prevalence, error) {
	switch strings.ToLower(x) {
	case "very_rare":
		return PrevalenceVeryRare, nil
	case "rare":
		return PrevalenceRare, nil
	case "moderate":
		return PrevalenceModerate, nil
	case "common":
		return PrevalenceCommon, nil
	default:
		return "", fmt.Errorf("unexpected value for Prevalence: %q", x)
	}
}

type Acuteness string

const (
	AcutenessChronic                  Acuteness = "chronic"
	AcutenessChronicWithExacerbations Acuteness = "chronic_with_exacerbations"
	AcutenessAcutePotentiallyChronic  Acuteness = "acute_potentially_chronic"
	AcutenessAcute                    Acuteness = "acute"
)

func (a Acuteness) Ptr() *Acuteness { return &a }
func (a Acuteness) String() string  { return string(a) }

func (a *Acuteness) IsValid() error {
	_, err := AcutenessFromString(a.String())
	if err != nil {
		return err
	}
	return nil
}

func AcutenessFromString(x string) (Acuteness, error) {
	switch strings.ToLower(x) {
	case "chronic":
		return AcutenessChronic, nil
	case "chronic_with_exacerbations":
		return AcutenessChronicWithExacerbations, nil
	case "acute_potentially_chronic":
		return AcutenessAcutePotentiallyChronic, nil
	case "acute":
		return AcutenessAcute, nil
	default:
		return "", fmt.Errorf("unexpected value for Acuteness: %q", x)
	}
}

type Severity string

const (
	SeverityMild     Severity = "mild"
	SeverityModerate Severity = "moderate"
	SeveritySevere   Severity = "severe"
)

func (s Severity) Ptr() *Severity { return &s }
func (s Severity) String() string { return string(s) }

func (s *Severity) IsValid() error {
	_, err := SeverityFromString(s.String())
	if err != nil {
		return err
	}
	return nil
}

func SeverityFromString(x string) (Severity, error) {
	switch strings.ToLower(x) {
	case "mild":
		return SeverityMild, nil
	case "moderate":
		return SeverityModerate, nil
	case "severe":
		return SeveritySevere, nil
	default:
		return "", fmt.Errorf("unexpected value for Severity: %q", x)
	}
}

type Condition struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
	ICD10Code  string `json:"icd10_code"`
}

type ConditionRes struct {
	Condition
	SexFilter   SexFilter       `json:"sex_filter"`
	Categories  []string        `json:"categories"`
	Prevalence  Prevalence      `json:"prevalence"`
	Acuteness   Acuteness       `json:"acuteness"`
	Severity    Severity        `json:"severity"`
	Extras      ConditionExtras `json:"extras"`
	TriageLevel string          `json:"triage_level"`
}

type ConditionExtras struct {
	Hint      string `json:"hint"`
	ICD10Code string `json:"icd10_code"`
}

func (a *App) Conditions() (*[]ConditionRes, error) {
	req, err := a.prepareRequest("GET", "conditions", nil)
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
	r := []ConditionRes{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (a *App) ConditionsIDMap() (*map[string]ConditionRes, error) {
	r, err := a.Conditions()
	if err != nil {
		return nil, err
	}
	rmap := make(map[string]ConditionRes)
	for _, sr := range *r {
		rmap[sr.ID] = sr
	}
	return &rmap, nil
}

func (a *App) ConditionByID(id string) (*ConditionRes, error) {
	req, err := a.prepareGETRequest("conditions/" + id, nil)
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
	r := ConditionRes{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

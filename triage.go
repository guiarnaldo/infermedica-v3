package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type TriageReq struct {
	Sex         Sex              `json:"sex"`
	Age         Age              `json:"age"`
	Evidences   []Evidence       `json:"evidence,omitempty"`
	EvaluatedAt string           `json:"evaluated_at,omitempty"`
	Extras      *TriageReqExtras `json:"extras,omitempty"`
}

type TriageReqExtras struct {
	EnableTriage3         bool `json:"enable_triage_3,omitempty"`         // Using this option disables the 5-level triage mode that is recommended for all applications
	EnableSymptomDuration bool `json:"enable_symptom_duration,omitempty"` // This flag enables questions of the type duration which contain a new field EvidenceID
}

type TriageRes struct {
	TriageLevel                TriageLevel `json:"triage_level"`                // A classification of the case provided
	Serious                    []Serious   `json:"serious"`                     // A list of serious observations
	TeleconsultationApplicable bool        `json:"teleconsultation_applicable"` // The teleconsultation_applicable flag has been deprecated and will stop being supported in the near future
	RootCause                  string      `json:"root_cause"`                  // A root cause that explains the internal rationale of the underlying triage algorithm
}

type Serious struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	CommonName  string      `json:"common_name"`
	Seriousness Seriousness `json:"seriousness"`
	IsEmergency bool        `json:"is_emergency"` // IsEmergency is redundant with Seriousness field and should be treated as deprecated
}

type RootCause string

const (
	RootCauseEmergencyEvidencePresent        RootCause = "emergency_evidence_present"
	RootCauseEmergencySeriousEvidencePresent RootCause = "serious_evidence_present"
	RootCauseEmergencyConditionLikely        RootCause = "emergency_condition_likely"
	RootCauseEmergencyConditionPossible      RootCause = "emergency_condition_possible"
	RootCauseConsultationConditionLikely     RootCause = "consultation_condition_likely"
	RootCauseSelfCareSufficient              RootCause = "self_care_sufficient"
	RootCauseDiagnosisUnknown                RootCause = "diagnosis_unknown"
)

type Seriousness string

const (
	SeriousnessSerious                   Seriousness = "serious"
	SeriousnessSeriousEmergency          Seriousness = "emergency"
	SeriousnessSeriousEmergencyAmbulance Seriousness = "emergency_ambulance"
)

type TriageLevel string

const (
	TriageLevelEmergencyAmbulance TriageLevel = "emergency_ambulance"
	TriageLevelEmergency          TriageLevel = "emergency"
	TriageLevelConsultation24     TriageLevel = "consultation_24"
	TriageLevelConsultation       TriageLevel = "consultation"
	TriageLevelSelfCare           TriageLevel = "self_care"
)

func (s *TriageLevel) IsValid() error {
	_, err := TriageLevelFromString(string(*s))
	if err != nil {
		return err
	}
	return nil
}

func TriageLevelFromString(x string) (TriageLevel, error) {
	switch strings.ToLower(x) {
	case "emergency_ambulance":
		return TriageLevelEmergencyAmbulance, nil
	case "emergency":
		return TriageLevelEmergency, nil
	case "consultation_24":
		return TriageLevelConsultation24, nil
	case "consultation":
		return TriageLevelConsultation, nil
	case "self_care":
		return TriageLevelSelfCare, nil
	default:
		return "", fmt.Errorf("infermedica: Unexpected value for triage level: %q", x)
	}
}

// Triage estimates triage level based on the provided patient information.
func (a *App) Triage(tr TriageReq) (*TriageRes, error) {
	if tr.Sex.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: Unexpected value for Sex")
	}
	req, err := a.prepareRequest("POST", "triage", tr)
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

	// Check response
	err = checkResponse(res)
	if err != nil {
		return nil, err
	}

	var r TriageRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

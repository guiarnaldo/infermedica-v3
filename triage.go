package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type TriageRes struct {
	TriageLevel                TriageLevel `json:"triage_level"`                // A classification of the case provided
	Serious                    []Serious   `json:"serious"`                     // A list of serious observations
	TeleconsultationApplicable bool        `json:"teleconsultation_applicable"` // The teleconsultation_applicable flag has been deprecated and will stop being supported in the near future. This functionality has been improved and extended upon and is now available via the /recommend_specialist
	RootCause                  string      `json:"root_cause"`                  // A root cause that explains the internal rationale of the underlying triage algorithm
}

type Serious struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CommonName  string `json:"common_name"`
	Seriousness string `json:"seriousness"`
	IsEmergency bool   `json:"is_emergency"` // IsEmergency is redundant with Seriousness field and should be treated as deprecated
}

type RootCause string

const (
	RootCauseEmergencyEvidencePresent        RootCause = "emergency_evidence_present"    // Emergency evidence was reported as present
	RootCauseEmergencySeriousEvidencePresent RootCause = "serious_evidence_present"      // Serious evidence was reported as present
	RootCauseEmergencyConditionLikely        RootCause = "emergency_condition_likely"    // Life-threatening condition is likely enough to recommend emergency triage
	RootCauseEmergencyConditionPossible      RootCause = "emergency_condition_possible"  // Life-threatening condition is possible
	RootCauseConsultationConditionLikely     RootCause = "consultation_condition_likely" // At least one condition which requires medical consultation is likely
	RootCauseSelfCareSufficient              RootCause = "self_care_sufficient"          // No identified reason for medical evaluation
	RootCauseDiagnosisUnknown                RootCause = "diagnosis_unknown"             // Assessment not possible
)

type Seriousness string

const (
	SeriousnessSerious                   Seriousness = "serious"
	SeriousnessSeriousEmergency          Seriousness = "emergency"
	SeriousnessSeriousEmergencyAmbulance Seriousness = "emergency_ambulance"
)

type TriageLevel string

const (
	TriageLevelEmergencyAmbulance TriageLevel = "emergency_ambulance" // T the reported symptoms are very serious and the patient may require emergency care. The patient should call an ambulance right now
	TriageLevelEmergency          TriageLevel = "emergency"           // The reported evidence appears serious and the patient should go to an emergency department. If the patient can't get to the nearest emergency department, they should call an ambulance
	TriageLevelConsultation24     TriageLevel = "consultation_24"     // The patient should see a doctor within 24 hours. If the symptoms suddenly get worse, the patient should go to the nearest emergency department
	TriageLevelConsultation       TriageLevel = "consultation"        // The patient may require medical evaluation and may need to schedule an appointment with a doctor. If the symptoms get worse, the patient should see a doctor immediately
	TriageLevelSelfCare           TriageLevel = "self_care"           // The declared symptoms may not require a medical evaluation and they usually resolve on their own. Patients should observe their symptoms and consult a doctor if the symptoms get worse or new ones appear
)

func (s TriageLevel) Ptr() *TriageLevel { return &s }
func (s TriageLevel) String() string    { return string(s) }

func (s *TriageLevel) IsValid() error {
	_, err := TriageLevelFromString(s.String())
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
	case "consultation24":
		return TriageLevelConsultation24, nil
	case "consultation":
		return TriageLevelConsultation, nil
	case "self_care":
		return TriageLevelSelfCare, nil
	default:
		return "", fmt.Errorf("infermedica: Unexpected value for triage level: %q", x)
	}
}

func (a *App) Triage(tr ObservationReq) (*TriageRes, error) {
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

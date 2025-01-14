package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func (p *Prevalence) IsValid() error {
	_, err := PrevalenceFromString(string(*p))
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
		return "", fmt.Errorf("infermedica: unexpected value for Prevalence: %q", x)
	}
}

type Acuteness string

const (
	AcutenessChronic                  Acuteness = "chronic"
	AcutenessChronicWithExacerbations Acuteness = "chronic_with_exacerbations"
	AcutenessAcutePotentiallyChronic  Acuteness = "acute_potentially_chronic"
	AcutenessAcute                    Acuteness = "acute"
)

func (a *Acuteness) IsValid() error {
	_, err := AcutenessFromString(string(*a))
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
		return "", fmt.Errorf("infermedica: unexpected value for Acuteness: %q", x)
	}
}

type Severity string

const (
	SeverityMild     Severity = "mild"
	SeverityModerate Severity = "moderate"
	SeveritySevere   Severity = "severe"
)

func (s *Severity) IsValid() error {
	_, err := SeverityFromString(string(*s))
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
		return "", fmt.Errorf("infermedica: unexpected value for Severity: %q", x)
	}
}

type ConditionRes struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	CommonName string   `json:"common_name"`
	SexFilter  string   `json:"sex_filter"`
	Categories []string `json:"categories"`
	Prevalence string   `json:"prevalence"`
	Acuteness  string   `json:"acuteness"`
	Severity   string   `json:"severity"`
	Extras     struct {
		Hint      string `json:"hint"`
		Icd10Code string `json:"icd10_code"`
	} `json:"extras"`
}

func (a *App) Conditions(age Age, enableTriage3 bool) (*[]ConditionRes, error) {
	req, err := a.prepareRequest("GET", "conditions?age.value="+strconv.Itoa(age.Value)+"&age.unit"+string(age.Unit)+"&enableTriage3="+strconv.FormatBool(enableTriage3), nil)
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
	var r []ConditionRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (a *App) ConditionByID(id string, age Age, enableTriage3 bool) (*ConditionRes, error) {
	req, err := a.prepareGETRequest("conditions/" + id + "?age.value=" + strconv.Itoa(age.Value) + "&age.unit" + string(age.Unit) + "&enableTriage3=" + strconv.FormatBool(enableTriage3))
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
	var r ConditionRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

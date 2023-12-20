package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RecommendSpecialistReq struct {
	Sex         Sex                           `json:"sex"`
	Age         Age                           `json:"age"`
	EvaluatedAt string                        `json:"evaluated_at,omitempty"`
	Evidences   []Evidence                    `json:"evidence,omitempty"`
	Extras      *RecommendSpecialistReqExtras `json:"extras,omitempty"`
}

type RecommendSpecialistReqExtras struct {
	EnableSymptomDuration bool              `json:"enable_symptom_duration,omitempty"` // This flag enables questions of the type duration which contain a new field EvidenceID
	SpecialistMapping     map[string]string `json:"specialist_mapping,omitempty"`
}

type RecommendSpecialistRes struct {
	RecommendedSpecialist struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"recommended_specialist"`
	RecommendedChannel string `json:"recommended_channel"`
}

type RecommendedChannel string

const (
	RecommendedChannelPersonalVisit         RecommendedChannel = "personal_visit"
	RecommendedChannelVideoTeleconsultation RecommendedChannel = "video_teleconsultation"
	RecommendedChannelAudioTeleconsultation RecommendedChannel = "audio_teleconsultation"
	RecommendedChannelTextTeleconsultation  RecommendedChannel = "text_teleconsultation"
)

func (a *App) RecommendSpecialist(tr RecommendSpecialistReq) (*RecommendSpecialistRes, error) {
	if tr.Sex.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: unexpected value for Sex")
	}
	req, err := a.prepareRequest("POST", "recommend_specialist", tr)
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

	var r RecommendSpecialistRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

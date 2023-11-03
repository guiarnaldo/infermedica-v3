package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RecommendSpecialistRes struct {
	RecommendedSpecialist struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"recommended_specialist"`
	RecommendedChannel string `json:"recommended_channel"`
}

type RecommendedChannel string

const (
	RecommendedChannelPersonalVisit         RecommendedChannel = "personal_visit"         // Personal visit (In-person visit)
	RecommendedChannelVideoTeleconsultation RecommendedChannel = "video_teleconsultation" // Video teleconsultation (Video consultation)
	RecommendedChannelAudioTeleconsultation RecommendedChannel = "audio_teleconsultation" // Audio teleconsultation (Telephone or any other consultation with audio only)
	RecommendedChannelTextTeleconsultation  RecommendedChannel = "text_teleconsultation"  // Text teleconsultation (Chat)
)

func (a *App) RecommendSpecialist(tr ObservationReq) (*RecommendSpecialistRes, error) {
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

package gortc

import (
	"encoding/json"

	"github.com/pion/webrtc/v3"
)

type ICECandidateInit webrtc.ICECandidateInit

func (ice *ICECandidateInit) ToJSON() string {
	str, err := json.Marshal(ice)
	if err != nil {
		return ""
	}

	return string(str)
}

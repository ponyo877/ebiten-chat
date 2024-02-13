package folks

import (
	"encoding/json"
	"log"

	"github.com/pion/webrtc/v3"
)

func (g *Game) syncWebRTC() {
	// Create new PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		log.Print(err)
		return
	}

	// When this frame returns close the PeerConnection
	defer peerConnection.Close() //nolint

	// Accept one audio and one video track incoming
	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			log.Print(err)
			return
		}
	}

	// Trickle ICE. Emit server candidate to client
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		data, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Println(err)
			return
		}
		if err := g.wss.Candidate(data); err != nil {
			log.Println(err)
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		// Trackを受信したときの処理
	})

	// message := &Signal{}

	// for {
	// 	raw, err := g.wss.ReadMessage()
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	} else if err := json.Unmarshal(raw, &message); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	switch message.Event {
	// 	case "candidate":
	// 		candidate := webrtc.ICECandidateInit{}
	// 		if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
	// 			log.Println(err)
	// 			return
	// 		}

	// 		if err := peerConnection.AddICECandidate(candidate); err != nil {
	// 			log.Println(err)
	// 			return
	// 		}
	// 	case "answer":
	// 		answer := webrtc.SessionDescription{}
	// 		if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
	// 			log.Println(err)
	// 			return
	// 		}

	// 		if err := peerConnection.SetRemoteDescription(answer); err != nil {
	// 			log.Println(err)
	// 			return
	// 		}
	// 	}
	// }
}

package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/opus"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"
	"github.com/ponyo877/folks-ui/entity"
	"nhooyr.io/websocket" // wasm対応のためgorilla/websockeから変更
	"nhooyr.io/websocket/wsjson"
)

type WebConnection struct {
	wcon *websocket.Conn
	pcon *webrtc.PeerConnection
}

func NewWebConnection(host, path string) (*WebConnection, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: path}
	wcon, _, err := websocket.Dial(context.Background(), u.String(), nil)
	if err != nil {
		fmt.Printf("Websocketサーバへの接続に失敗しました: %v\n", err)
		return nil, err
	}
	opusParams, err := opus.NewParams()
	if err != nil {
		panic(err)
	}
	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithAudioEncoders(&opusParams),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	pcon, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameFormat = prop.FrameFormat(frame.FormatI420)
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
		},
		Audio: func(c *mediadevices.MediaTrackConstraints) {
		},
		Codec: codecSelector,
	})
	if err != nil {
		panic(err)
	}
	for _, track := range s.GetTracks() {
		track.OnEnded(func(err error) {
			fmt.Printf("Track (ID: %s) ended with error: %v\n", track.ID(), err)
		})

		_, err = pcon.AddTransceiverFromTrack(track,
			webrtc.RtpTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionSendonly,
			},
		)
		if err != nil {
			panic(err)
		}
	}
	return &WebConnection{wcon, pcon}, nil
}

func (w *WebConnection) Send(message *entity.Message) error {
	return wsjson.Write(context.Background(), w.wcon, MarshalMessage(message))
}

func (w *WebConnection) candidate(ICECandidate []byte) error {
	signalPresenter := NewSignalPresenter("candidate", string(ICECandidate))
	if err := wsjson.Write(context.Background(), w.wcon, signalPresenter); err != nil {
		return err
	}
	return nil
}

func (w *WebConnection) answer(answer webrtc.SessionDescription) error {
	signalPresenter, err := NewSignalPresenter("answer", answer)
	if err != nil {
		return err
	}
	if err := wsjson.Write(context.Background(), w.wcon, signalPresenter); err != nil {
		return err
	}
	return nil
}

func (w *WebConnection) WaitCandidate() {
	w.pcon.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		data, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Println(err)
			return
		}
		if err := w.candidate(data); err != nil {
			log.Println(err)
		}
	})
}

func (w *WebConnection) SetAudioConfig() error {
	_, err := w.pcon.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionSendrecv,
	})
	if err != nil {
		return err
	}
	return nil
}

func (w *WebConnection) ReadSignal() (*entity.Signal, error) {
	_, bin, err := w.wcon.Read(context.Background())
	if err != nil {
		return nil, err
	}
	signal, err := UnmarshalSignal(bin)
	if err != nil {
		return nil, err
	}
	return signal, nil
}

func (w *WebConnection) SetSignal(signal *entity.Signal) error {
	switch signal.Event() {
	case "candidate":
		return w.SetCandidate(signal)
	case "answer":
		return w.SetAnswer(signal)
	}
	return nil
}

func (w *WebConnection) SetCandidate(signal *entity.Signal) error {
	candidate := webrtc.ICECandidateInit{}
	if err := json.Unmarshal(signal.Data(), &candidate); err != nil {
		return err
	}

	if err := w.pcon.AddICECandidate(candidate); err != nil {
		return err
	}
	return nil
}

func (w *WebConnection) SetAnswer(signal *entity.Signal) error {
	answer := webrtc.SessionDescription{}
	if err := json.Unmarshal(signal.Data(), &answer); err != nil {
		return err
	}
	if err := w.pcon.SetRemoteDescription(answer); err != nil {
		return err
	}
	return nil
}

func (w *WebConnection) ReadTrack() {
	w.pcon.OnTrack(func(t *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		fmt.Printf("OnTrack: %v, %v\n", t, r)
	})
}

func (w *WebConnection) Something() {
	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	Decode(MustReadStdin(), &offer)

	// Set the remote SessionDescription
	if err := w.pcon.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	// Create an answer
	answer, err := w.pcon.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// signalingServer.send(JSON.stringify({ event: 'answer', data: JSON.stringify(answer) }));
	w.

		// Create channel that is blocked until ICE Gathering is complete
		gatherComplete := webrtc.GatheringCompletePromise(w.pcon)

	// Sets the LocalDescription, and starts our UDP listeners
	if err := w.pcon.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete
}

func (w *WebConnection) Receive(f func(*entity.Message)) {
	var messagePresenter MessagePresenter
	for {
		if err := wsjson.Read(context.Background(), w.wcon, &messagePresenter); err != nil {
			fmt.Printf("Messageの読み込みに失敗しました: %v\n", err)
			return
		}
		f(messagePresenter.Unmarshal())
	}
}

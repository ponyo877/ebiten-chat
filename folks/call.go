package folks

import (
	"log"
)

func (g *Game) syncWebRTC() {
	if err := g.webc.SetAudioConfig(); err != nil {
		log.Println(err)
		return
	}
	g.webc.WaitCandidate()
	g.webc.ReadTrack()
	g.webc.Something()

	for {
		signal, err := g.webc.ReadSignal()
		if err != nil {
			log.Println(err)
			return
		}
		if err := g.webc.SetSignal(signal); err != nil {
			log.Println(err)
			return
		}
	}
}

package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
)

type GameVariant uint8

func (gv GameVariant) String() string {
	switch gv {
	case TexasHoldem:
		return "TEXAS HOLD'EM"
	case Other:
		return "other"
	default:
		return "unknown"
	}
}

const (
	TexasHoldem GameVariant = iota
	Other
)

type ServerConfig struct {
	Version     string
	ListenAddr  string
	GameVariant GameVariant
}

type Server struct {
	ServerConfig
	transport *TCPTransport
	peers     map[net.Addr]*Peer
	addPeer   chan *Peer
	delPeer   chan *Peer
	msgCh     chan *Message
	gameState *GameState
}

func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
		gameState:    NewGameState(),
	}

	tr := NewTCPTransport(s.ListenAddr)
	s.transport = tr

	tr.AddPeer = s.addPeer
	tr.DelPeer = s.delPeer
	return s
}

func (s *Server) Start() {
	go s.loop()

	logrus.WithFields(logrus.Fields{
		"port":    s.ListenAddr,
		"variant": s.GameVariant,
	}).Info("Started New Game Server")

	s.transport.ListenAndAccept()
}

func (s *Server) SendHandShake(p *Peer) error {
	hs := &Handshake{
		GameVariant: s.GameVariant,
		Version:     s.Version,
		GameStatus:  s.gameState.gameStatus,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		return err
	}

	return p.Send(buf.Bytes())
}

// TODO:right now we have some redundant code in registering new peers to the game network
// maybe construct a new peer and handshake protocol after registering a plain connection
func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn:     conn,
		outbound: true,
	}

	s.addPeer <- peer
	return s.SendHandShake(peer)
}

func (s *Server) loop() {
	for {
		select {

		case peer := <-s.delPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("Player disconnected")

			delete(s.peers, peer.conn.RemoteAddr())

		//If a new peer connects to the server we send our handshake message
		//and wait for his reply.
		case peer := <-s.addPeer:
			if err := s.handshake(peer); err != nil {
				logrus.Errorf("%s:Handshake Failed with incomming player: %s", s.ListenAddr, err)
				peer.conn.Close()
				delete(s.peers, peer.conn.RemoteAddr())
				continue
			}

			//TODO:check max playera and other game state logic.
			go peer.ReadLoop(s.msgCh)

			if !peer.outbound {
				if err := s.SendHandShake(peer); err != nil {
					logrus.Errorf("Failed to send Handshake with peer: %s", err)
					peer.conn.Close()
					delete(s.peers, peer.conn.RemoteAddr())
					continue
				}
			}

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("Handshake Successful: New Player Connected")

			s.peers[peer.conn.RemoteAddr()] = peer

		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}

func (s *Server) handshake(p *Peer) error {
	hs := &Handshake{}
	if err := gob.NewDecoder(p.conn).Decode(hs); err != nil {
		return err
	}

	if s.GameVariant != hs.GameVariant {
		return fmt.Errorf("GameVariant does not match %s", hs.GameVariant)
	}

	if s.Version != hs.Version {
		return fmt.Errorf("invalid Version %s", hs.Version)
	}

	logrus.WithFields(logrus.Fields{
		"peer":       p.conn.RemoteAddr(),
		"version":    hs.Version,
		"variant":    hs.GameVariant,
		"gameStatus": hs.GameStatus,
	}).Info("recieved handshake")

	return nil
}

func (s *Server) handleMessage(msg *Message) error {
	fmt.Printf("%+v\n", msg)
	return nil
}

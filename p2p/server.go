package p2p

import (
	"bytes"
	"encoding/binary"
	// "encoding/gob"
	// "errors"
	"fmt"
	"io"
	"net"
	"github.com/sirupsen/logrus"
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
	}

	buf := new(bytes.Buffer)
	// if err := gob.NewEncoder(buf).Encode(hs); err != nil {
	// 	return err
	// }

	if err:=hs.Encode(buf);err!=nil{
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
		conn: conn,
	}

	s.addPeer <- peer
	return peer.Send([]byte(s.Version))
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
			go s.SendHandShake(peer)
			if err := s.handshake(peer); err != nil {
				logrus.Errorf("Handshake Failed with incomming player: %s", err)
				continue
			}

			//TODO:check max playera and other game state logic.
			go peer.ReadLoop(s.msgCh)

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

type Handshake struct {
	Version     string
	GameVariant GameVariant
}

func (hs *Handshake) Encode(w io.Writer) error {
	if err:=binary.Write(w, binary.LittleEndian, []byte(hs.Version));err != nil{
		return err
	}

	return binary.Write(w,binary.LittleEndian,hs.GameVariant)
}

func (hs *Handshake) Decode(r io.Reader) error {
	if err:=binary.Read(r,binary.LittleEndian,[]byte(hs.Version));err!=nil{
		return err
	}

	return binary.Read(r,binary.LittleEndian,hs.GameVariant)
}

func (s *Server) handshake(p *Peer) error {
	hs := &Handshake{}
	// if err := gob.NewDecoder(p.conn).Decode(hs); err != nil {
	// 	return err
	// }

	if err :=hs.Decode(p.conn);err!=nil{
		return err
	} 

	if s.GameVariant != hs.GameVariant {
		return fmt.Errorf("invalid GameVariant %s", hs.GameVariant)
	}

	if s.Version != hs.Version {
		return fmt.Errorf("invalid Version %s", hs.Version)
	}

	logrus.WithFields(logrus.Fields{
		"peer":    p.conn.RemoteAddr(),
		"version": hs.Version,
		"variant": hs.GameVariant,
	}).Info("recieved handshake")

	return nil
}

func (s *Server) handleMessage(msg *Message) error {
	fmt.Printf("%+v\n", msg)
	return nil
}

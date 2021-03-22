package firewall

import (
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/multiformats/go-multiaddr"
	"github.com/zarbchain/zarb-go/logger"
)

//Firewall implments notifee interface to inject policies and rules
//it requires
type Firewall struct {
	ps     PeerSelector
	logger *logger.Logger
}

//NewFirewall will spinup a new firewall with default data
func NewFirewall(ps PeerSelector, lg *logger.Logger) (*Firewall, error) {
	if ps == nil {
		ps = NewSelector()
	}

	fw := &Firewall{
		ps:     ps,
		logger: lg,
	}
	if lg == nil {
		fw.logger = logger.NewLogger("_network_firewall", *fw)
	}

	fw.logger.Info("Network Firewall started")

	if err := fw.ps.Load(); err != nil {
		fw.logger.Debug("PeerSelector Couldn't load policy configuration", "err", err)
		fw.logger.Warn("couldn't Load firewall policy may allow any one to connect")
		// return nil, err
	}

	return fw, nil
}

/***************************************Notifee Interface*******************************************/
// called when network starts listening on an addr
func (fw *Firewall) Listen(n network.Network, m multiaddr.Multiaddr) {
	// libp2p.ConnectionGater(connmgr.ConnectionGater)
	//COntinure with
	fw.logger.Trace("Network is up Listening", "Address", m.String())
}

// called when network stops listening on an addr
func (fw *Firewall) ListenClose(n network.Network, m multiaddr.Multiaddr) {
	fw.logger.Trace("Network Closing Connection", "Address", m.String())
}

// called when a connection opened
func (fw *Firewall) Connected(n network.Network, c network.Conn) {
	fw.logger.Info("New Node Tring To Connected ", "peerid", c.RemotePeer(), "Addrs", c.RemoteMultiaddr())

	if fw.ps.CanConnect(n.Peerstore().PeerInfo(c.RemotePeer())) {
		fw.logger.Debug("Node Disconnected Due to selector Policy", "Peerid", c.RemotePeer())
		n.ClosePeer(c.RemotePeer())
		c.Close()
	}
}

// called when a connection closed
func (fw *Firewall) Disconnected(n network.Network, c network.Conn) {
	fw.logger.Debug("Node Disconnected", "Peerid", c.RemotePeer(), "Addrs", c.RemoteMultiaddr())
}

// called when a stream opened
func (fw *Firewall) OpenedStream(n network.Network, s network.Stream) {
	fw.logger.Trace("Stream Opened", "Id", s.ID(), "Peerid", s.Conn().RemotePeer(), "stats", s.Stat())
}

// called when a stream closed
func (fw *Firewall) ClosedStream(n network.Network, s network.Stream) {
	fw.logger.Trace("Stream Closed", "Id", s.ID(), "Peerid", s.Conn().RemotePeer(), "stats", s.Stat())
}

/***************************************Notifee Interface*******************************************/

func (fw *Firewall) Stop() error {
	fw.logger.Info("Cosing firewall ...")
	err := fw.ps.Store()
	if err != nil {
		fw.logger.Error("PeerSelector Couldn't Store policy configuration", "err", err)
	}
	return err
}

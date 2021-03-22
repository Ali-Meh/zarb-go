package firewall

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/util"
)

//PeerSelector filters wether or not a peer is allowed and trusted to be connected or not
//its main firewall interface on p2p layer to be implemented
type PeerSelector interface {
	//check to see if Peer is allowed to connect
	CanConnect(peer.AddrInfo) bool
	//Store Selection related Configuration it will be called upon closing process
	Store() error
	//Load Selection related Configuration it will be invoced after firewall creation
	Load() error
}

const configFile = "selector.conf"

func init() {
	// This type must match exactly what youre going to be using,
	// down to whether or not its a pointer
	gob.Register(&peer.IDSlice{})
	gob.Register(&[]peer.ID{})
}

type DefaultSelector struct {
	trustedNodes *peer.Set
}

func NewSelector( /* , logger *logger.Logger */ ) *DefaultSelector {
	return &DefaultSelector{
		trustedNodes: peer.NewSet(),
	}
}

func (ps *DefaultSelector) encodeConfig() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(ps.trustedNodes.Peers())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decodeConfig(t io.Reader) ([]peer.ID, error) {
	conf := make([]peer.ID, 0)
	decoder := gob.NewDecoder(t)
	if err := decoder.Decode(&conf); err != nil {
		return nil, err
	}
	return conf, nil
}

/***************************************PeerSelector Interface*******************************************/

func (ps *DefaultSelector) CanConnect(peer peer.AddrInfo) bool {
	if ps.trustedNodes.Size() > 0 {
		return !ps.trustedNodes.Contains(peer.ID)
	}
	return true
}

func (ps *DefaultSelector) Store() error {
	buf, err := ps.encodeConfig()
	if err != nil {
		return err
	}
	return util.WriteFile(configFile, buf)
}

func (ps *DefaultSelector) Load() error {
	if !util.PathExists(configFile) {
		return nil
	}
	buf, err := util.ReadFile(configFile)
	if err != nil {
		return err
	}

	peerids, err := decodeConfig(bytes.NewReader(buf))
	if err != nil {
		return err
	}

	for _, pid := range peerids {
		ps.trustedNodes.TryAdd(pid)
	}

	return nil
}

/***************************************PeerSelector Interface*******************************************/

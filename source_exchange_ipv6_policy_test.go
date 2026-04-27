package goed2k

import (
	"net"
	"testing"

	"github.com/goed2k/core/protocol"
)

func TestPeersForSourceExchangeIncludesIPv6DialOnlyPeer(t *testing.T) {
	st := NewSettings()
	s := NewSession(st)
	tf, err := NewTransfer(s, AddTransferParams{
		Hash:       protocol.MustHashFromString("31D6CFE0D16AE931B73C59D7E0C089C0"),
		CreateTime: 1,
		Size:       1000,
	})
	if err != nil {
		t.Fatal(err)
	}
	ep, err := protocol.EndpointFromString("192.0.2.1", 4662)
	if err != nil {
		t.Fatal(err)
	}
	v6 := NewPeerFromTCPAddr(&net.TCPAddr{IP: net.ParseIP("2001:db8::2"), Port: 4662}, true, 0)
	_, _ = tf.policy.AddPeer(v6)
	list := tf.policy.PeersForSourceExchange(ep, 10)
	if len(list) != 1 {
		t.Fatalf("want 1 ipv6 dial peer, got %d", len(list))
	}
}

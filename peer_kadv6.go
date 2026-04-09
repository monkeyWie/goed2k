package goed2k

import (
	kadv6 "github.com/monkeyWie/goed2k/protocol/kadv6"
)

// PeerFromKADV6SearchEntry 将 KADV6 SearchEntry 中的 TCP 源转为 Policy 用 Peer（Connectable=true）。
// 若条目不含有效 IPv6 源地址则返回 false。供 KADV6Tracker 接入后调用。
func PeerFromKADV6SearchEntry(se kadv6.SearchEntry, sourceFlag int) (Peer, bool) {
	tcp, ok := se.SourceAddr()
	if !ok {
		return Peer{}, false
	}
	return NewPeerFromTCPAddr(tcp, true, sourceFlag), true
}

// AddPeerFromKADV6Search 将单条 KADV6 搜索结果并入当前任务策略表（去重规则同 AddPeer）。
func (t *Transfer) AddPeerFromKADV6Search(entry kadv6.SearchEntry) (bool, error) {
	if t == nil {
		return false, NewError(IllegalArgument)
	}
	p, ok := PeerFromKADV6SearchEntry(entry, int(PeerDHT))
	if !ok {
		return false, nil
	}
	return t.policy.AddPeer(p)
}

package goCache

// PeerPicker is a interface for locating the peer
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter用于查询缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

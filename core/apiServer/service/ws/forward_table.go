package ws

type wrapMsg struct {
	msgID      string
	msg        []byte
	msgSrc     string
	msgDes     string
	msgRecount int
	msgErr     error
}

func NewWrapMsg(msgID int64, msg []byte) {

}

//func (w *wrapMsg) Decode(msg []byte) ([]byte, error) {
//	return utils.UPacket(w.msg)
//}

//func (w *wrapMsg) Encode(msg []byte) ([]byte, error) {
//	return utils.Packet(msg)
//}

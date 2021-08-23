
package mc_util

import (
	bufio "bufio"
	bytes "bytes"
	errors "errors"
	io    "io"
)

type Writer interface{
	BeginWritePacket(id int32)(n int, err error)
	WriteInt8(num int8)(n int, err error)
	WriteInt16(num int16)(n int, err error)
	WriteInt32(num int32)(n int, err error)
	WriteInt64(num int64)(n int, err error)
	WriteUint8(num uint8)(n int, err error)
	WriteUint16(num uint16)(n int, err error)
	WriteUint32(num uint32)(n int, err error)
	WriteUint64(num uint64)(n int, err error)
	WriteVarInt32(num int32)(n int, err error)
	WriteVarInt64(num int64)(n int, err error)
	WriteString(str string)(n int, err error)
	EndWritePacket()(leng int32, n int, err error)
}

type writer struct{
	bw *bufio.Writer
	buf *bytes.Buffer
}

func NewWriter(w io.Writer)(*writer){
	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}
	return &writer{
		bw: bw,
		buf: nil,
	}
}

func (w *writer)BeginWritePacket(id int32)(n int, err error){
	if w.buf != nil {
		return 0, errors.New("PACKET_NOT_SEND_ERR")
	}
	w.buf = bytes.NewBuffer([]byte{})
	return w.WriteVarInt32(id)
}

func (w *writer)WriteInt8(num int8)(n int, err error){
	return w.buf.Write(EncodeInt8(num))
}

func (w *writer)WriteInt16(num int16)(n int, err error){
	return w.buf.Write(EncodeInt16(num))
}

func (w *writer)WriteInt32(num int32)(n int, err error){
	return w.buf.Write(EncodeInt32(num))
}

func (w *writer)WriteInt64(num int64)(n int, err error){
	return w.buf.Write(EncodeInt64(num))
}

func (w *writer)WriteUint8(num uint8)(n int, err error){
	return w.buf.Write(EncodeUint8(num))
}

func (w *writer)WriteUint16(num uint16)(n int, err error){
	return w.buf.Write(EncodeUint16(num))
}

func (w *writer)WriteUint32(num uint32)(n int, err error){
	return w.buf.Write(EncodeUint32(num))
}

func (w *writer)WriteUint64(num uint64)(n int, err error){
	return w.buf.Write(EncodeUint64(num))
}

func (w *writer)WriteVarInt32(num int32)(n int, err error){
	return w.buf.Write(EncodeVarInt32(num))
}

func (w *writer)WriteVarInt64(num int64)(n int, err error){
	return w.buf.Write(EncodeVarInt64(num))
}

func (w *writer)WriteString(str string)(n int, err error){
	return w.buf.Write(EncodeString(str))
}

func (w *writer)EndWritePacket()(leng int32, n int, err error){
	if w.buf == nil {
		return 0, 0, errors.New("PACKET_NOT_BEGIN_WRITE_ERR")
	}
	leng = (int32)(w.buf.Len())
	w.bw.Write(EncodeVarInt32(leng))

	n, err = w.bw.Write(w.buf.Bytes())
	if err != nil {
		return 0, 0, err
	}
	w.bw.Flush()
	w.buf = nil
	return leng, n, err
}


package mc_util

import (
	bufio "bufio"
	bytes "bytes"
	errors "errors"
	io    "io"
)


type Reader interface{
	BeginReadPacket()(id int32, err error)
	ReadInt8()(num int8, err error)
	ReadInt16()(num int16, err error)
	ReadInt32()(num int32, err error)
	ReadInt64()(num int64, err error)
	ReadUint8()(num uint8, err error)
	ReadUint16()(num uint16, err error)
	ReadUint32()(num uint32, err error)
	ReadUint64()(num uint64, err error)
	ReadVarInt32()(num int32, n int, err error)
	ReadVarInt64()(num int64, n int, err error)
	ReadString()(str string, n int, err error)
	EndReadPacket()(data []byte, err error)
}

type reader struct{
	br *bufio.Reader
	data_buf *bytes.Reader
}

func NewReader(r io.Reader)(*reader){
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	return &reader{
		br: br,
		data_buf: nil,
	}
}

func readVarInt32(r io.ByteReader)(num int32, n int, err error){
	var b byte
	for n = 0; n < 5 ;n++ {
		b, err = r.ReadByte()
		num |= ((int32)(b & 0x7f) << (n * 7))
		if err != nil {
			return 0, 0, err
		}
		if b & 0x80 == 0 {
			return num, n + 1, nil
		}
	}
	return 0, n + 1, errors.New("VAR_INT_TOO_LONG_ERR")
}

func readVarInt64(r io.ByteReader)(num int64, n int, err error){
	var b byte
	for n = 0; n < 10 ;n++ {
		b, err = r.ReadByte()
		num |= ((int64)(b & 0x7f) << (n * 7))
		if err != nil {
			return 0, 0, err
		}
		if b & 0x80 == 0 {
			return num, n + 1, nil
		}
	}
	return 0, n + 1, errors.New("VAR_INT_TOO_LONG_ERR")
}

func (r *reader)BeginReadPacket()(id int32, err error){
	var leng int32
	leng, _, err = readVarInt32(r.br)
	if err != nil {
		return 0, err
	}
	var n int
	id, n, err = readVarInt32(r.br)
	if err != nil {
		return 0, err
	}
	leng -= int32(n)
	var buf []byte = make([]byte, leng)
	for ind := 0; ind < (int)(leng) ;ind += n {
		n, err = r.br.Read(buf[ind:])
		if err != nil {
			return 0, err
		}
	}
	r.data_buf = bytes.NewReader(buf)
	return id, nil
}

func (r *reader)ReadInt8()(num int8, err error){
	var buf []byte = make([]byte, 1)
	_, err = r.data_buf.Read(buf)
	if err != nil {
		return 0, err
	}
	return DecodeInt8(buf), nil
}

func (r *reader)ReadInt16()(num int16, err error){
	var buf []byte = make([]byte, 2)
	_, err = r.data_buf.Read(buf)
	if err != nil {
		return 0, err
	}
	return DecodeInt16(buf), nil
}

func (r *reader)ReadInt32()(num int32, err error){
	var buf []byte = make([]byte, 4)
	_, err = r.data_buf.Read(buf)
	if err != nil {
		return 0, err
	}
	return DecodeInt32(buf), nil
}

func (r *reader)ReadInt64()(num int64, err error){
	var buf []byte = make([]byte, 8)
	_, err = r.data_buf.Read(buf)
	if err != nil {
		return 0, err
	}
	return DecodeInt64(buf), nil
}

func (r *reader)ReadUint8()(num uint8, err error){
	var buf []byte = make([]byte, 1)
	_, err = r.data_buf.Read(buf)
	if err != nil {
		return 0, err
	}
	return DecodeUint8(buf), nil
}

func (r *reader)ReadUint16()(num uint16, err error){
	var buf []byte = make([]byte, 2)
	_, err = r.data_buf.Read(buf)
	if err != nil {
		return 0, err
	}
	return DecodeUint16(buf), nil
}

func (r *reader)ReadUint32()(num uint32, err error){
	var buf []byte = make([]byte, 4)
	_, err = r.data_buf.Read(buf)
	if err != nil {
		return 0, err
	}
	return DecodeUint32(buf), nil
}

func (r *reader)ReadUint64()(num uint64, err error){
	var buf []byte = make([]byte, 8)
	_, err = r.data_buf.Read(buf)
	if err != nil {
		return 0, err
	}
	return DecodeUint64(buf), nil
}

func (r *reader)ReadVarInt32()(num int32, n int, err error){
	return readVarInt32(r.data_buf)
}

func (r *reader)ReadVarInt64()(num int64, n int, err error){
	return readVarInt64(r.data_buf)
}

func (r *reader)ReadString()(str string, n int, err error){
	if r.data_buf == nil {
		return "", 0, errors.New("PACKET_NOT_BEGIN_READ_ERR")
	}
	var leng int32
	var n0 int
	leng, n0, err = r.ReadVarInt32()
	if err != nil {
		return "", 0, err
	}
	n = n0
	var buf []rune = make([]rune, leng)
	for ind := 0; ind < (int)(leng) ;ind++{
		buf[ind], n0, err = r.data_buf.ReadRune()
		if err != nil {
			if err == io.EOF{
				return "", 0, errors.New("PACKET_LENGTH_OUT_ERROR")
			}
			return "", 0, err
		}
		n += n0
	}
	return (string)(buf), n, err
}

func (r *reader)EndReadPacket()(data []byte, err error){
	data = make([]byte, r.data_buf.Len())
	var n int
	n, err = r.data_buf.Read(data)
	if err != nil {
		return nil, err
	}
	r.data_buf = nil
	return data[:n], nil
}

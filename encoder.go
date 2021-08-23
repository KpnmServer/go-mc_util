
package mc_util

func EncodeBool(b bool)([]byte){
	if b {
		return []byte{0x01}
	}
	return []byte{0x00}
}

func EncodeUint8(n uint8)([]byte){
	return []byte{(byte)(n)}
}

func EncodeUint16(n uint16)([]byte){
	return []byte{
		(byte)((n >> 8) & 0xff),
		(byte)(n & 0xff),
	}
}

func EncodeUint32(n uint32)([]byte){
	return []byte{
		(byte)((n >> 24)),
		(byte)((n >> 16) & 0xff),
		(byte)((n >> 8) & 0xff),
		(byte)(n & 0xff),
	}
}

func EncodeUint64(n uint64)([]byte){
	return []byte{
		(byte)((n >> 56)),
		(byte)((n >> 48) & 0xff),
		(byte)((n >> 40) & 0xff),
		(byte)((n >> 32) & 0xff),
		(byte)((n >> 24) & 0xff),
		(byte)((n >> 16) & 0xff),
		(byte)((n >> 8) & 0xff),
		(byte)(n & 0xff),
	}
}

func EncodeInt8(n int8)([]byte){
	return EncodeUint8((uint8)(n))
}

func EncodeInt16(n int16)([]byte){
	return EncodeUint16((uint16)(n))
}

func EncodeInt32(n int32)([]byte){
	return EncodeUint32((uint32)(n))
}

func EncodeInt64(n int64)([]byte){
	return EncodeUint64((uint64)(n))
}

func EncodeVarInt32(num int32)(bts []byte){
	n := (uint32)(num)
	bts = make([]byte, 0)
	var b byte
	for{
		b = (byte)(n & 0x7f)
		n >>= 7
		if n != 0 {
			b |= 0x80
		}
		bts = append(bts, b)
		if n == 0 {
			break
		}
	}
	return bts
}

func EncodeVarInt64(num int64)(bts []byte){
	n := (uint64)(num)
	bts = make([]byte, 0)
	var b byte
	for{
		b = (byte)(n & 0x7f)
		n >>= 7
		if n != 0 {
			b |= 0x80
		}
		bts = append(bts, b)
		if n == 0 {
			break
		}
	}
	return bts
}

func EncodeString(str string)(bts []byte){
	head := EncodeVarInt32((int32)(len(([]rune)(str))))
	bts = make([]byte, 0, len(str) + len(head))
	bts = append(bts, head...)
	bts = append(bts, ([]byte)(str)...)
	return bts
}

func DecodeUint8(bts []byte)(uint8){
	return (uint8)(bts[0])
}

func DecodeUint16(bts []byte)(uint16){
	return (uint16)(
		((uint16)(bts[0]) << 8) |
		(uint16)(bts[1]))
}

func DecodeUint32(bts []byte)(uint32){
	return (uint32)(
		((uint32)(bts[0]) << 24) |
		((uint32)(bts[1]) << 16) |
		((uint32)(bts[2]) << 8) |
		(uint32)(bts[3]))
}

func DecodeUint64(bts []byte)(uint64){
	return (uint64)(
		((uint64)(bts[0]) << 56) |
		((uint64)(bts[0]) << 48) |
		((uint64)(bts[0]) << 40) |
		((uint64)(bts[0]) << 32) |
		((uint64)(bts[0]) << 24) |
		((uint64)(bts[1]) << 16) |
		((uint64)(bts[2]) << 8) |
		(uint64)(bts[3]))
}

func DecodeInt8(bts []byte)(int8){
	return (int8)(DecodeUint8(bts))
}

func DecodeInt16(bts []byte)(int16){
	return (int16)(DecodeUint16(bts))
}

func DecodeInt32(bts []byte)(int32){
	return (int32)(DecodeUint32(bts))
}

func DecodeInt64(bts []byte)(int64){
	return (int64)(DecodeUint64(bts))
}


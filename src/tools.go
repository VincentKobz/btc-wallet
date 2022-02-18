package src

import "encoding/binary"

func VarIntConverter(count int) []byte {
	if count <= 255 {
		return []byte{byte(count)}
	}

	if count <= 65535 {
		res := []byte{253}
		value := make([]byte, 2)
		binary.LittleEndian.PutUint16(value, uint16(count))
		res = append(res, value...)
		return res
	}

	if count <= 4294967295 {
		res := []byte{254}
		value := make([]byte, 4)
		binary.LittleEndian.PutUint32(value, uint32(count))
		res = append(res, value...)
		return res
	}

	res := []byte{255}
	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, uint64(count))
	res = append(res, value...)
	return res
}

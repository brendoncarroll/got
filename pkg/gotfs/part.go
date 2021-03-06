package gotfs

import (
	"encoding/binary"

	"github.com/brendoncarroll/got/pkg/gotkv"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func (p *Part) marshal() []byte {
	data, err := proto.Marshal(p)
	if err != nil {
		panic(err)
	}
	return data
}

func parsePart(data []byte) (*Part, error) {
	p := &Part{}
	if err := proto.Unmarshal(data, p); err != nil {
		return nil, err
	}
	return p, nil
}

func splitPartKey(k []byte) (p string, offset uint64, err error) {
	if len(k) < 9 {
		return "", 0, errors.Errorf("key too short")
	}
	if k[len(k)-9] != 0x00 {
		return "", 0, errors.Errorf("not part key, no NULL")
	}
	p = string(k[:len(k)-9])
	offset = binary.BigEndian.Uint64(k[len(k)-8:])
	return p, offset, nil
}

func makePartKey(p string, offset uint64) []byte {
	x := []byte(p)
	x = append(x, 0x00)
	x = appendUint64(x, offset)
	return x
}

func fileSpanEnd(p string) []byte {
	return gotkv.PrefixEnd(append([]byte(p), 0x00))
}

func appendUint64(buf []byte, n uint64) []byte {
	nbytes := [8]byte{}
	binary.BigEndian.PutUint64(nbytes[:], n)
	return append(buf, nbytes[:]...)
}

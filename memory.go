package vm

import (
	"errors"
	"fmt"
)

var ErrInvalidAddress = errors.New("invalid address")

type Memory struct {
	mem [64 * 1024]byte
}

func (m *Memory) StoreByte(addr uint16, value byte) error {
	if int(addr) >= len(m.mem) {
		fmt.Errorf("%w: %d", ErrInvalidAddress, addr)
	}

	m.mem[int(addr)] = value
	return nil
}

func (m *Memory) FetchByte(addr uint16) (byte, error) {
	if int(addr) >= len(m.mem) {
		return 0, fmt.Errorf("%w: %d", ErrInvalidAddress, addr)
	}

	return m.mem[addr], nil
}

func (m *Memory) FetchDword(addr uint16) (byte, error) {
	if int(addr)+3 >= len(m.mem) {
		return 0, fmt.Errorf("%w: %d", ErrInvalidAddress, addr)
	}

	return m.mem[addr] |
			m.mem[addr+1]<<8 |
			m.mem[addr+2]<<16 |
			m.mem[addr+3]<<24,
		nil
}

func (m *Memory) StoreDword(addr uint16, value uint32) error {
	if int(addr)+3 >= len(m.mem) {
		return fmt.Errorf("%w: %d", ErrInvalidAddress, addr)
	}

	m.mem[int(addr)] = byte(value & 0xff)
	m.mem[int(addr)+1] = byte((value >> 8) & 0xff)
	m.mem[int(addr)+2] = byte((value >> 16) & 0xff)
	m.mem[int(addr)+3] = byte((value >> 24) & 0xff)

	return nil
}

func (m *Memory) FetchMany(addr uint16, size int) ([]byte, error) {
	if int(addr)+size >= len(m.mem) {
		return nil, fmt.Errorf("%w: %d", ErrInvalidAddress, addr)
	}

	return m.mem[addr : int(addr)+size], nil
}

func (m *Memory) StoreMany(addr uint16, data []byte) error {
	if int(addr)+len(data) >= len(m.mem) {
		return fmt.Errorf("%w: %d", ErrInvalidAddress, addr)
	}

	for i := 0; i < len(data); i++ {
		m.mem[int(addr)+i] = data[i]
	}

	return nil
}

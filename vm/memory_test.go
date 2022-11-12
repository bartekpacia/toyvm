package vm

import (
	"bytes"
	"errors"
	"testing"
)

func TestStoreByte(t *testing.T) {
	testCases := []struct {
		addr    uint16
		value   byte
		want    []byte
		wantErr error
	}{
		{addr: 0, value: 2, want: []byte{2, 0, 0, 0}, wantErr: nil},
		{addr: 3, value: 1, want: []byte{0, 0, 0, 1}, wantErr: nil},
		{addr: 4, value: 0, want: []byte{0, 0, 0, 0}, wantErr: ErrInvalidAddress},
	}

	for _, tc := range testCases {
		memory := &Memory{mem: make([]byte, 4)}

		err := memory.StoreByte(tc.addr, tc.value)
		if !errors.Is(err, tc.wantErr) {
			t.Fatalf("want %#v, got %#v", tc.wantErr, err)
		}

		if !bytes.Equal(memory.mem, tc.want) {
			t.Fatalf("want %v, got %v", tc.want, memory.mem)
		}
	}
}

func TestFetchByte(t *testing.T) {
	testCases := []struct {
		mem     *Memory
		addr    uint16
		want    byte
		wantErr error
	}{
		{mem: &Memory{mem: []byte{2, 0, 0, 0}}, addr: 0, want: 2, wantErr: nil},
		{mem: &Memory{mem: []byte{0, 0, 0, 1}}, addr: 3, want: 1, wantErr: nil},
		{mem: &Memory{mem: []byte{0, 0, 0, 0}}, addr: 5, want: 0, wantErr: ErrInvalidAddress},
	}

	for _, tc := range testCases {
		got, err := tc.mem.FetchByte(tc.addr)
		if !errors.Is(err, tc.wantErr) {
			t.Fatalf("want %#v, got %#v", tc.wantErr, err)
		}

		if got != tc.want {
			t.Fatalf("want %v, got %v", tc.want, got)
		}
	}
}

func TestFetchDword(t *testing.T) {
	testCases := []struct {
		mem     *Memory
		addr    uint16
		want    uint32
		wantErr error
	}{
		{mem: &Memory{mem: []byte{0xd5, 0, 0, 0}}, addr: 0, want: 213, wantErr: nil},
		{mem: &Memory{mem: []byte{0xcc, 0xdc, 0xbd, 0xc}}, addr: 0, want: 213769420, wantErr: nil},
		{mem: &Memory{mem: []byte{0, 0, 0, 0}}, addr: 1, want: 0, wantErr: ErrInvalidAddress},
		{mem: &Memory{mem: []byte{0, 0, 0, 0}}, addr: 5, want: 0, wantErr: ErrInvalidAddress},
	}

	for _, tc := range testCases {
		got, err := tc.mem.FetchDword(tc.addr)
		if !errors.Is(err, tc.wantErr) {
			t.Fatalf("want %#v, got %#v", tc.wantErr, err)
		}

		if got != tc.want {
			t.Fatalf("want %v, got %v", tc.want, got)
		}
	}
}

func TestStoreDword(t *testing.T) {
	testCases := []struct {
		addr    uint16
		value   uint32
		want    []byte
		wantErr error
	}{
		{addr: 0, value: 213, want: []byte{213, 0, 0, 0}, wantErr: nil},
		{addr: 0, value: 213769420, want: []byte{0xcc, 0xdc, 0xbd, 0xc}, wantErr: nil},
		{addr: 1, value: 0, want: []byte{0, 0, 0, 0}, wantErr: ErrInvalidAddress},
		{addr: 5, value: 0, want: []byte{0, 0, 0, 0}, wantErr: ErrInvalidAddress},
	}

	for _, tc := range testCases {
		memory := &Memory{mem: make([]byte, 4)}

		err := memory.StoreDword(tc.addr, tc.value)
		if !errors.Is(err, tc.wantErr) {
			t.Fatalf("want %#v, got %#v", tc.wantErr, err)
		}

		if !bytes.Equal(memory.mem, tc.want) {
			t.Fatalf("want %v, got %v", tc.want, memory.mem)
		}
	}
}

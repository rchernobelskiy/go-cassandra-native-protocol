package primitives

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-cassandra-native-protocol/cassandraprotocol"
	"net"
	"testing"
)

func TestReadByte(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  byte
		remaining []byte
		err       error
	}{
		{"simple byte", []byte{5}, byte(5), []byte{}, nil},
		{"zero byte", []byte{0}, byte(0), []byte{}, nil},
		{"byte with remaining", []byte{5, 1, 2, 3, 4}, byte(5), []byte{1, 2, 3, 4}, nil},
		{"cannot read byte", []byte{}, byte(0), []byte{}, errors.New("not enough bytes to read [byte]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadByte(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteByte(t *testing.T) {
	tests := []struct {
		name      string
		input     byte
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{"simple byte", byte(5), make([]byte, 1), []byte{5}, []byte{}, nil},
		{"byte with remaining", byte(5), make([]byte, 4), []byte{5, 0, 0, 0}, []byte{0, 0, 0}, nil},
		{"cannot write byte", byte(5), []byte{}, []byte{}, []byte{}, errors.New("not enough capacity to write [byte]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteByte(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestReadShort(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  uint16
		remaining []byte
		err       error
	}{
		{"simple short", []byte{0, 5}, uint16(5), []byte{}, nil},
		{"zero short", []byte{0, 0}, uint16(0), []byte{}, nil},
		{"short with remaining", []byte{0, 5, 1, 2, 3, 4}, uint16(5), []byte{1, 2, 3, 4}, nil},
		{"cannot read short", []byte{0}, uint16(0), []byte{0}, errors.New("not enough bytes to read [short]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadShort(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteShort(t *testing.T) {
	tests := []struct {
		name      string
		input     uint16
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{"simple short", uint16(5), make([]byte, LengthOfShort), []byte{0, 5}, []byte{}, nil},
		{"short with remaining", uint16(5), make([]byte, LengthOfShort+1), []byte{0, 5, 0}, []byte{0}, nil},
		{"cannot write short", uint16(5), make([]byte, LengthOfShort-1), []byte{0}, []byte{0}, errors.New("not enough capacity to write [short]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteShort(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestReadInt(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  int32
		remaining []byte
		err       error
	}{
		{"simple int", []byte{0, 0, 0, 5}, int32(5), []byte{}, nil},
		{"zero int", []byte{0, 0, 0, 0}, int32(0), []byte{}, nil},
		{"negative int", []byte{0xff, 0xff, 0xff, 0xff & -5}, int32(-5), []byte{}, nil},
		{"int with remaining", []byte{0, 0, 0, 5, 1, 2, 3, 4}, int32(5), []byte{1, 2, 3, 4}, nil},
		{"cannot read int", []byte{0, 0, 0}, int32(0), []byte{0, 0, 0}, errors.New("not enough bytes to read [int]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadInt(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteInt(t *testing.T) {
	tests := []struct {
		name      string
		input     int32
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{"simple int", int32(5), make([]byte, LengthOfInt), []byte{0, 0, 0, 5}, []byte{}, nil},
		{"negative int", int32(-5), make([]byte, LengthOfInt), []byte{0xff, 0xff, 0xff, 0xff & -5}, []byte{}, nil},
		{"int with remaining", int32(5), make([]byte, LengthOfInt+1), []byte{0, 0, 0, 5, 0}, []byte{0}, nil},
		{"cannot write int", int32(5), make([]byte, LengthOfInt-1), []byte{0, 0, 0}, []byte{0, 0, 0}, errors.New("not enough capacity to write [int]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteInt(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestReadLong(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  int64
		remaining []byte
		err       error
	}{
		{"simple long", []byte{0, 0, 0, 0, 0, 0, 0, 5}, int64(5), []byte{}, nil},
		{"zero long", []byte{0, 0, 0, 0, 0, 0, 0, 0}, int64(0), []byte{}, nil},
		{"negative long", []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff & -5}, int64(-5), []byte{}, nil},
		{"long with remaining", []byte{0, 0, 0, 0, 0, 0, 0, 5, 1, 2, 3, 4}, int64(5), []byte{1, 2, 3, 4}, nil},
		{"cannot read long", []byte{0, 0, 0, 0, 0, 0, 0}, int64(0), []byte{0, 0, 0, 0, 0, 0, 0}, errors.New("not enough bytes to read [long]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadLong(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteLong(t *testing.T) {
	tests := []struct {
		name      string
		input     int64
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{"simple long", int64(5), make([]byte, LengthOfLong), []byte{0, 0, 0, 0, 0, 0, 0, 5}, []byte{}, nil},
		{"negative long", int64(-5), make([]byte, LengthOfLong), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff & -5}, []byte{}, nil},
		{"long with remaining", int64(5), make([]byte, LengthOfLong+1), []byte{0, 0, 0, 0, 0, 0, 0, 5, 0}, []byte{0}, nil},
		{"cannot write long", int64(5), make([]byte, LengthOfLong-1), []byte{0, 0, 0, 0, 0, 0, 0}, []byte{0, 0, 0, 0, 0, 0, 0}, errors.New("not enough capacity to write [long]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteLong(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

const (
	d = byte('d')
	e = byte('e')
	h = byte('h')
	k = byte('k')
	l = byte('l')
	m = byte('m')
	n = byte('n')
	o = byte('o')
	r = byte('r')
	u = byte('u')
	w = byte('w')
)

func TestReadString(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  string
		remaining []byte
		err       error
	}{
		{"simple string", []byte{0, 5, h, e, l, l, o}, "hello", []byte{}, nil},
		{"string with remaining", []byte{0, 5, h, e, l, l, o, 1, 2, 3, 4}, "hello", []byte{1, 2, 3, 4}, nil},
		{"empty string", []byte{0, 0}, "", []byte{}, nil},
		{"non-ASCII string", []byte{
			0, 15, // length
			0xce, 0xb3, 0xce, 0xb5, 0xce, 0xb9, 0xce, 0xac, //γειά
			0x20,                               // space
			0xcf, 0x83, 0xce, 0xbf, 0xcf, 0x85, // σου
		}, "γειά σου", []byte{}, nil},
		{
			"cannot read length",
			[]byte{0},
			"",
			[]byte{0},
			fmt.Errorf("cannot read [string] length: %w", errors.New("not enough bytes to read [short]")),
		},
		{
			"cannot read string",
			[]byte{0, 5, h, e, l, l},
			"",
			[]byte{h, e, l, l},
			errors.New("not enough bytes to read [string] content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadString(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{
			"simple string",
			"hello",
			make([]byte, LengthOfString("hello")),
			[]byte{0, 5, h, e, l, l, o},
			[]byte{},
			nil,
		},
		{"empty string", "", make([]byte, LengthOfString("")), []byte{0, 0}, []byte{}, nil},
		{"non-ASCII string", "γειά σου", make([]byte, LengthOfString("γειά σου")), []byte{
			0, 15, // length
			0xce, 0xb3, 0xce, 0xb5, 0xce, 0xb9, 0xce, 0xac, //γειά
			0x20,                               // space
			0xcf, 0x83, 0xce, 0xbf, 0xcf, 0x85, // σου
		}, []byte{}, nil},
		{
			"string with remaining",
			"hello",
			make([]byte, LengthOfString("hello")+1),
			[]byte{0, 5, h, e, l, l, o, 0},
			[]byte{0},
			nil,
		},
		{
			"cannot write string length",
			"hello",
			make([]byte, LengthOfShort-1),
			[]byte{0},
			[]byte{0},
			fmt.Errorf("cannot write [string] length: %w", errors.New("not enough capacity to write [short]")),
		},
		{
			"cannot write string",
			"hello",
			make([]byte, LengthOfString("hello")-1),
			[]byte{0, 5, 0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			errors.New("not enough capacity to write [string] content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteString(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestReadLongString(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  string
		remaining []byte
		err       error
	}{
		{"simple string", []byte{0, 0, 0, 5, h, e, l, l, o}, "hello", []byte{}, nil},
		{"string with remaining", []byte{0, 0, 0, 5, h, e, l, l, o, 1, 2, 3, 4}, "hello", []byte{1, 2, 3, 4}, nil},
		{"empty string", []byte{0, 0, 0, 0}, "", []byte{}, nil},
		{"non-ASCII string", []byte{
			0, 0, 0, 15, // length
			0xce, 0xb3, 0xce, 0xb5, 0xce, 0xb9, 0xce, 0xac, //γειά
			0x20,                               // space
			0xcf, 0x83, 0xce, 0xbf, 0xcf, 0x85, // σου
		}, "γειά σου", []byte{}, nil},
		{
			"cannot read length",
			[]byte{0, 0, 0},
			"",
			[]byte{0, 0, 0},
			fmt.Errorf("cannot read [long string] length: %w", errors.New("not enough bytes to read [int]")),
		},
		{
			"cannot read string",
			[]byte{0, 0, 0, 5, h, e, l, l},
			"",
			[]byte{h, e, l, l},
			errors.New("not enough bytes to read [long string] content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadLongString(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteLongString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{"simple string", "hello", make([]byte, LengthOfLongString("hello")), []byte{0, 0, 0, 5, h, e, l, l, o}, []byte{}, nil},
		{"empty string", "", make([]byte, LengthOfLongString("")), []byte{0, 0, 0, 0}, []byte{}, nil},
		{"non-ASCII string", "γειά σου", make([]byte, LengthOfLongString("γειά σου")), []byte{
			0, 0, 0, 15, // length
			0xce, 0xb3, 0xce, 0xb5, 0xce, 0xb9, 0xce, 0xac, //γειά
			0x20,                               // space
			0xcf, 0x83, 0xce, 0xbf, 0xcf, 0x85, // σου
		}, []byte{}, nil},
		{
			"string with remaining",
			"hello",
			make([]byte, LengthOfLongString("hello")+1),
			[]byte{0, 0, 0, 5, h, e, l, l, o, 0},
			[]byte{0},
			nil,
		},
		{
			"cannot write string length",
			"hello",
			make([]byte, LengthOfInt-1),
			[]byte{0, 0, 0},
			[]byte{0, 0, 0},
			fmt.Errorf("cannot write [long string] length: %w", errors.New("not enough capacity to write [int]")),
		},
		{
			"cannot write string",
			"hello",
			make([]byte, LengthOfLongString("hello")-1),
			[]byte{0, 0, 0, 5, 0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			errors.New("not enough capacity to write [long string] content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteLongString(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestReadStringList(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  []string
		remaining []byte
		err       error
	}{
		{"empty string list", []byte{0, 0}, []string{}, []byte{}, nil},
		{"singleton string list", []byte{
			0, 1, // length
			0, 5, h, e, l, l, o, // hello
		}, []string{"hello"}, []byte{}, nil},
		{"simple string list", []byte{
			0, 2, // length
			0, 5, h, e, l, l, o, // hello
			0, 5, w, o, r, l, d, // world
		}, []string{"hello", "world"}, []byte{}, nil},
		{"empty elements", []byte{
			0, 2, // length
			0, 0, // elt 1
			0, 0, // elt 2
		}, []string{"", ""}, []byte{}, nil},
		{
			"cannot read list length",
			[]byte{0},
			nil,
			[]byte{0},
			fmt.Errorf("cannot read [string list] length: %w", errors.New("not enough bytes to read [short]")),
		},
		{
			"cannot read list element",
			[]byte{0, 1, 0, 5, h, e, l, l},
			nil,
			[]byte{h, e, l, l},
			fmt.Errorf("cannot read [string list] element: %w", errors.New("not enough bytes to read [string] content")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadStringList(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteStringList(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{
			"empty string list",
			[]string{},
			make([]byte, LengthOfStringList([]string{})),
			[]byte{0, 0},
			[]byte{},
			nil,
		},
		{
			"nil string list",
			nil,
			make([]byte, LengthOfStringList(nil)),
			[]byte{0, 0},
			[]byte{},
			nil,
		},
		{
			"singleton string list",
			[]string{"hello"},
			make([]byte, LengthOfStringList([]string{"hello"})),
			[]byte{
				0, 1, // length
				0, 5, h, e, l, l, o, // hello
			},
			[]byte{},
			nil,
		},
		{
			"simple string list",
			[]string{"hello", "world"},
			make([]byte, LengthOfStringList([]string{"hello", "world"})),
			[]byte{
				0, 2, // length
				0, 5, h, e, l, l, o, // hello
				0, 5, w, o, r, l, d, // world
			},
			[]byte{},
			nil,
		},
		{
			"empty elements",
			[]string{"", ""},
			make([]byte, LengthOfStringList([]string{"", ""})),
			[]byte{
				0, 2, // length
				0, 0, // elt 1
				0, 0, // elt 2
			},
			[]byte{},
			nil,
		},
		{
			"cannot write list length",
			[]string{"hello"},
			make([]byte, LengthOfShort-1),
			[]byte{0},
			[]byte{0},
			fmt.Errorf("cannot write [string list] length: %w", errors.New("not enough capacity to write [short]")),
		},
		{
			"cannot write list element",
			[]string{"hello"},
			make([]byte, LengthOfStringList([]string{"hello"})-1),
			[]byte{0, 1, 0, 5, 0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			fmt.Errorf("cannot write [string list] element: %w", errors.New("not enough capacity to write [string] content")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteStringList(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestReadBytes(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{"empty bytes", []byte{0, 0, 0, 0}, []byte{}, []byte{}, nil},
		{"nil bytes", []byte{0xff, 0xff, 0xff, 0xff}, nil, []byte{}, nil},
		{"singleton bytes", []byte{0, 0, 0, 1, 1}, []byte{1}, []byte{}, nil},
		{"simple bytes", []byte{0, 0, 0, 2, 1, 2}, []byte{1, 2}, []byte{}, nil},
		{
			"cannot read bytes length",
			[]byte{0, 0, 0},
			nil,
			[]byte{0, 0, 0},
			fmt.Errorf("cannot read [bytes] length: %w", errors.New("not enough bytes to read [int]")),
		},
		{
			"cannot read bytes content",
			[]byte{0, 0, 0, 2, 1},
			nil,
			[]byte{1},
			fmt.Errorf("not enough bytes to read [bytes] content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadBytes(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteBytes(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{
			"empty bytes",
			[]byte{},
			make([]byte, LengthOfBytes([]byte{})),
			[]byte{0, 0, 0, 0},
			[]byte{},
			nil,
		},
		{
			"nil bytes",
			nil,
			make([]byte, LengthOfBytes([]byte{})),
			[]byte{0xff, 0xff, 0xff, 0xff},
			[]byte{},
			nil,
		},
		{
			"singleton bytes",
			[]byte{1},
			make([]byte, LengthOfBytes([]byte{1})),
			[]byte{0, 0, 0, 1, 1},
			[]byte{},
			nil,
		},
		{
			"simple bytes",
			[]byte{1, 2},
			make([]byte, LengthOfBytes([]byte{1, 2})),
			[]byte{0, 0, 0, 2, 1, 2},
			[]byte{},
			nil,
		},
		{
			"cannot write bytes length",
			[]byte{1},
			make([]byte, LengthOfInt-1),
			[]byte{0, 0, 0},
			[]byte{0, 0, 0},
			fmt.Errorf("cannot write [bytes] length: %w", errors.New("not enough capacity to write [int]")),
		},
		{
			"cannot write list element",
			[]byte{1, 2},
			make([]byte, LengthOfBytes([]byte{1, 2})-1),
			[]byte{0, 0, 0, 2, 0},
			[]byte{0},
			fmt.Errorf("not enough capacity to write [bytes] content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteBytes(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestReadShortBytes(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{"empty short bytes", []byte{0, 0}, []byte{}, []byte{}, nil},
		{"singleton short bytes", []byte{0, 1, 1}, []byte{1}, []byte{}, nil},
		{"simple short bytes", []byte{0, 2, 1, 2}, []byte{1, 2}, []byte{}, nil},
		{
			"cannot read short bytes length",
			[]byte{0},
			nil,
			[]byte{0},
			fmt.Errorf("cannot read [short bytes] length: %w", errors.New("not enough bytes to read [short]")),
		},
		{
			"cannot read short bytes content",
			[]byte{0, 2, 1},
			nil,
			[]byte{1},
			fmt.Errorf("not enough bytes to read [short bytes] content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadShortBytes(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteShortBytes(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{
			"empty short bytes",
			[]byte{},
			make([]byte, LengthOfShortBytes([]byte{})),
			[]byte{0, 0},
			[]byte{},
			nil,
		},
		// not officially allowed by the specs, but better safe than sorry
		{
			"nil short bytes",
			nil,
			make([]byte, LengthOfShortBytes(nil)),
			[]byte{0, 0},
			[]byte{},
			nil,
		},
		{
			"singleton short bytes",
			[]byte{1},
			make([]byte, LengthOfShortBytes([]byte{1})),
			[]byte{0, 1, 1},
			[]byte{},
			nil,
		},
		{
			"simple short bytes",
			[]byte{1, 2},
			make([]byte, LengthOfShortBytes([]byte{1, 2})),
			[]byte{0, 2, 1, 2},
			[]byte{},
			nil,
		},
		{
			"cannot write short bytes length",
			[]byte{1},
			make([]byte, LengthOfShort-1),
			[]byte{0},
			[]byte{0},
			fmt.Errorf("cannot write [short bytes] length: %w", errors.New("not enough capacity to write [short]")),
		},
		{
			"cannot write list element",
			[]byte{1, 2},
			make([]byte, LengthOfShortBytes([]byte{1, 2})-1),
			[]byte{0, 2, 0},
			[]byte{0},
			fmt.Errorf("not enough capacity to write [short bytes] content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteShortBytes(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

var uuid = cassandraprotocol.UUID{0xC0, 0xD1, 0xD2, 0x1E, 0xBB, 0x01, 0x41, 0x96, 0x86, 0xDB, 0xBC, 0x31, 0x7B, 0xC1, 0x79, 0x6A}
var uuidBytes = [16]byte{0xC0, 0xD1, 0xD2, 0x1E, 0xBB, 0x01, 0x41, 0x96, 0x86, 0xDB, 0xBC, 0x31, 0x7B, 0xC1, 0x79, 0x6A}

func TestReadUuid(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  *cassandraprotocol.UUID
		remaining []byte
		err       error
	}{
		{"simple UUID", uuidBytes[:], &uuid, []byte{}, nil},
		{"UUID with remaining", append(uuidBytes[:], 1, 2, 3, 4), &uuid, []byte{1, 2, 3, 4}, nil},
		{
			"cannot read UUID",
			uuidBytes[:15],
			nil,
			uuidBytes[:15],
			errors.New("not enough bytes to read [uuid] content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadUuid(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteUuid(t *testing.T) {
	tests := []struct {
		name      string
		input     *cassandraprotocol.UUID
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{
			"simple UUID",
			&uuid,
			make([]byte, LengthOfUuid),
			uuidBytes[:],
			[]byte{},
			nil,
		},
		{
			"UUID with remaining",
			&uuid,
			make([]byte, LengthOfUuid+1),
			append(uuidBytes[:], 0),
			[]byte{0},
			nil,
		},
		{
			"nil UUID",
			nil,
			make([]byte, LengthOfUuid),
			make([]byte, LengthOfUuid),
			make([]byte, LengthOfUuid),
			errors.New("cannot write nil as [uuid]"),
		},
		{
			"cannot write UUID content",
			&uuid,
			make([]byte, LengthOfUuid-1),
			make([]byte, LengthOfUuid-1),
			make([]byte, LengthOfUuid-1),
			errors.New("not enough capacity to write [uuid] content"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteUuid(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

var inet4 = cassandraprotocol.Inet{
	Addr: net.IPv4(192, 168, 1, 1),
	Port: 9042,
}
var inet4Bytes = []byte{
	4,              // length of IP
	192, 168, 1, 1, // IP
	0, 0, 0x23, 0x52, //port
}

var inet6 = cassandraprotocol.Inet{
	// 2001:0db8:85a3:0000:0000:8a2e:0370:7334
	Addr: net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34},
	Port: 9042,
}
var inet6Bytes = []byte{
	16,                                                                                             // length of IP
	0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34, // IP
	0, 0, 0x23, 0x52, //port
}

func TestReadInet(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  *cassandraprotocol.Inet
		remaining []byte
		err       error
	}{
		{"IPv4 INET", inet4Bytes[:], &inet4, []byte{}, nil},
		{"IPv6 INET", inet6Bytes[:], &inet6, []byte{}, nil},
		{"INET with remaining", append(inet4Bytes[:], 1, 2, 3, 4), &inet4, []byte{1, 2, 3, 4}, nil},
		{
			"cannot read INET length",
			[]byte{},
			nil,
			[]byte{},
			fmt.Errorf("cannot read [inet] length: %w", errors.New("not enough bytes to read [byte]")),
		},
		{
			"not enough bytes to read [inet] IPv4 content",
			[]byte{4, 192, 168, 1},
			nil,
			[]byte{192, 168, 1},
			errors.New("not enough bytes to read [inet] IPv4 content"),
		},
		{
			"not enough bytes to read [inet] IPv6 content",
			[]byte{16, 0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73},
			nil,
			[]byte{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73},
			errors.New("not enough bytes to read [inet] IPv6 content"),
		},
		{
			"cannot read [inet] port number",
			[]byte{4, 192, 168, 1, 1, 0, 0, 0},
			nil,
			[]byte{0, 0, 0},
			fmt.Errorf("cannot read [inet] port number: %w", errors.New("not enough bytes to read [int]")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadInet(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

var inet4Length, _ = LengthOfInet(&inet4)
var inet6Length, _ = LengthOfInet(&inet6)

func TestWriteInet(t *testing.T) {
	tests := []struct {
		name      string
		input     *cassandraprotocol.Inet
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{
			"IPv4 INET",
			&inet4,
			make([]byte, inet4Length),
			inet4Bytes,
			[]byte{},
			nil,
		},
		{
			"IPv6 INET",
			&inet6,
			make([]byte, inet6Length),
			inet6Bytes,
			[]byte{},
			nil,
		},
		{
			"INET with remaining",
			&inet4,
			make([]byte, inet4Length+1),
			append(inet4Bytes, 0),
			[]byte{0},
			nil,
		},
		{
			"cannot write nil INET",
			nil,
			[]byte{},
			[]byte{},
			[]byte{},
			errors.New("cannot write nil as [inet]"),
		},
		{
			"cannot write INET length",
			&inet4,
			[]byte{},
			[]byte{},
			[]byte{},
			fmt.Errorf("cannot write [inet] length: %w", errors.New("not enough capacity to write [byte]")),
		},
		{
			"not enough capacity to write [inet] IPv4 content",
			&inet4,
			make([]byte, inet4Length-LengthOfInt-1),
			[]byte{4, 0, 0, 0},
			[]byte{0, 0, 0},
			errors.New("not enough capacity to write [inet] IPv4 content"),
		},
		{
			"not enough capacity to write [inet] IPv6 content",
			&inet6,
			make([]byte, inet6Length-LengthOfInt-1),
			[]byte{16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			errors.New("not enough capacity to write [inet] IPv6 content"),
		},
		{
			"cannot write port number",
			&inet6,
			make([]byte, inet6Length-1),
			[]byte{
				16,                                                                                             // length of IP
				0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34, // IP
				0, 0, 0,
			},
			[]byte{0, 0, 0},
			fmt.Errorf("cannot write [inet] port number: %w", errors.New("not enough capacity to write [int]")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteInet(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestLengthOfInet(t *testing.T) {
	tests := []struct {
		name     string
		input    *cassandraprotocol.Inet
		expected int
		err      error
	}{
		{
			"IPv4 INET",
			&inet4,
			LengthOfByte + net.IPv4len + LengthOfInt,
			nil,
		},
		{
			"IPv6 INET",
			&inet6,
			LengthOfByte + net.IPv6len + LengthOfInt,
			nil,
		},
		{
			"nil INET",
			nil,
			-1,
			errors.New("cannot compute nil [inet] length"),
		},
		{
			"nil INET addr",
			&cassandraprotocol.Inet{},
			-1,
			errors.New("cannot compute nil [inet] length"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := LengthOfInet(tt.input)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestReadStringMap(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  map[string]string
		remaining []byte
		err       error
	}{
		{"empty string map", []byte{0, 0}, map[string]string{}, []byte{}, nil},
		{"map 1 key", []byte{
			0, 1, // map length
			0, 5, h, e, l, l, o, // key: hello
			0, 5, w, o, r, l, d, // value1: world
		}, map[string]string{"hello": "world"}, []byte{}, nil},
		{"map 2 keys", []byte{
			0, 2, // map length
			0, 5, h, e, l, l, o, // key1: hello
			0, 5, w, o, r, l, d, // value1: world
			0, 6, h, o, l, 0xc3, 0xa0, 0x21, // key2: holà!
			0, 5, m, u, n, d, o, // value2: mundo
		}, map[string]string{
			"hello": "world",
			"holà!": "mundo",
		}, []byte{}, nil},
		{
			"cannot read map length",
			[]byte{0},
			nil,
			[]byte{0},
			fmt.Errorf(
				"cannot read [string map] length: %w",
				errors.New("not enough bytes to read [short]"),
			),
		},
		{
			"cannot read key length",
			[]byte{0, 1, 0},
			nil,
			[]byte{0},
			fmt.Errorf(
				"cannot read [string map] key: %w",
				fmt.Errorf("cannot read [string] length: %w",
					errors.New("not enough bytes to read [short]")),
			),
		},
		{
			"cannot read key",
			[]byte{0, 1, 0, 2, 0},
			nil,
			[]byte{0},
			fmt.Errorf(
				"cannot read [string map] key: %w",
				errors.New("not enough bytes to read [string] content"),
			),
		},
		{
			"cannot read value length",
			[]byte{0, 1, 0, 1, k, 0},
			nil,
			[]byte{0},
			fmt.Errorf(
				"cannot read [string map] value: %w",
				fmt.Errorf("cannot read [string] length: %w",
					errors.New("not enough bytes to read [short]")),
			),
		},
		{
			"cannot read value",
			[]byte{0, 1, 0, 1, k, 0, 2, 0},
			nil,
			[]byte{0},
			fmt.Errorf(
				"cannot read [string map] value: %w",
				errors.New("not enough bytes to read [string] content"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadStringMap(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteStringMap(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]string
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{
			"empty string map",
			map[string]string{},
			make([]byte, LengthOfStringMap(map[string]string{})),
			[]byte{0, 0},
			[]byte{},
			nil,
		},
		// not officially allowed by the specs, but better safe than sorry
		{
			"nil string map",
			nil,
			make([]byte, LengthOfStringMap(nil)),
			[]byte{0, 0},
			[]byte{},
			nil,
		},
		{
			"map 1 key",
			map[string]string{"hello": "world"},
			make([]byte, LengthOfStringMap(map[string]string{"hello": "world"})),
			[]byte{
				0, 1, // map length
				0, 5, h, e, l, l, o, // key: hello
				0, 5, w, o, r, l, d, // value1: world
			},
			[]byte{},
			nil,
		},
		{
			"map 1 key with remaining",
			map[string]string{"hello": "world"},
			make([]byte, LengthOfStringMap(map[string]string{"hello": "world"})+1),
			[]byte{
				0, 1, // map length
				0, 5, h, e, l, l, o, // key: hello
				0, 5, w, o, r, l, d, // value1: world
				0, // extra
			},
			[]byte{0},
			nil,
		},
		// Cannot test maps with > 1 key since map entry iteration order is not deterministic :-(
		{
			"cannot write map length",
			map[string]string{},
			make([]byte, LengthOfShort-1),
			[]byte{0},
			[]byte{0},
			fmt.Errorf("cannot write [string map] length: %w",
				errors.New("not enough capacity to write [short]")),
		},
		{
			"cannot write key length",
			map[string]string{"hello": "world"},
			make([]byte, LengthOfShort+LengthOfShort-1),
			[]byte{0, 1, 0},
			[]byte{0},
			fmt.Errorf("cannot write [string map] key: %w",
				fmt.Errorf("cannot write [string] length: %w",
					errors.New("not enough capacity to write [short]"))),
		},
		{
			"cannot write key",
			map[string]string{"hello": "world"},
			make([]byte, LengthOfShort+LengthOfString("hello")-1),
			[]byte{0, 1, 0, 5, 0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			fmt.Errorf("cannot write [string map] key: %w",
				errors.New("not enough capacity to write [string] content")),
		},
		{
			"cannot write value length",
			map[string]string{"hello": "world"},
			make([]byte, LengthOfShort+LengthOfString("hello")+LengthOfShort-1),
			[]byte{0, 1, 0, 5, h, e, l, l, o, 0},
			[]byte{0},
			fmt.Errorf("cannot write [string map] value: %w",
				fmt.Errorf("cannot write [string] length: %w",
					errors.New("not enough capacity to write [short]"))),
		},
		{
			"cannot write value",
			map[string]string{"hello": "world"},
			make([]byte, LengthOfShort+LengthOfString("hello")+LengthOfString("world")-1),
			[]byte{0, 1, 0, 5, h, e, l, l, o, 0, 5, 0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			fmt.Errorf(
				"cannot write [string map] value: %w",
				errors.New("not enough capacity to write [string] content")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteStringMap(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestReadStringMultiMap(t *testing.T) {
	tests := []struct {
		name      string
		source    []byte
		expected  map[string][]string
		remaining []byte
		err       error
	}{
		{"empty string multimap", []byte{0, 0}, map[string][]string{}, []byte{}, nil},
		{"multimap 1 key 1 value", []byte{
			0, 1, // map length
			0, 5, h, e, l, l, o, // key: hello
			0, 1, // list length
			0, 5, w, o, r, l, d, // value1: world
		}, map[string][]string{"hello": {"world"}}, []byte{}, nil},
		{"multimap 1 key 2 values", []byte{
			0, 1, // map length
			0, 5, h, e, l, l, o, // key: hello
			0, 2, // list length
			0, 5, w, o, r, l, d, // value1: world
			0, 5, m, u, n, d, o, // value2: mundo
		}, map[string][]string{"hello": {"world", "mundo"}}, []byte{}, nil},
		{"multimap 2 keys 2 values", []byte{
			0, 2, // map length
			0, 5, h, e, l, l, o, // key1: hello
			0, 2, // list length
			0, 5, w, o, r, l, d, // value1: world
			0, 5, m, u, n, d, o, // value2: mundo
			0, 6, h, o, l, 0xc3, 0xa0, 0x21, // key2: holà!
			0, 2, // list length
			0, 5, w, o, r, l, d, // value1: world
			0, 5, m, u, n, d, o, // value2: mundo
		}, map[string][]string{
			"hello": {"world", "mundo"},
			"holà!": {"world", "mundo"},
		}, []byte{}, nil},
		{
			"cannot read map length",
			[]byte{0},
			nil,
			[]byte{0},
			fmt.Errorf(
				"cannot read [string multimap] length: %w",
				errors.New("not enough bytes to read [short]"),
			),
		},
		{
			"cannot read key length",
			[]byte{0, 1, 0},
			nil,
			[]byte{0},
			fmt.Errorf(
				"cannot read [string multimap] key: %w",
				fmt.Errorf("cannot read [string] length: %w",
					errors.New("not enough bytes to read [short]")),
			),
		},
		{
			"cannot read list length",
			[]byte{0, 1, 0, 1, k, 0},
			nil,
			[]byte{0},
			fmt.Errorf(
				"cannot read [string multimap] value: %w",
				fmt.Errorf("cannot read [string list] length: %w",
					errors.New("not enough bytes to read [short]")),
			),
		},
		{
			"cannot read element length",
			[]byte{0, 1, 0, 1, k, 0, 1, 0},
			nil,
			[]byte{0},
			fmt.Errorf(
				"cannot read [string multimap] value: %w",
				fmt.Errorf("cannot read [string list] element: %w",
					fmt.Errorf("cannot read [string] length: %w",
						errors.New("not enough bytes to read [short]"))),
			),
		},
		{
			"cannot read list",
			[]byte{0, 1, 0, 1, k, 0, 1, 0, 5, h, e, l, l},
			nil,
			[]byte{h, e, l, l},
			fmt.Errorf("cannot read [string multimap] value: %w",
				fmt.Errorf("cannot read [string list] element: %w",
					errors.New("not enough bytes to read [string] content"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, remaining, err := ReadStringMultiMap(tt.source)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestWriteStringMultiMap(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string][]string
		dest      []byte
		expected  []byte
		remaining []byte
		err       error
	}{
		{
			"empty string multimap",
			map[string][]string{},
			make([]byte, LengthOfStringMultiMap(map[string][]string{})),
			[]byte{0, 0},
			[]byte{},
			nil,
		},
		// not officially allowed by the specs, but better safe than sorry
		{
			"nil string multimap",
			nil,
			make([]byte, LengthOfStringMultiMap(nil)),
			[]byte{0, 0},
			[]byte{},
			nil,
		},
		{
			"multimap 1 key 1 value",
			map[string][]string{"hello": {"world"}},
			make([]byte, LengthOfStringMultiMap(map[string][]string{"hello": {"world"}})),
			[]byte{
				0, 1, // map length
				0, 5, h, e, l, l, o, // key: hello
				0, 1, // list length
				0, 5, w, o, r, l, d, // value1: world
			},
			[]byte{},
			nil,
		},
		{
			"multimap 1 key 1 value with remaining",
			map[string][]string{"hello": {"world"}},
			make([]byte, LengthOfStringMultiMap(map[string][]string{"hello": {"world"}})+1),
			[]byte{
				0, 1, // map length
				0, 5, h, e, l, l, o, // key: hello
				0, 1, // list length
				0, 5, w, o, r, l, d, // value1: world
				0, // extra
			},
			[]byte{0},
			nil,
		},
		{
			"multimap 1 key 2 values",
			map[string][]string{"hello": {"world", "mundo"}},
			make([]byte, LengthOfStringMultiMap(map[string][]string{"hello": {"world", "mundo"}})),
			[]byte{
				0, 1, // map length
				0, 5, h, e, l, l, o, // key: hello
				0, 2, // list length
				0, 5, w, o, r, l, d, // value1: world
				0, 5, m, u, n, d, o, // value2: mundo
			},
			[]byte{},
			nil,
		},
		// Cannot test maps with > 1 key since map entry iteration order is not deterministic :-(
		{
			"cannot write map length",
			map[string][]string{},
			make([]byte, LengthOfShort-1),
			[]byte{0},
			[]byte{0},
			fmt.Errorf("cannot write [string multimap] length: %w",
				errors.New("not enough capacity to write [short]")),
		},
		{
			"cannot write key length",
			map[string][]string{"hello": {"world"}},
			make([]byte, LengthOfShort+LengthOfShort-1),
			[]byte{0, 1, 0},
			[]byte{0},
			fmt.Errorf("cannot write [string multimap] key: %w",
				fmt.Errorf("cannot write [string] length: %w",
					errors.New("not enough capacity to write [short]"))),
		},
		{
			"cannot write key",
			map[string][]string{"hello": {"world"}},
			make([]byte, LengthOfShort+LengthOfString("hello")-1),
			[]byte{0, 1, 0, 5, 0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			fmt.Errorf("cannot write [string multimap] key: %w",
				errors.New("not enough capacity to write [string] content")),
		},
		{
			"cannot write list length",
			map[string][]string{"hello": {"world"}},
			make([]byte, LengthOfShort+LengthOfString("hello")+LengthOfShort-1),
			[]byte{0, 1, 0, 5, h, e, l, l, o, 0},
			[]byte{0},
			fmt.Errorf("cannot write [string multimap] value: %w",
				fmt.Errorf("cannot write [string list] length: %w",
					errors.New("not enough capacity to write [short]"))),
		},
		{
			"cannot write element length",
			map[string][]string{"hello": {"world"}},
			make([]byte, LengthOfShort+LengthOfString("hello")+LengthOfShort+LengthOfShort-1),
			[]byte{0, 1, 0, 5, h, e, l, l, o, 0, 1, 0},
			[]byte{0},
			fmt.Errorf("cannot write [string multimap] value: %w",
				fmt.Errorf("cannot write [string list] element: %w",
					fmt.Errorf("cannot write [string] length: %w",
						errors.New("not enough capacity to write [short]")))),
		},
		{
			"cannot write list element",
			map[string][]string{"hello": {"world"}},
			make([]byte, LengthOfShort+LengthOfString("hello")+LengthOfShort+LengthOfString("world")-1),
			[]byte{0, 1, 0, 5, h, e, l, l, o, 0, 1, 0, 5, 0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			fmt.Errorf(
				"cannot write [string multimap] value: %w",
				fmt.Errorf("cannot write [string list] element: %w",
					errors.New("not enough capacity to write [string] content")),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, err := WriteStringMultiMap(tt.input, tt.dest)
			assert.Equal(t, tt.expected, tt.dest)
			assert.Equal(t, tt.remaining, remaining)
			assert.Equal(t, tt.err, err)
		})
	}
}

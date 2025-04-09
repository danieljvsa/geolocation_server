package main

import (
	"encoding/hex"
	"testing"
)

func TestDecodeIos(t *testing.T) {
	tests := []struct {
		name       string
		data       []byte
		dataLength int64
		startByte  int64
		wantErr    bool
		want       ResponseDecode
	}{
		{
			name: "Valid data with one-byte IO",
			data: func() []byte {
				hexStr := "0101024A" // 1 IO, 1 one-byte IO, ID=2, Value=4A
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			dataLength: 4,
			startByte:  0,
			wantErr:    false,
			want:       ResponseDecode{IOs: []IOData{{IO: 2, Value: "4a"}}, NumberOfIOs: 1},
		},
		{
			name:       "No IOs",
			data:       []byte{0x00}, // Zero IO count
			dataLength: 1,
			startByte:  0,
			wantErr:    false,
			want:       ResponseDecode{IOs: []IOData{}, NumberOfIOs: 0},
		},
		{
			name: "Invalid data length (out of bounds)",
			data: func() []byte {
				hexStr := "01" // Incomplete data
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			dataLength: 1,
			startByte:  0,
			wantErr:    false,
			want:       ResponseDecode{IOs: []IOData{}, NumberOfIOs: 0},
		},
		{
			name: "Incorrect byte parsing",
			data: func() []byte {
				hexStr := "01010502" // 1 IO, 5 one-byte IOs, but missing bytes
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			dataLength: 4,
			startByte:  0,
			wantErr:    false,
			want:       ResponseDecode{IOs: []IOData{{IO: 5, Value: "02"}}, NumberOfIOs: 1},
		},
		{
			name: "Test with all types",
			data: func() []byte {
				hexStr := "05021503010101425E0F01F10000601A014E0000000000000000" // 1 IO, 5 one-byte IOs, but missing bytes
				data, _ := hex.DecodeString(hexStr)
				return data
			}(),
			dataLength: 26,
			startByte:  0,
			wantErr:    false,
			want:       ResponseDecode{IOs: []IOData{{IO: 21, Value: "03"}, {IO: 1, Value: "01"}, {IO: 66, Value: "5e0f"}, {IO: 241, Value: "0000601a"}, {IO: 78, Value: "0000000000000000"}}, NumberOfIOs: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeIos(tt.data, tt.dataLength, tt.startByte)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeIos() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !equalIOData(got.IOs, tt.want.IOs) {
				t.Errorf("decodeIos() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare slices
func equalIOData(a, b []IOData) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].IO != b[i].IO || a[i].Value != b[i].Value {
			return false
		}
	}
	return true
}

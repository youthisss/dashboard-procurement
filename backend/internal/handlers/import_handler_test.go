package handlers

import "testing"

func TestSanitizeUploadedExcelFilename(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid xlsx", input: "pricing.xlsx", want: "pricing.xlsx"},
		{name: "valid uppercase extension", input: "pricing.XLSX", want: "pricing.XLSX"},
		{name: "reject parent path", input: "../pricing.xlsx", wantErr: true},
		{name: "reject windows path", input: `..\pricing.xlsx`, wantErr: true},
		{name: "reject nested path", input: "nested/pricing.xlsx", wantErr: true},
		{name: "reject non excel", input: "pricing.csv", wantErr: true},
		{name: "reject empty", input: " ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitizeUploadedExcelFilename(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("sanitizeUploadedExcelFilename(%q) error = nil, want error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("sanitizeUploadedExcelFilename(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("sanitizeUploadedExcelFilename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

package session

import "testing"

func TestSession_CloseSession(t *testing.T) {
	tests := []struct {
		name    string
		s       *Session
		wantErr bool
	}{
		{
			"It should close the session",
			&Session{
				"id123",
				make(chan bool),
				&eventEmitterMock{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.CloseSession(); (err != nil) != tt.wantErr {
				t.Errorf("Session.CloseSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

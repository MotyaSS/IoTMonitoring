package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
)

type producerMock struct {
	lastMsg any
	err     error
}

func (m *producerMock) Produce(_ context.Context, msg any) error {
	m.lastMsg = msg
	return m.err
}

func TestSendTelemetry_ForwardsMessageToProducer(t *testing.T) {
	m := &producerMock{}
	svc := NewScrapperService(m, slog.New(slog.NewTextHandler(io.Discard, nil)))
	input := &pb.Telemetry{SenderId: 7}

	if err := svc.SendTelemetry(context.Background(), input); err != nil {
		t.Fatalf("SendTelemetry returned error: %v", err)
	}

	got, ok := m.lastMsg.(*pb.Telemetry)
	if !ok {
		t.Fatalf("expected *Telemetry passed to producer, got %T", m.lastMsg)
	}
	if got.GetSenderId() != 7 {
		t.Fatalf("expected sender_id=7, got %d", got.GetSenderId())
	}
}

func TestSendTelemetry_PropagatesProducerError(t *testing.T) {
	wantErr := errors.New("boom")
	m := &producerMock{err: wantErr}
	svc := NewScrapperService(m, slog.New(slog.NewTextHandler(io.Discard, nil)))

	err := svc.SendTelemetry(context.Background(), &pb.Telemetry{SenderId: 1})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected producer error, got %v", err)
	}
}

func TestAuthenticate_ReturnsNil(t *testing.T) {
	svc := NewScrapperService(&producerMock{}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err := svc.Authenticate("token"); err != nil {
		t.Fatalf("expected nil auth error, got %v", err)
	}
}

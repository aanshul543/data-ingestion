package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anshul543/data-ingestion/internal/transformer"
	"github.com/anshul543/data-ingestion/internal/uploader"
)

// MockUploader implements uploader.S3Uploader
type MockUploader struct {
	ShouldFail bool
	Called     bool
}

func (m *MockUploader) Upload(posts []transformer.TransformedPost) error {
	m.Called = true
	if m.ShouldFail {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// MockAPI sets up API with mock uploader (skipping real S3 client)
/*func MockAPI(mock MockUploader) *API {
	return &API{
		uploader: &mock,
	}
}*/

func MockAPI(u uploader.S3Uploader) *API {
	return &API{
		uploader: u,
	}
}

func TestIngestData_Success(t *testing.T) {
	mock := &MockUploader{}
	api := MockAPI(mock)

	req := httptest.NewRequest(http.MethodPost, "/ingest", nil)
	w := httptest.NewRecorder()

	api.IngestData(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Ingestion completed") {
		t.Errorf("Unexpected body: %s", string(body))
	}
	if !mock.Called {
		t.Errorf("Expected uploader to be called")
	}
}

func TestIngestData_Failure(t *testing.T) {
	mock := &MockUploader{ShouldFail: true}
	api := MockAPI(mock)

	req := httptest.NewRequest(http.MethodPost, "/ingest", nil)
	w := httptest.NewRecorder()

	api.IngestData(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected 500 Internal Server Error, got %d", resp.StatusCode)
	}
}

func TestIngestData_WrongMethod(t *testing.T) {
	mock := &MockUploader{}
	api := MockAPI(mock)

	req := httptest.NewRequest(http.MethodGet, "/ingest", nil)
	w := httptest.NewRecorder()

	api.IngestData(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 Method Not Allowed, got %d", resp.StatusCode)
	}
}

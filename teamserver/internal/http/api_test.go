package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
	thttp "github.com/sentientbottleofwine/osmium/teamserver/internal/http"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

const port = 2137
const addr = "http://localhost:2137"

func startServer(t *testing.T) (*thttp.Server, <-chan error) {
	dbHandle, err := database.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to get setup db: %v", err)
	}

	server := thttp.NewServer(port, dbHandle)

	serverErrCh := make(chan error)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		serverErrCh <- nil
	}()
	return server, serverErrCh
}

func TestAddingAgents(t *testing.T) {
	server, serverErrCh := startServer(t)

	agentCreationResponse, err := http.Post(addr+"/api/register", "application/json", bytes.NewBufferString(""))
	if err != nil {
		t.Fatal(err)
	}
	var agent api.RegisterResponse
	err = json.NewDecoder(agentCreationResponse.Body).Decode(&agent)
	if err != nil {
		t.Fatal(err)
	}

	// Check if public key is valid
	_, err = tools.PemToPubRsa(agent.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	server.Close()

	err, ok := <-serverErrCh
	if !ok {
		t.Fatal("Channel closed somehow")
	}

	if err != nil {
		t.Fatal(err)
	}
}

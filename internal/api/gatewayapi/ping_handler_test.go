package gatewayapi

import (
	"encoding/json"
	"net/http"
	"testing"

	"bou.ke/monkey"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/stretchr/testify/assert"
)

type mockResponseWriter struct {
	payload []byte
	header  map[string][]string
}

func (m *mockResponseWriter) WriteHeader(_ int) {
}

func (m *mockResponseWriter) Header() http.Header {
	return m.header
}

func (m *mockResponseWriter) EncodeJson(v interface{}) ([]byte, error) {
	m.payload, _ = json.Marshal(v)

	return m.payload, nil
}

func (m *mockResponseWriter) WriteJson(v interface{}) error {
	b, _ := m.EncodeJson(v)
	_, _ = m.Write(b)

	return nil
}
func (w *mockResponseWriter) Write(i []byte) (int, error) {
	return len(i), nil
}

// TestHandleGatewayPingRequest success test
func TestHandleGatewayPingRequest(t *testing.T) {
	mockNodeID, _ := nodeid.NewNodeIDFromHexString("42")
	mockNonce := int64(42)
	mockTTL := int64(43)
	mockPrivateKey, _ := fcrcrypto.GenerateRetrievalV1KeyPair()
	mockKeyVersion := fcrcrypto.InitialKeyVersion()

	patchGetSingleInstance := monkey.Patch(core.GetSingleInstance, func(_ ...*settings.AppSettings) *core.Core {
		return &core.Core{
			GatewayPrivateKey:        mockPrivateKey,
			GatewayPrivateKeyVersion: mockKeyVersion,
		}
	})
	defer patchGetSingleInstance.Unpatch()

	request, _ := fcrmessages.EncodeGatewayPingRequest(mockNodeID, mockNonce, mockTTL)

	mockRW := &mockResponseWriter{}

	// func to test
	HandleGatewayPingRequest(mockRW, request)

	assert.NotEmpty(t, mockRW.payload)

	response := &fcrmessages.FCRMessage{}
	assert.Nil(t, json.Unmarshal(mockRW.payload, response))

	nonce, isAlive, err := fcrmessages.DecodeGatewayPingResponse(response)
	assert.Nil(t, err)
	assert.True(t, isAlive)
	assert.Equal(t, mockNonce, nonce)
}

package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Vaibhav2154/ShadowNet/internal/shared/proto"
)

// Client is a control plane API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new control plane client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Register registers this node with the control plane
func (c *Client) Register(info *proto.PeerInfo) error {
	req := proto.RegisterRequest{
		ID:           info.ID,
		WGPublicKey:  info.WGPublicKey,
		EndpointIP:   info.EndpointIP,
		EndpointPort: info.EndpointPort,
	}

	var resp proto.RegisterResponse
	if err := c.post("/register", req, &resp); err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("registration failed: %s", resp.Message)
	}

	return nil
}

// GetPeers retrieves the list of active peers
func (c *Client) GetPeers(excludeID string) ([]*proto.PeerInfo, error) {
	url := c.baseURL + "/peers"
	if excludeID != "" {
		url += "?exclude=" + excludeID
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get peers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get peers failed: %s - %s", resp.Status, string(body))
	}

	var peersResp proto.PeersResponse
	if err := json.NewDecoder(resp.Body).Decode(&peersResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to pointer slice
	var peers []*proto.PeerInfo
	for i := range peersResp.Peers {
		peers = append(peers, &peersResp.Peers[i])
	}

	return peers, nil
}

// SendHeartbeat sends a heartbeat to the control plane
func (c *Client) SendHeartbeat(id string) error {
	req := proto.HeartbeatRequest{
		ID: id,
	}

	var resp proto.HeartbeatResponse
	if err := c.post("/heartbeat", req, &resp); err != nil {
		return fmt.Errorf("heartbeat failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("heartbeat failed: %s", resp.Message)
	}

	return nil
}

// GetMetrics retrieves control plane metrics
func (c *Client) GetMetrics() (*proto.MetricsResponse, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/metrics")
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get metrics failed: %s", resp.Status)
	}

	var metrics proto.MetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &metrics, nil
}

// post sends a POST request and decodes the response
func (c *Client) post(path string, request, response interface{}) error {
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+path,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s - %s", resp.Status, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

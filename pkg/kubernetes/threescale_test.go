package kubernetes

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/ghodss/yaml"

	"github.com/3scale/3scale-istio-adapter/config"
	"istio.io/api/policy/v1beta1"
)

func TestNewThreescaleHandlerSpec(t *testing.T) {
	inputs := []struct {
		name      string
		token     string
		svcID     string
		url       string
		expectErr bool
		expectRes *HandlerSpec
	}{
		{
			name:      "Test expect error -  no url",
			token:     "1234",
			svcID:     "12345",
			expectErr: true,
		},
		{
			name:      "Test expect error -  no token",
			svcID:     "12345",
			url:       "https://test.com",
			expectErr: true,
		},
		{
			name:      "Test expect error -  invalid url",
			token:     "12345",
			url:       "https://<t est.com",
			expectErr: true,
		},
		{
			name:      "Test happy path",
			token:     "12345",
			svcID:     "54321",
			url:       "https://test.com",
			expectErr: false,
			expectRes: &HandlerSpec{
				Adapter: defaultThreescaleAdapterName,
				Params: config.Params{
					ServiceId:   "54321",
					SystemUrl:   "https://test.com",
					AccessToken: "12345",
				},
				Connection: v1beta1.Connection{
					Address: defaultThreescaleAdapterListenAddress + ":" + strconv.Itoa(defaultThreescaleAdapterListenPort),
				},
			},
		},
		{
			name:      "Test happy path with port",
			token:     "12345",
			svcID:     "54321",
			url:       "https://test.com:8080",
			expectErr: false,
			expectRes: &HandlerSpec{
				Adapter: defaultThreescaleAdapterName,
				Params: config.Params{
					ServiceId:   "54321",
					SystemUrl:   "https://test.com:8080",
					AccessToken: "12345",
				},
				Connection: v1beta1.Connection{
					Address: defaultThreescaleAdapterListenAddress + ":" + strconv.Itoa(defaultThreescaleAdapterListenPort),
				},
			},
		},
	}
	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			result, err := NewThreescaleHandlerSpec(input.token, input.url, input.svcID)
			if input.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if !reflect.DeepEqual(result, input.expectRes) {
				t.Errorf("returned handler does not match expected handler")
			}
		})
	}
}

func TestNewApiKeyInstance(t *testing.T) {
	expect := `params:
  action:
    method: request.method | "get"
    path: request.url_path
    service: destination.labels["service-mesh.3scale.net/service-id"] | ""
  subject:
    user: request.query_params["user_key"] | request.headers["x-user-key"] | ""
template: threescale-authorization`
	instance := NewApiKeyInstance(getDefaultApiKeyAttributeString())
	b, err := json.MarshalIndent(instance, "", "  ")
	if err != nil {
		t.Errorf("unexpected error when converting to JSON")
	}

	b, err = yaml.JSONToYAML(b)
	if err != nil {
		t.Errorf("unexpected error when converting JSON to YAML")
	}

	if strings.TrimSpace(string(b)) != expect {
		t.Errorf("unexpected YAML returned.\nWanted:\n%s\nGot:\n%s", expect, string(b))
	}
}

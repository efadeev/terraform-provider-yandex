package yandex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/apploadbalancer/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

func Test_flattenALBRateLimit(t *testing.T) {
	t.Parallel()

	testsTable := []struct {
		name           string
		rateLimit      *apploadbalancer.RateLimit
		expectedResult []map[string]interface{}
	}{
		{
			name: "nil rate limit",
		},
		{
			name:           "empty rate limit",
			rateLimit:      &apploadbalancer.RateLimit{},
			expectedResult: []map[string]interface{}{{}},
		},
		{
			name: "empty all requests limit",
			rateLimit: &apploadbalancer.RateLimit{
				AllRequests: &apploadbalancer.RateLimit_Limit{},
			},
			expectedResult: []map[string]interface{}{
				{
					allRequestsSchemaKey: []map[string]interface{}{{}},
				},
			},
		},
		{
			name: "all requests rps",
			rateLimit: &apploadbalancer.RateLimit{
				AllRequests: &apploadbalancer.RateLimit_Limit{
					Rate: &apploadbalancer.RateLimit_Limit_PerSecond{
						PerSecond: 15,
					},
				},
			},
			expectedResult: []map[string]interface{}{
				{
					allRequestsSchemaKey: []map[string]interface{}{
						{
							perSecondSchemaKey: 15,
						},
					},
				},
			},
		},
		{
			name: "all requests rpm",
			rateLimit: &apploadbalancer.RateLimit{
				AllRequests: &apploadbalancer.RateLimit_Limit{
					Rate: &apploadbalancer.RateLimit_Limit_PerMinute{
						PerMinute: 15,
					},
				},
			},
			expectedResult: []map[string]interface{}{
				{
					allRequestsSchemaKey: []map[string]interface{}{
						{
							perMinuteSchemaKey: 15,
						},
					},
				},
			},
		},
		{
			name: "empty requests per ip limit",
			rateLimit: &apploadbalancer.RateLimit{
				RequestsPerIp: &apploadbalancer.RateLimit_Limit{},
			},
			expectedResult: []map[string]interface{}{
				{
					requestsPerIPSchemaKey: []map[string]interface{}{{}},
				},
			},
		},
		{
			name: "requests per ip rps",
			rateLimit: &apploadbalancer.RateLimit{
				RequestsPerIp: &apploadbalancer.RateLimit_Limit{
					Rate: &apploadbalancer.RateLimit_Limit_PerSecond{
						PerSecond: 15,
					},
				},
			},
			expectedResult: []map[string]interface{}{
				{
					requestsPerIPSchemaKey: []map[string]interface{}{
						{
							perSecondSchemaKey: 15,
						},
					},
				},
			},
		},
		{
			name: "all requests rpm",
			rateLimit: &apploadbalancer.RateLimit{
				RequestsPerIp: &apploadbalancer.RateLimit_Limit{
					Rate: &apploadbalancer.RateLimit_Limit_PerMinute{
						PerMinute: 15,
					},
				},
			},
			expectedResult: []map[string]interface{}{
				{
					requestsPerIPSchemaKey: []map[string]interface{}{
						{
							perMinuteSchemaKey: 15,
						},
					},
				},
			},
		},
		{
			name: "all requests and requests per ip limits",
			rateLimit: &apploadbalancer.RateLimit{
				AllRequests: &apploadbalancer.RateLimit_Limit{
					Rate: &apploadbalancer.RateLimit_Limit_PerSecond{
						PerSecond: 10,
					},
				},
				RequestsPerIp: &apploadbalancer.RateLimit_Limit{
					Rate: &apploadbalancer.RateLimit_Limit_PerMinute{
						PerMinute: 15,
					},
				},
			},
			expectedResult: []map[string]interface{}{
				{
					allRequestsSchemaKey: []map[string]interface{}{
						{
							perSecondSchemaKey: 10,
						},
					},
					requestsPerIPSchemaKey: []map[string]interface{}{
						{
							perMinuteSchemaKey: 15,
						},
					},
				},
			},
		},
	}

	for _, testCase := range testsTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actualResult := flattenALBRateLimit(testCase.rateLimit)

			assert.Equal(t, testCase.expectedResult, actualResult)
		})
	}
}

func Test_flattenALBRegexMatchAndSubstitute(t *testing.T) {
	t.Parallel()

	testsTable := []struct {
		name           string
		regexRewrite   *apploadbalancer.RegexMatchAndSubstitute
		expectedResult []map[string]interface{}
	}{
		{
			name: "nil regex rewrite",
		},
		{
			name:           "empty regex rewrite",
			regexRewrite:   &apploadbalancer.RegexMatchAndSubstitute{},
			expectedResult: []map[string]interface{}{{}},
		},
		{
			name: "regex rewrite",
			regexRewrite: &apploadbalancer.RegexMatchAndSubstitute{
				Regex:      "regex",
				Substitute: "substitute",
			},
			expectedResult: []map[string]interface{}{
				{
					regexSchemaKey:      "regex",
					substituteSchemaKey: "substitute",
				},
			},
		},
		{
			name: "regex rewrite: empty regex",
			regexRewrite: &apploadbalancer.RegexMatchAndSubstitute{
				Substitute: "substitute",
			},
			expectedResult: []map[string]interface{}{
				{
					substituteSchemaKey: "substitute",
				},
			},
		},
		{
			name: "regex rewrite: empty substitute",
			regexRewrite: &apploadbalancer.RegexMatchAndSubstitute{
				Regex: "regex",
			},
			expectedResult: []map[string]interface{}{
				{
					regexSchemaKey: "regex",
				},
			},
		},
	}

	for _, testCase := range testsTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actualResult := flattenALBRegexMatchAndSubstitute(testCase.regexRewrite)

			assert.Equal(t, testCase.expectedResult, actualResult)
		})
	}
}

func Test_flattenALBStreamBackends(t *testing.T) {
	t.Parallel()

	testsTable := []struct {
		name           string
		backendGroup   *apploadbalancer.BackendGroup
		expectedResult []interface{}
		expectErr      bool
	}{
		{
			name: "stream backend: keep_connections_on_host_health_failure set to false",
			backendGroup: &apploadbalancer.BackendGroup{
				Name:        "backend-group",
				Description: "some-backend-group",
				Backend: &apploadbalancer.BackendGroup_Stream{
					Stream: &apploadbalancer.StreamBackendGroup{
						Backends: []*apploadbalancer.StreamBackend{
							{
								Name:                               "stream-backend",
								KeepConnectionsOnHostHealthFailure: false,
							},
						},
					},
				},
			},
			expectedResult: []interface{}{
				map[string]interface{}{
					"name":                  "stream-backend",
					"port":                  0,
					"weight":                1,
					"tls":                   []map[string]interface{}{},
					"healthcheck":           []interface{}(nil),
					"load_balancing_config": []map[string]interface{}{},
					"enable_proxy_protocol": false,
					keepConnectionsOnHostHealthFailureSchemaKey: false,
				},
			},
		},
		{
			name: "stream backend: keep_connections_on_host_health_failure set to true",
			backendGroup: &apploadbalancer.BackendGroup{
				Name:        "backend-group",
				Description: "some-backend-group",
				Backend: &apploadbalancer.BackendGroup_Stream{
					Stream: &apploadbalancer.StreamBackendGroup{
						Backends: []*apploadbalancer.StreamBackend{
							{
								Name:                               "stream-backend",
								KeepConnectionsOnHostHealthFailure: true,
							},
						},
					},
				},
			},
			expectedResult: []interface{}{
				map[string]interface{}{
					"name":                  "stream-backend",
					"port":                  0,
					"weight":                1,
					"tls":                   []map[string]interface{}{},
					"healthcheck":           []interface{}(nil),
					"load_balancing_config": []map[string]interface{}{},
					"enable_proxy_protocol": false,
					keepConnectionsOnHostHealthFailureSchemaKey: true,
				},
			},
		},
	}

	for _, testCase := range testsTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actualResult, err := flattenALBStreamBackends(testCase.backendGroup)

			if testCase.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedResult, actualResult)
			}
		})
	}
}

func Test_flattenALBHealthChecks(t *testing.T) {
	t.Parallel()

	testsTable := []struct {
		name           string
		healthchecks   []*apploadbalancer.HealthCheck
		expectedResult []interface{}
		expectErr      bool
	}{
		{
			name: "http backend: use nil expected statuses slice",
			healthchecks: []*apploadbalancer.HealthCheck{
				{
					Timeout:  durationpb.New(time.Second),
					Interval: durationpb.New(time.Second),
					Healthcheck: &apploadbalancer.HealthCheck_Http{
						Http: &apploadbalancer.HealthCheck_HttpHealthCheck{
							Path: "/",
						},
					},
				},
			},
			expectedResult: []interface{}{
				map[string]interface{}{
					"timeout":                 formatDuration(durationpb.New(time.Second)),
					"interval":                formatDuration(durationpb.New(time.Second)),
					"interval_jitter_percent": float64(0),
					"healthy_threshold":       int64(0),
					"unhealthy_threshold":     int64(0),
					"healthcheck_port":        0,
					"http_healthcheck": []map[string]interface{}{
						{
							"host":                    "",
							"path":                    "/",
							"http2":                   false,
							expectedStatusesSchemaKey: []int64(nil),
						},
					},
				},
			},
		},
		{
			name: "http backend: use empty expected statuses slice",
			healthchecks: []*apploadbalancer.HealthCheck{
				{
					Timeout:  durationpb.New(time.Second),
					Interval: durationpb.New(time.Second),
					Healthcheck: &apploadbalancer.HealthCheck_Http{
						Http: &apploadbalancer.HealthCheck_HttpHealthCheck{
							Path:             "/",
							ExpectedStatuses: []int64{},
						},
					},
				},
			},
			expectedResult: []interface{}{
				map[string]interface{}{
					"timeout":                 formatDuration(durationpb.New(time.Second)),
					"interval":                formatDuration(durationpb.New(time.Second)),
					"interval_jitter_percent": float64(0),
					"healthy_threshold":       int64(0),
					"unhealthy_threshold":     int64(0),
					"healthcheck_port":        0,
					"http_healthcheck": []map[string]interface{}{
						{
							"host":                    "",
							"path":                    "/",
							"http2":                   false,
							expectedStatusesSchemaKey: []int64{},
						},
					},
				},
			},
		},
		{
			name: "http backend: use expected statuses",
			healthchecks: []*apploadbalancer.HealthCheck{
				{
					Timeout:  durationpb.New(time.Second),
					Interval: durationpb.New(time.Second),
					Healthcheck: &apploadbalancer.HealthCheck_Http{
						Http: &apploadbalancer.HealthCheck_HttpHealthCheck{
							Path:             "/",
							ExpectedStatuses: []int64{200, 201, 300, 302},
						},
					},
				},
			},
			expectedResult: []interface{}{
				map[string]interface{}{
					"timeout":                 formatDuration(durationpb.New(time.Second)),
					"interval":                formatDuration(durationpb.New(time.Second)),
					"interval_jitter_percent": float64(0),
					"healthy_threshold":       int64(0),
					"unhealthy_threshold":     int64(0),
					"healthcheck_port":        0,
					"http_healthcheck": []map[string]interface{}{
						{
							"host":                    "",
							"path":                    "/",
							"http2":                   false,
							expectedStatusesSchemaKey: []int64{200, 201, 300, 302},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testsTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actualResult := flattenALBHealthChecks(testCase.healthchecks)

			assert.Equal(t, testCase.expectedResult, actualResult)
		})
	}
}

func Test_flattenALBAutoscalePolicy(t *testing.T) {
	t.Parallel()

	testsTable := []struct {
		name           string
		autoscale      *apploadbalancer.AutoScalePolicy
		expectedResult []map[string]interface{}
		expectErr      bool
	}{
		{
			name:           "nil value",
			autoscale:      nil,
			expectedResult: nil,
		},
		{
			name: "some sane values",
			autoscale: &apploadbalancer.AutoScalePolicy{
				MinZoneSize: 10,
				MaxSize:     3,
			},
			expectedResult: []map[string]interface{}{{
				"min_zone_size": interface{}(10),
				"max_size":      interface{}(3),
			}},
		},
		{
			name: "only min_zone_size specified",
			autoscale: &apploadbalancer.AutoScalePolicy{
				MinZoneSize: 10,
			},
			expectedResult: []map[string]interface{}{{
				"min_zone_size": interface{}(10),
			}},
		},
		{
			name: "only max_size specified",
			autoscale: &apploadbalancer.AutoScalePolicy{
				MaxSize: 10,
			},
			expectedResult: []map[string]interface{}{{
				"max_size": interface{}(10),
			}},
		},
	}

	for _, testCase := range testsTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actualResult, err := flattenALBAutoscalePolicy(&apploadbalancer.LoadBalancer{
				Id:              "1",
				AutoScalePolicy: testCase.autoscale,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResult, actualResult)
		})
	}
}

func Test_flattenALBRoutes(t *testing.T) {
	t.Parallel()

	testsTable := []struct {
		name           string
		routes         []*apploadbalancer.Route
		expectedResult []map[string]any
		expectErr      bool
	}{
		{
			name: "http route: disable security profile: true",
			routes: []*apploadbalancer.Route{
				{
					Name: "my_little_route",
					Route: &apploadbalancer.Route_Http{
						Http: &apploadbalancer.HttpRoute{
							Match: &apploadbalancer.HttpRouteMatch{
								Path: &apploadbalancer.StringMatch{
									Match: &apploadbalancer.StringMatch_ExactMatch{
										ExactMatch: "/",
									},
								},
							},
							Action: &apploadbalancer.HttpRoute_DirectResponse{
								DirectResponse: &apploadbalancer.DirectResponseAction{
									Status: 200,
									Body: &apploadbalancer.Payload{
										Payload: &apploadbalancer.Payload_Text{
											Text: "hello world",
										},
									},
								},
							},
						},
					},
					DisableSecurityProfile: true,
				},
			},
			expectedResult: []map[string]any{
				{
					"name": "my_little_route",
					"http_route": []map[string]any{
						{
							"http_match": []map[string]any{
								{
									"http_method": []string(nil),
									"path": []map[string]any{
										{
											"exact": "/",
										},
									},
								},
							},
							"direct_response_action": []map[string]any{
								{
									"body":   "hello world",
									"status": 200,
								},
							},
						},
					},
					"route_options":            []map[string]any(nil),
					"disable_security_profile": true,
				},
			},
		},
		{
			name: "http route: disable security profile: false",
			routes: []*apploadbalancer.Route{
				{
					Name: "my_little_route",
					Route: &apploadbalancer.Route_Http{
						Http: &apploadbalancer.HttpRoute{
							Match: &apploadbalancer.HttpRouteMatch{
								Path: &apploadbalancer.StringMatch{
									Match: &apploadbalancer.StringMatch_ExactMatch{
										ExactMatch: "/",
									},
								},
							},
							Action: &apploadbalancer.HttpRoute_DirectResponse{
								DirectResponse: &apploadbalancer.DirectResponseAction{
									Status: 200,
									Body: &apploadbalancer.Payload{
										Payload: &apploadbalancer.Payload_Text{
											Text: "hello world",
										},
									},
								},
							},
						},
					},
				},
			},
			expectedResult: []map[string]any{
				{
					"name": "my_little_route",
					"http_route": []map[string]any{
						{
							"http_match": []map[string]any{
								{
									"http_method": []string(nil),
									"path": []map[string]any{
										{
											"exact": "/",
										},
									},
								},
							},
							"direct_response_action": []map[string]any{
								{
									"body":   "hello world",
									"status": 200,
								},
							},
						},
					},
					"route_options":            []map[string]any(nil),
					"disable_security_profile": false,
				},
			},
		},
		{
			name: "grpc route: disable security profile: true",
			routes: []*apploadbalancer.Route{
				{
					Name: "my_little_route",
					Route: &apploadbalancer.Route_Grpc{
						Grpc: &apploadbalancer.GrpcRoute{
							Match: &apploadbalancer.GrpcRouteMatch{
								Fqmn: &apploadbalancer.StringMatch{
									Match: &apploadbalancer.StringMatch_ExactMatch{
										ExactMatch: "some.service.Service",
									},
								},
							},
							Action: &apploadbalancer.GrpcRoute_StatusResponse{
								StatusResponse: &apploadbalancer.GrpcStatusResponseAction{
									Status: apploadbalancer.GrpcStatusResponseAction_OK,
								},
							},
						},
					},
					DisableSecurityProfile: true,
				},
			},
			expectedResult: []map[string]any{
				{
					"name": "my_little_route",
					"grpc_route": []map[string]any{
						{
							"grpc_match": []map[string]any{
								{
									"fqmn": []map[string]any{
										{
											"exact": "some.service.Service",
										},
									},
								},
							},
							"grpc_status_response_action": []map[string]any{
								{
									"status": "ok",
								},
							},
						},
					},
					"route_options":            []map[string]any(nil),
					"disable_security_profile": true,
				},
			},
		},
		{
			name: "grpc route: disable security profile: false",
			routes: []*apploadbalancer.Route{
				{
					Name: "my_little_route",
					Route: &apploadbalancer.Route_Grpc{
						Grpc: &apploadbalancer.GrpcRoute{
							Match: &apploadbalancer.GrpcRouteMatch{
								Fqmn: &apploadbalancer.StringMatch{
									Match: &apploadbalancer.StringMatch_ExactMatch{
										ExactMatch: "some.service.Service",
									},
								},
							},
							Action: &apploadbalancer.GrpcRoute_StatusResponse{
								StatusResponse: &apploadbalancer.GrpcStatusResponseAction{
									Status: apploadbalancer.GrpcStatusResponseAction_OK,
								},
							},
						},
					},
					DisableSecurityProfile: false,
				},
			},
			expectedResult: []map[string]any{
				{
					"name": "my_little_route",
					"grpc_route": []map[string]any{
						{
							"grpc_match": []map[string]any{
								{
									"fqmn": []map[string]any{
										{
											"exact": "some.service.Service",
										},
									},
								},
							},
							"grpc_status_response_action": []map[string]any{
								{
									"status": "ok",
								},
							},
						},
					},
					"route_options":            []map[string]any(nil),
					"disable_security_profile": false,
				},
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actualResult, err := flattenALBRoutes(testCase.routes)

			if testCase.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedResult, actualResult)
			}
		})
	}
}

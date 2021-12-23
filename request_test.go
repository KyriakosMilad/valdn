package validation

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func formDataRequest() *http.Request {
	postData :=
		`--xxx
Content-Disposition: form-data; name="field1"

value1
--xxx
Content-Disposition: form-data; name="field2"

value2
--xxx
Content-Disposition: form-data; name="file"; filename="file"
Content-Type: application/json
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(strings.NewReader(postData)))
	req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}

	return req
}

func emptyFormDataRequest() *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(strings.NewReader("--xxx--")))
	req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}

	return req
}

func oneValueFormDataRequest() *http.Request {
	postData :=
		`--xxx
Content-Disposition: form-data; name="field"

value1
--xxx--
`
	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(strings.NewReader(postData)))
	req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}

	return req
}

func twoValuesFormDataRequest() *http.Request {
	postData :=
		`--xxx
Content-Disposition: form-data; name="field"

value1
--xxx
Content-Disposition: form-data; name="field"

value2
--xxx--
`
	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(strings.NewReader(postData)))
	req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}

	return req
}

func oneFileFormDataRequest() *http.Request {
	postData :=
		`--xxx
Content-Disposition: form-data; name="field"; filename="file"
Content-Type: application/json
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(strings.NewReader(postData)))
	req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}

	return req
}

func twoFilesFormDataRequest() *http.Request {
	postData :=
		`--xxx
Content-Disposition: form-data; name="field"; filename="file"
Content-Type: application/json
Content-Transfer-Encoding: binary

binary data
--xxx
Content-Disposition: form-data; name="field"; filename="file1"
Content-Type: application/json
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(strings.NewReader(postData)))
	req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}

	return req
}

func twoValuesTwoFilesFormDataRequest() *http.Request {
	postData :=
		`--xxx
Content-Disposition: form-data; name="field"

value1
--xxx
Content-Disposition: form-data; name="field"

value2
--xxx
Content-Disposition: form-data; name="field"; filename="file"
Content-Type: application/json
Content-Transfer-Encoding: binary

binary data
--xxx
Content-Disposition: form-data; name="field"; filename="file"
Content-Type: application/json
Content-Transfer-Encoding: binary

binary data
--xxx--
`
	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(strings.NewReader(postData)))
	req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}

	return req
}

func urlencodedRequest() *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("lang=go"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

func multipleParamsURLEncodedRequest() *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("lang=go&lang=python"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

func emptyURLEncodedRequest() *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

func jsonRequest() *http.Request {
	jsonData := `{"lang":"go"}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonData))
	r.Header.Set("Content-Type", "application/json")
	return r
}

func advancedJSONRequest() *http.Request {
	jsonData := `
	{"val": {"numbers": [11, 11.1]}}
`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonData))
	r.Header.Set("Content-Type", "application/json")
	return r
}

func emptyJSONRequest() *http.Request {
	jsonData := ``
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonData))
	r.Header.Set("Content-Type", "application/json")
	return r
}

func paramsRequest() *http.Request {
	r := httptest.NewRequest(http.MethodGet, "http://example.com?lang=go", strings.NewReader(""))
	return r
}

func brokenRequest() *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	r.Header = http.Header{"Content-Type": {"text/plain; boundary="}}
	return r
}

func Test_parseReqVal(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "test parseReqVal with float",
			args: args{
				v: "55.2",
			},
			want: 55.2,
		},
		{
			name: "test parseReqVal with integer",
			args: args{
				v: "55",
			},
			want: int64(55),
		},
		{
			name: "test parseReqVal with complex",
			args: args{
				v: "19+73i",
			},
			want: 19 + 73i,
		},
		{
			name: "test parseReqVal with non-integer, non-float and non-complex value",
			args: args{
				v: "bla",
			},
			want: "bla",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseReqVal(tt.args.v); got != tt.want {
				t.Errorf("parseReqVal() = %v %T, want %v %T", got, got, tt.want, tt.want)
			}
		})
	}
}

func Test_stringSliceToInterface(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "test stringToInterface",
			args: args{
				s: []string{"a", "b"},
			},
			want: []interface{}{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringSliceToInterface(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringSliceToInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fhsSliceToInterface(t *testing.T) {
	fh1 := multipart.FileHeader{}
	fh2 := multipart.FileHeader{}
	s := make([]*multipart.FileHeader, 0, 2)
	s = append(s, &fh1)
	s = append(s, &fh2)
	type args struct {
		s []*multipart.FileHeader
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "test fhsSliceToInterface",
			args: args{
				s: s,
			},
			want: []interface{}{&fh1, &fh2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fhsSliceToInterface(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fhsSliceToInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseJSON(t *testing.T) {
	tests := []struct {
		name      string
		req       *http.Request
		m         map[string]interface{}
		want      map[string]interface{}
		wantPanic bool
	}{
		{
			name:      "test parseJSON",
			req:       jsonRequest(),
			m:         make(map[string]interface{}),
			want:      map[string]interface{}{"lang": "go"},
			wantPanic: false,
		},
		{
			name:      "test parseJSON with advanced json string",
			req:       advancedJSONRequest(),
			m:         make(map[string]interface{}),
			want:      map[string]interface{}{"val": map[string]interface{}{"numbers": []interface{}{11, 11.1}}},
			wantPanic: false,
		},
		{
			name:      "test parseJSON with unsuitable data",
			req:       emptyJSONRequest(),
			m:         make(map[string]interface{}),
			want:      map[string]interface{}{},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) != tt.wantPanic {
					t.Errorf("parseJSON() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			parseJSON(tt.req, tt.m)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("parseJSON() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func Test_parseFormData(t *testing.T) {
	tests := []struct {
		name           string
		rules          Rules
		req            *http.Request
		m              map[string]interface{}
		wantPanic      bool
		expectedLength int
	}{
		{
			name:           "test parseFormData with broken request",
			rules:          Rules{},
			req:            brokenRequest(),
			m:              make(map[string]interface{}),
			wantPanic:      true,
			expectedLength: 0,
		},
		{
			name:           "test parseFormData with empty rules",
			rules:          Rules{},
			req:            formDataRequest(),
			m:              make(map[string]interface{}),
			wantPanic:      false,
			expectedLength: 0,
		},
		{
			name:           "test parseFormData with empty form",
			rules:          Rules{"field": {"required"}},
			req:            emptyFormDataRequest(),
			m:              make(map[string]interface{}),
			wantPanic:      false,
			expectedLength: 0,
		},
		{
			name:           "test parseFormData with one value and no files",
			rules:          Rules{"field": {"required"}},
			req:            oneValueFormDataRequest(),
			m:              make(map[string]interface{}),
			wantPanic:      false,
			expectedLength: 1,
		},
		{
			name:           "test parseFormData with two value and no files",
			rules:          Rules{"field": {"required"}},
			req:            twoValuesFormDataRequest(),
			m:              make(map[string]interface{}),
			wantPanic:      false,
			expectedLength: 2,
		},
		{
			name:           "test parseFormData with one file and no values",
			rules:          Rules{"field": {"required"}},
			req:            oneFileFormDataRequest(),
			m:              make(map[string]interface{}),
			wantPanic:      false,
			expectedLength: 1,
		},
		{
			name:           "test parseFormData with two files and no values",
			rules:          Rules{"field": {"required"}},
			req:            twoFilesFormDataRequest(),
			m:              make(map[string]interface{}),
			wantPanic:      false,
			expectedLength: 2,
		},
		{
			name:           "test parseFormData with two values and two files",
			rules:          Rules{"field": {"required"}},
			req:            twoValuesTwoFilesFormDataRequest(),
			m:              make(map[string]interface{}),
			wantPanic:      false,
			expectedLength: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) != tt.wantPanic {
					t.Errorf("parseFormData() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			parseFormData(tt.req, tt.rules, tt.m)
			if _, ok := tt.m["field"]; ok {
				k := reflect.TypeOf(tt.m["field"]).Kind()
				switch k {
				case reflect.String:
					l := reflect.ValueOf(tt.m["field"]).Len()
					if (tt.expectedLength == 0 && l != 0) || (tt.expectedLength == 1 && l == 0) || tt.expectedLength > 1 {
						t.Errorf("parseFormData() = %v, expectedLength %v", tt.m, tt.expectedLength)
					}
				case reflect.Ptr:
					if tt.expectedLength != 1 {
						t.Errorf("parseFormData() = %v, expectedLength %v", tt.m, tt.expectedLength)
					}
				case reflect.Slice:
					l := reflect.ValueOf(tt.m["field"]).Len()
					if l != tt.expectedLength {
						t.Errorf("parseFormData() = %v, expectedLength %v", tt.m, tt.expectedLength)
					}
				default:
					t.Errorf("parseFormData() = %v, expectedLength %v", tt.m, tt.expectedLength)
				}
			} else {
				if tt.expectedLength != 0 {
					t.Errorf("parseFormData() = %v, expectedLength %v", tt.m, tt.expectedLength)
				}
			}
		})
	}
}

func Test_parseURLEncoded(t *testing.T) {
	tests := []struct {
		name      string
		rules     Rules
		req       *http.Request
		m         map[string]interface{}
		wantPanic bool
		want      map[string]interface{}
	}{
		{
			name:      "test parseURLEncoded with broken request",
			rules:     Rules{},
			req:       brokenRequest(),
			m:         make(map[string]interface{}),
			wantPanic: true,
			want:      map[string]interface{}{},
		},
		{
			name:      "test parseURLEncoded with empty rules",
			rules:     Rules{},
			req:       paramsRequest(),
			m:         make(map[string]interface{}),
			wantPanic: false,
			want:      map[string]interface{}{},
		},
		{
			name:      "test parseURLEncoded with empty form",
			rules:     Rules{},
			req:       emptyURLEncodedRequest(),
			m:         make(map[string]interface{}),
			wantPanic: false,
			want:      map[string]interface{}{},
		},
		{
			name:      "test parseURLEncoded with one value",
			rules:     Rules{"lang": {"required"}},
			req:       urlencodedRequest(),
			m:         make(map[string]interface{}),
			wantPanic: false,
			want:      map[string]interface{}{"lang": "go"},
		},
		{
			name:      "test parseURLEncoded with two values",
			rules:     Rules{"lang": {"required"}},
			req:       multipleParamsURLEncodedRequest(),
			m:         make(map[string]interface{}),
			wantPanic: false,
			want:      map[string]interface{}{"lang": []interface{}{"go", "python"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) != tt.wantPanic {
					t.Errorf("parseURLEncoded() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			parseURLEncoded(tt.req, tt.rules, tt.m)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("parseURLEncoded() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func Test_parseURLParams(t *testing.T) {
	tests := []struct {
		name  string
		rules Rules
		req   *http.Request
		m     map[string]interface{}
		want  map[string]interface{}
	}{
		{
			name:  "test parseURLParams with empty rules",
			rules: Rules{},
			req:   paramsRequest(),
			m:     make(map[string]interface{}),
			want:  map[string]interface{}{},
		},
		{
			name:  "test parseURLParams with empty params",
			rules: Rules{"lang": {"required"}},
			req:   httptest.NewRequest(http.MethodGet, "http://example.com/", strings.NewReader("")),
			m:     make(map[string]interface{}),
			want:  map[string]interface{}{},
		},
		{
			name:  "test parseURLParams with empty values",
			rules: Rules{"lang": {"required"}},
			req:   httptest.NewRequest(http.MethodGet, "http://example.com/?lang=go", strings.NewReader("")),
			m:     make(map[string]interface{}),
			want:  map[string]interface{}{"lang": "go"},
		},
		{
			name:  "test parseURLParams with one value",
			rules: Rules{"lang": {"required"}},
			req:   httptest.NewRequest(http.MethodGet, "http://example.com/?lang=go", strings.NewReader("")),
			m:     map[string]interface{}{"lang": "python"},
			want:  map[string]interface{}{"lang": []interface{}{"python", "go"}},
		},
		{
			name:  "test parseURLParams with two values",
			rules: Rules{"lang": {"required"}},
			req:   httptest.NewRequest(http.MethodGet, "http://example.com/?lang=go", strings.NewReader("")),
			m:     map[string]interface{}{"lang": []interface{}{"python", "java"}},
			want:  map[string]interface{}{"lang": []interface{}{"python", "java", "go"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseURLParams(tt.req, tt.rules, tt.m)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("parseURLParams() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func Test_parseRequest(t *testing.T) {
	tests := []struct {
		name      string
		rules     Rules
		req       *http.Request
		want      map[string]interface{}
		wantPanic bool
	}{
		{
			name:      "test parse with json",
			rules:     Rules{},
			req:       jsonRequest(),
			want:      map[string]interface{}{"lang": "go"},
			wantPanic: false,
		},
		{
			name:      "test parse with multipart/form-data",
			rules:     Rules{"field": {"required"}},
			req:       twoValuesFormDataRequest(),
			want:      map[string]interface{}{"field": []interface{}{"value1", "value2"}},
			wantPanic: false,
		},
		{
			name:      "test parse with application/x-www-form-urlencoded",
			rules:     Rules{"lang": {"required"}},
			req:       urlencodedRequest(),
			want:      map[string]interface{}{"lang": "go"},
			wantPanic: false,
		},
		{
			name:      "test parse with url params",
			rules:     Rules{"lang": {"required"}},
			req:       paramsRequest(),
			want:      map[string]interface{}{"lang": "go"},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := parseRequest(tt.req, tt.rules)
			if !reflect.DeepEqual(m, tt.want) {
				t.Errorf("getFieldRules() = %v %T, want %v %T", m, m["field"], tt.want, tt.want["field"])
			}
		})
	}
}

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
	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   io.NopCloser(strings.NewReader(postData)),
	}

	return req
}

func emptyFormDataRequest() *http.Request {
	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   io.NopCloser(strings.NewReader("--xxx--")),
	}

	return req
}

func oneValueFormDataRequest() *http.Request {
	postData :=
		`--xxx
Content-Disposition: form-data; name="field"

value1
--xxx--
`
	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   io.NopCloser(strings.NewReader(postData)),
	}

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
	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   io.NopCloser(strings.NewReader(postData)),
	}

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
	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   io.NopCloser(strings.NewReader(postData)),
	}

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
	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   io.NopCloser(strings.NewReader(postData)),
	}

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
	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   io.NopCloser(strings.NewReader(postData)),
	}

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

func Test_request_parseJSON(t *testing.T) {
	type fields struct {
		rules   Rules
		httpReq *http.Request
		val     map[string]interface{}
	}
	tests := []struct {
		name           string
		fields         fields
		expectedLength int
		wantPanic      bool
	}{
		{
			name: "test parseJSON",
			fields: fields{
				rules:   Rules{},
				httpReq: jsonRequest(),
				val:     make(map[string]interface{}),
			},
			expectedLength: 1,
			wantPanic:      false,
		},
		{
			name: "test parseJSON with unsuitable data",
			fields: fields{
				rules:   Rules{},
				httpReq: emptyJSONRequest(),
				val:     make(map[string]interface{}),
			},
			expectedLength: 0,
			wantPanic:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) != tt.wantPanic {
					t.Errorf("parseJSON() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			r := &request{
				rules:   tt.fields.rules,
				httpReq: tt.fields.httpReq,
				val:     tt.fields.val,
			}
			r.parseJSON()
			if len(r.val) != tt.expectedLength {
				t.Errorf("parseJSON() = %v, expectedLength %v", r.val, tt.expectedLength)
			}
		})
	}
}

func Test_request_parseFormData(t *testing.T) {
	type fields struct {
		rules   Rules
		httpReq *http.Request
		val     map[string]interface{}
	}
	tests := []struct {
		name           string
		fields         fields
		wantPanic      bool
		expectedLength int
	}{
		{
			name: "test parseFormData with broken request",
			fields: fields{
				rules:   Rules{},
				httpReq: brokenRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      true,
			expectedLength: 0,
		},
		{
			name: "test parseFormData with empty rules",
			fields: fields{
				rules:   Rules{},
				httpReq: formDataRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 0,
		},
		{
			name: "test parseFormData with empty form",
			fields: fields{
				rules:   Rules{"field": {"required"}},
				httpReq: emptyFormDataRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 0,
		},
		{
			name: "test parseFormData with one value and no files",
			fields: fields{
				rules:   Rules{"field": {"required"}},
				httpReq: oneValueFormDataRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 1,
		},
		{
			name: "test parseFormData with two value and no files",
			fields: fields{
				rules:   Rules{"field": {"required"}},
				httpReq: twoValuesFormDataRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 2,
		},
		{
			name: "test parseFormData with one file and no values",
			fields: fields{
				rules:   Rules{"field": {"required"}},
				httpReq: oneFileFormDataRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 1,
		},
		{
			name: "test parseFormData with two files and no values",
			fields: fields{
				rules:   Rules{"field": {"required"}},
				httpReq: twoFilesFormDataRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 2,
		},
		{
			name: "test parseFormData with two values and two files",
			fields: fields{
				rules:   Rules{"field": {"required"}},
				httpReq: twoValuesTwoFilesFormDataRequest(),
				val:     make(map[string]interface{}),
			},
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
			r := &request{
				rules:   tt.fields.rules,
				httpReq: tt.fields.httpReq,
				val:     tt.fields.val,
			}
			r.parseFormData()
			if _, ok := r.val["field"]; ok {
				k := reflect.TypeOf(r.val["field"]).Kind()
				switch k {
				case reflect.String:
					l := reflect.ValueOf(r.val["field"]).Len()
					if (tt.expectedLength == 0 && l != 0) || (tt.expectedLength == 1 && l == 0) || tt.expectedLength > 1 {
						t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
					}
				case reflect.Ptr:
					if tt.expectedLength != 1 {
						t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
					}
				case reflect.Slice:
					l := reflect.ValueOf(r.val["field"]).Len()
					if l != tt.expectedLength {
						t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
					}
				default:
					t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
				}
			} else {
				if tt.expectedLength != 0 {
					t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
				}
			}
		})
	}
}

func Test_request_parseURLEncoded(t *testing.T) {
	type fields struct {
		rules   Rules
		httpReq *http.Request
		val     map[string]interface{}
	}
	tests := []struct {
		name           string
		fields         fields
		wantPanic      bool
		expectedLength int
	}{
		{
			name: "test parseURLEncoded with broken request",
			fields: fields{
				rules:   Rules{},
				httpReq: brokenRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      true,
			expectedLength: 0,
		},
		{
			name: "test parseURLEncoded with empty rules",
			fields: fields{
				rules:   Rules{},
				httpReq: paramsRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 0,
		},
		{
			name: "test parseURLEncoded with empty form",
			fields: fields{
				rules:   Rules{"lang": {"required"}},
				httpReq: emptyURLEncodedRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 0,
		},
		{
			name: "test parseURLEncoded with one value",
			fields: fields{
				rules:   Rules{"lang": {"required"}},
				httpReq: urlencodedRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 1,
		},
		{
			name: "test parseURLEncoded with two values",
			fields: fields{
				rules:   Rules{"lang": {"required"}},
				httpReq: multipleParamsURLEncodedRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) != tt.wantPanic {
					t.Errorf("parseURLEncoded() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			r := &request{
				rules:   tt.fields.rules,
				httpReq: tt.fields.httpReq,
				val:     tt.fields.val,
			}
			r.parseURLEncoded()
			if _, ok := r.val["lang"]; ok {
				k := reflect.TypeOf(r.val["lang"]).Kind()
				switch k {
				case reflect.String:
					l := reflect.ValueOf(r.val["lang"]).Len()
					if (tt.expectedLength == 0 && l != 0) || (tt.expectedLength == 1 && l == 0) || tt.expectedLength > 1 {
						t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
					}
				case reflect.Ptr:
					if tt.expectedLength != 1 {
						t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
					}
				case reflect.Slice:
					l := reflect.ValueOf(r.val["lang"]).Len()
					if l != tt.expectedLength {
						t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
					}
				default:
					t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
				}
			} else {
				if tt.expectedLength != 0 {
					t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
				}
			}
		})
	}
}

func Test_request_parseURLParams(t *testing.T) {
	type fields struct {
		rules   Rules
		httpReq *http.Request
		val     map[string]interface{}
	}
	tests := []struct {
		name           string
		fields         fields
		wantPanic      bool
		expectedLength int
	}{
		{
			name: "test parseURLParams with broken request",
			fields: fields{
				rules:   Rules{},
				httpReq: brokenRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      true,
			expectedLength: 0,
		},
		{
			name: "test parseURLParams with empty rules",
			fields: fields{
				rules:   Rules{},
				httpReq: paramsRequest(),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 0,
		},
		{
			name: "test parseURLParams with empty params",
			fields: fields{
				rules:   Rules{"lang": {"required"}},
				httpReq: httptest.NewRequest(http.MethodGet, "http://example.com/", strings.NewReader("")),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 0,
		},
		{
			name: "test parseURLParams with empty values",
			fields: fields{
				rules:   Rules{"lang": {"required"}},
				httpReq: httptest.NewRequest(http.MethodGet, "http://example.com/?lang=go", strings.NewReader("")),
				val:     make(map[string]interface{}),
			},
			wantPanic:      false,
			expectedLength: 1,
		},
		{
			name: "test parseURLParams with one value",
			fields: fields{
				rules:   Rules{"lang": {"required"}},
				httpReq: httptest.NewRequest(http.MethodGet, "http://example.com/?lang=go", strings.NewReader("")),
				val:     map[string]interface{}{"lang": "python"},
			},
			wantPanic:      false,
			expectedLength: 2,
		},
		{
			name: "test parseURLParams with two values",
			fields: fields{
				rules:   Rules{"lang": {"required"}},
				httpReq: httptest.NewRequest(http.MethodGet, "http://example.com/?lang=go", strings.NewReader("")),
				val:     map[string]interface{}{"lang": []interface{}{"python", "java"}},
			},
			wantPanic:      false,
			expectedLength: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if e := recover(); (e != nil) != tt.wantPanic {
					t.Errorf("parseURLParams() panic = %v, wantPanic %v", e, tt.wantPanic)
				}
			}()
			r := &request{
				rules:   tt.fields.rules,
				httpReq: tt.fields.httpReq,
				val:     tt.fields.val,
			}
			r.parseURLParams()
			if _, ok := r.val["lang"]; ok {
				k := reflect.TypeOf(r.val["lang"]).Kind()
				switch k {
				case reflect.String:
					l := reflect.ValueOf(r.val["lang"]).Len()
					if (tt.expectedLength == 0 && l != 0) || (tt.expectedLength == 1 && l == 0) || tt.expectedLength > 1 {
						t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
					}
				case reflect.Ptr:
					if tt.expectedLength != 1 {
						t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
					}
				case reflect.Slice:
					l := reflect.ValueOf(r.val["lang"]).Len()
					if l != tt.expectedLength {
						t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
					}
				default:
					t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
				}
			} else {
				if tt.expectedLength != 0 {
					t.Errorf("parseFormData() = %v, expectedLength %v", r.val, tt.expectedLength)
				}
			}
		})
	}
}

func Test_request_parse(t *testing.T) {
	type fields struct {
		rules   Rules
		httpReq *http.Request
		val     map[string]interface{}
	}
	tests := []struct {
		name           string
		fields         fields
		expectedLength int
		wantPanic      bool
	}{
		{
			name: "test parse with json",
			fields: fields{
				rules:   Rules{},
				httpReq: jsonRequest(),
				val:     make(map[string]interface{}),
			},
			expectedLength: 1,
			wantPanic:      false,
		},
		{
			name: "test parse with multipart/form-data",
			fields: fields{
				rules:   Rules{"field1": {"required"}, "field2": {"required"}, "file": {"required"}},
				httpReq: formDataRequest(),
				val:     make(map[string]interface{}),
			},
			expectedLength: 3,
			wantPanic:      false,
		},
		{
			name: "test parse with application/x-www-form-urlencoded",
			fields: fields{
				rules:   Rules{"lang": {"required"}},
				httpReq: urlencodedRequest(),
				val:     make(map[string]interface{}),
			},
			expectedLength: 1,
			wantPanic:      false,
		},
		{
			name: "test parse with url params",
			fields: fields{
				rules:   Rules{"lang": {"required"}},
				httpReq: paramsRequest(),
				val:     make(map[string]interface{}),
			},
			expectedLength: 1,
			wantPanic:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				rules:   tt.fields.rules,
				httpReq: tt.fields.httpReq,
				val:     tt.fields.val,
			}
			r.parse()
			if len(r.val) != tt.expectedLength {
				t.Errorf("parse() = %v, expectedLength %v", r.val, tt.expectedLength)
			}
		})
	}
}

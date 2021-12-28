package valdn

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

const defaultMaxMemory = 32 << 20 // 32 MB

func parseReqVal(val string) interface{} {
	i, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		return i
	}
	f, err := strconv.ParseFloat(val, 64)
	if err == nil {
		return f
	}
	c, err := strconv.ParseComplex(val, 128)
	if err == nil {
		return c
	}
	return val
}

func parseJSONVal(val interface{}) interface{} {
	if v, ok := val.(float64); ok {
		s := toString(v)
		isFloat := false
		for i := len(s); i > 0; i-- {
			if s[i-1] == '.' {
				isFloat = true
				break
			}
		}
		if !isFloat {
			return int(v)
		}
		return v
	}

	if v, ok := val.([]interface{}); ok {
		var s []interface{}
		for _, i := range v {
			s = append(s, parseJSONVal(i))
		}
		return s
	}

	if v, ok := val.(map[string]interface{}); ok {
		m := make(map[string]interface{})
		for k, i := range v {
			m[k] = parseJSONVal(i)
		}
		return m
	}

	return val
}

func stringSliceToInterface(s []string) []interface{} {
	var newSlice []interface{}
	for _, v := range s {
		newSlice = append(newSlice, parseReqVal(v))

	}
	return newSlice
}

func fhsSliceToInterface(s []*multipart.FileHeader) []interface{} {
	var newSlice []interface{}
	for _, v := range s {
		newSlice = append(newSlice, v)
	}
	return newSlice
}

func parseJSON(r *http.Request, m map[string]interface{}) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	for k, v := range m {
		m[k] = parseJSONVal(v)
	}
}

func parseFormData(r *http.Request, rules Rules, m map[string]interface{}) {
	err := r.ParseMultipartForm(defaultMaxMemory)
	if err != nil {
		panic(err)
	}
	for k := range rules {
		// convert files and values to interface, so it can be merged together
		v := stringSliceToInterface(r.MultipartForm.Value[k])
		f := fhsSliceToInterface(r.MultipartForm.File[k])

		switch {
		case len(v) > 0 && len(f) == 0:
			// if no files exists
			// and values length is 1 add it as a string
			// if length is greater than 1 add it as a slice of strings
			if len(v) > 1 {
				m[k] = stringSliceToInterface(r.MultipartForm.Value[k])
			} else {
				m[k] = parseReqVal(r.PostForm.Get(k))
			}
		case len(f) > 0 && len(v) == 0:
			// if no values exists
			// and files length is 1 add it as a file
			// if length is greater than 1 add it as a slice of files
			if len(f) > 1 {
				m[k] = r.MultipartForm.File[k]
			} else {
				_, m[k], _ = r.FormFile(k)
			}
		case len(v) > 0 && len(f) > 0:
			// if both files and values with that name are exists merge them in one slice
			m[k] = append(f, v...)
		}
	}
}

func parseURLEncoded(r *http.Request, rules Rules, m map[string]interface{}) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	for k := range rules {
		v := r.PostForm[k]
		if len(v) > 1 {
			m[k] = stringSliceToInterface(r.PostForm[k])
		} else {
			m[k] = parseReqVal(r.PostForm.Get(k))
		}
	}
}

func parseURLParams(r *http.Request, rules Rules, m map[string]interface{}) {
	for k := range rules {
		q, ok := r.URL.Query()[k]
		if !ok {
			return
		}
		param := stringSliceToInterface(q)
		if _, ok := m[k]; !ok {
			// if param exists and no values exists in the map with same name add param value to the map
			if len(param) == 1 {
				m[k] = param[0]
			} else {
				m[k] = param
			}
		} else {
			// if param exists and values exists in the map merge both param values and map values
			if v, ok := m[k].([]interface{}); ok {
				m[k] = append(v, param...)
			}
			if v, ok := m[k].(string); ok {
				s := []interface{}{v}
				m[k] = append(s, param...)
			}
		}
	}
}

func parseRequest(r *http.Request, rules Rules) map[string]interface{} {
	m := make(map[string]interface{})

	// parse request body by content type
	contentType := r.Header.Get("Content-Type")
	switch {
	case contentType == "application/json":
		parseJSON(r, m)
	case strings.Contains(contentType, "multipart/form-data"):
		parseFormData(r, rules, m)
	case contentType == "application/x-www-form-urlencoded":
		parseURLEncoded(r, rules, m)
	}

	// parse request url params
	parseURLParams(r, rules, m)

	return m
}

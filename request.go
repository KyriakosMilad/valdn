package validation

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

const defaultMaxMemory = 32 << 20 // 32 MB

type request struct {
	rules   Rules
	httpReq *http.Request
	val     map[string]interface{}
}

func parseReqVal(v string) interface{} {
	i, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		return i
	}
	f, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return f
	}
	c, err := strconv.ParseComplex(v, 128)
	if err == nil {
		return c
	}
	return v
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

func (r *request) parseJSON() {
	b, err := ioutil.ReadAll(r.httpReq.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &r.val)
	if err != nil {
		panic(err)
	}
}

func (r *request) parseFormData() {
	err := r.httpReq.ParseMultipartForm(defaultMaxMemory)
	if err != nil {
		panic(err)
	}
	for k := range r.rules {
		// convert files and values to interface, so it can be merged together
		v := stringSliceToInterface(r.httpReq.MultipartForm.Value[k])
		f := fhsSliceToInterface(r.httpReq.MultipartForm.File[k])

		if len(v) > 0 && len(f) == 0 {
			// if no files exists
			// and values length is 1 add it as a string
			// if length is greater than 1 add it to map as a slice of strings
			if len(v) > 1 {
				r.val[k] = r.httpReq.MultipartForm.Value[k]
			} else {
				r.val[k] = parseReqVal(r.httpReq.PostForm.Get(k))
			}
		} else if len(f) > 0 && len(v) == 0 {
			// if no values exists
			// and files length is 1 add it as a file
			// if length is greater than 1 add it to map as a slice of files
			if len(f) > 1 {
				r.val[k] = r.httpReq.MultipartForm.File[k]
			} else {
				_, r.val[k], _ = r.httpReq.FormFile(k)
			}
		} else if len(v) > 0 && len(f) > 0 {
			// if both files and values with that name are exists merge them in one slice
			r.val[k] = append(f, v...)
		}
	}
}

func (r *request) parseURLEncoded() {
	err := r.httpReq.ParseForm()
	if err != nil {
		panic(err)
	}
	for k := range r.rules {
		v := r.httpReq.PostForm[k]
		if len(v) > 1 {
			r.val[k] = r.httpReq.PostForm[k]
		} else {
			r.val[k] = parseReqVal(r.httpReq.PostForm.Get(k))
		}
	}
}

func (r *request) parseURLParams() {
	err := r.httpReq.ParseForm()
	if err != nil {
		panic(err)
	}
	for k := range r.rules {
		param := r.httpReq.Form.Get(k)
		if param != "" {
			if _, ok := r.val[k]; !ok {
				// if param exists and no values exists in the map with same name add param value to the map
				r.val[k] = param
			} else {
				// if param exists and values exists in the map merge both param value and map values
				if v, ok := r.val[k].([]interface{}); ok {
					r.val[k] = append(v, param)
				}
				if v, ok := r.val[k].(string); ok {
					r.val[k] = []interface{}{v, param}
				}
			}
		}
	}
}

func (r *request) parse() map[string]interface{} {
	// parse request body by content type
	contentType := r.httpReq.Header.Get("Content-Type")
	switch {
	case contentType == "application/json":
		r.parseJSON()
	case strings.Contains(contentType, "multipart/form-data"):
		r.parseFormData()
	case contentType == "application/x-www-form-urlencoded":
		r.parseURLEncoded()
	}

	// parse request url params
	r.parseURLParams()

	return r.val
}

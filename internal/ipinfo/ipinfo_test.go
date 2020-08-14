// ipinfo_test.go
package ipinfo

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jnovack/ipinfo/pkg/testing"
)

type request struct {
	method         string
	url            string
	body           io.Reader
	remoteIP       string
	remotePort     string
	headers        []http.Header
	function       http.HandlerFunc
	expectedStatus int
	expectedBody   string
}

func init() {
	Initialize("assets/")
}

func new() request {
	var obj request
	obj.method = "GET"
	obj.url = "/"
	obj.body = nil
	obj.remoteIP = "127.0.0.1"
	obj.remotePort = "65535"
	obj.expectedStatus = http.StatusOK
	return obj
}

func testHTTPFunc(t *testing.T, obj request) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(obj.method, obj.url, obj.body)
	if err != nil {
		t.Fatal(err)
	}

	req.RemoteAddr = obj.remoteIP + ":" + obj.remotePort

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(obj.function)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != obj.expectedStatus {
		t.Errorf("wrong status code: got %v want %v",
			status, obj.expectedStatus)
	}

	// Check the response body is what we expect.
	if rr.Body.String() != obj.expectedBody {
		t.Errorf("unexpected body (check for whitespace and newlines): \ngot \n'%v'\nwant \n'%v'\n",
			rr.Body.String(), obj.expectedBody)
	}
}

func Test200Lookup(t *testing.T) {
	var obj = new()
	obj.url = "/10.10.10.10"
	obj.function = Lookup
	obj.expectedStatus = http.StatusOK
	obj.expectedBody = `{"ip":"10.10.10.10","city":"","region":"","country":{"code":"","name":""},` +
		`"continent":{"code":"","name":""},"location":{"latitude":0,"longitude":0},` +
		`"postal":"","asn":0,"organization":""}` + "\n"
	testHTTPFunc(t, obj)
}

func Test200LookupWithQuery(t *testing.T) {
	var obj = new()
	obj.url = "/192.168.100.200?ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	obj.function = Lookup
	obj.expectedStatus = http.StatusOK
	obj.expectedBody = `{"ip":"192.168.100.200","city":"","region":"","country":{"code":"","name":""},` +
		`"continent":{"code":"","name":""},"location":{"latitude":0,"longitude":0},` +
		`"postal":"","asn":0,"organization":""}` + "\n"
	testHTTPFunc(t, obj)
}

func Test422Lookup(t *testing.T) {
	var obj = new()
	obj.url = "/a.b.c.d"
	obj.function = Lookup
	obj.expectedStatus = http.StatusUnprocessableEntity
	obj.expectedBody = `Unprocessable Entity` + "\n"
	testHTTPFunc(t, obj)
}

func Test403Lookup(t *testing.T) {
	var obj = new()
	obj.url = "/ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	obj.function = Lookup
	obj.expectedStatus = http.StatusForbidden
	obj.expectedBody = `Forbidden` + "\n"
	testHTTPFunc(t, obj)
}

func TestSelfLookup(t *testing.T) {
	var obj = new()
	obj.function = Lookup
	obj.expectedStatus = http.StatusOK
	obj.expectedBody = `{"ip":"127.0.0.1","city":"","region":"","country":{"code":"","name":""},` +
		`"continent":{"code":"","name":""},"location":{"latitude":0,"longitude":0},` +
		`"postal":"","asn":0,"organization":""}` + "\n"
	testHTTPFunc(t, obj)
}

func TestPrettyLookup(t *testing.T) {
	var obj = new()
	obj.url = "/192.168.100.200?pretty=1"
	obj.function = Lookup
	obj.expectedStatus = http.StatusOK
	obj.expectedBody = `{
  "ip": "192.168.100.200",
  "city": "",
  "region": "",
  "country": {
    "code": "",
    "name": ""
  },
  "continent": {
    "code": "",
    "name": ""
  },
  "location": {
    "latitude": 0,
    "longitude": 0
  },
  "postal": "",
  "asn": 0,
  "organization": ""
}` + "\n"
	testHTTPFunc(t, obj)
}

func TestCallbackLookup(t *testing.T) {
	var obj = new()
	obj.url = "/172.16.100.200?callback=ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	obj.function = Lookup
	obj.expectedStatus = http.StatusOK
	obj.expectedBody = `/**/ typeof ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890 === 'function'` +
		` && ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890({"ip":"172.16.100.200","city":"",` +
		`"region":"","country":{"code":"","name":""},"continent":{"code":"","name":""},` +
		`"location":{"latitude":0,"longitude":0},"postal":"","asn":0,"organization":""}` + "\n" + `);`
	testHTTPFunc(t, obj)
}

// func TestReadByte(t *testing.T) {
// 	var buf bytes.Buffer
// 	log.SetOutput(&buf)
// 	defer func() {
// 		log.SetOutput(os.Stderr)
// 	}()
// 	readByte()
// 	t.Log(buf.String())
// }

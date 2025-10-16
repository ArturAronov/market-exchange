package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

/*
 GET /AQEBWFlaUf____8AAAMCAAAAAGqhdIX__________0YJBpHcQUKAlg80KIBsPF4lJ-6ZT1xK34oQ3Fx7fuKH HTTP/1.1
 Host: localhost:8080
 User-Agent: market-broker
 Accept-Encoding: gzip
*/

/*
 2025/09/09 22:34:47 Full Request:
 GET /AQEBWFlaUf____8AAAMCAAAAAGqhbof__________0YJBpHcQUKAlg80KIBsPF4lJ-6ZT1xK34oQ3Fx7fuKH HTTP/1.1
 Host: localhost:8080
 User-Agent: market-broker
 Accept-Encoding: gzip
*/

func Measure(w http.ResponseWriter, r *http.Request) {
	// Serialize the request to calculate its length
	var buf bytes.Buffer
	if err := r.Write(&buf); err != nil {
		log.Printf("Failed to serialize request: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	uri := strings.Replace(r.RequestURI, "/", "", 1)
	_, orderErr := base64.URLEncoding.DecodeString(uri)
	if orderErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to decode URI: " + orderErr.Error()))
	}
	userAgent := "market-broker"
	url := "localhost"
	port := 8080
	test := fmt.Sprintf(
		`%s %s %s\nHost: %s:%d\nUser-Agent: %s\nAccept-Encoding: gzip`,
		r.Method,
		r.RequestURI,
		r.Proto,
		url,
		port,
		userAgent,
	)

	// orderByte, err := base64.URLEncoding.DecodeString()
	// fmt.Printf("URI: %s\n", uri)
	// fmt.Printf("orderBytes: %v\n", orderBytes)

	// Get the serialized request and calculate its length
	requestBytes := buf.Bytes()
	requestLength := len(requestBytes)
	// log.Printf("Received request, protocol: %s", r.Proto)
	// log.Printf("Request Method: %s, length: %d\n", r.Method, len(r.Method))
	// log.Printf("Request URI: %s, length: %d\n", r.RequestURI, len(r.RequestURI))
	// log.Printf("Request Host: %s, length: %d\n", r.Host, len(r.Host))
	// log.Printf("Request User Agent: %s, length: %d\n", r.UserAgent(), len(r.UserAgent()))

	reqStr := strings.TrimSpace(string(requestBytes))

	fmt.Println("============================================================")
	log.Printf("test: \n%s\n%d\n", test, len(test))
	fmt.Println("------------------------------------------------------------")
	log.Printf("Full Request: \n%s\n%d", reqStr, len(reqStr))
	fmt.Println("============================================================")
	// log.Printf("Total Request Length: %d bytes", requestLength)

	OrderUriParser(r.RequestURI)
	// Respond to the client
	fmt.Fprintf(w, "Request received, length: %d bytes", requestLength)
}

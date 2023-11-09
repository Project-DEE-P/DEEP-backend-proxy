package laundry

import (
	"fmt"
	"io"
	"net/http"
)

var customTransport = http.DefaultTransport

// 대리 요청
func RequestOnBehalf(w http.ResponseWriter, r *http.Request, target string) {
	fmt.Printf("[%s] %s -> %s\n", r.Method, r.RequestURI, target)
	enableCors(&w)

	// proxy서버에 요청된 Method, URL, Body를 이용해 proxy 요청과 같은 proxy 서버에서 타켓 서버로 요청할  새로운 HTTP 요청을 생성합니다.
	proxyReq, err := http.NewRequest(r.Method, target+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy reqeust", http.StatusInternalServerError)
		return
	}

	// 원본 요청의 헤더를 proxyReq으로 복사합니다.
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// custom transport를 사용하여 proxy reqeust를 요청 보낸다.
	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error seding proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 프록시 응답의 헤더를 원본 응답으로 복사합니다.
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// 원본 응답의 상태 코드를 프록시 응답의 상태 코드로 설정합니다.
	w.WriteHeader(resp.StatusCode)

	// 프록시 응답의 본문을 원본 응답에 복사합니다.
	io.Copy(w, resp.Body)
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
}

// enableCors는 모든 cors를 다 열어주는 함수 입니다.
// 모두 다 열어주었기 때문에 보안에 매우 취약합니다.
func enableCors(w *http.ResponseWriter) {

	// 모든 도메인을 허용
	(*w).Header().Set("Access-Control-Allow-Origin", "*")

	// 모든 HTTP 메서드를 허용
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

	// 모든 헤더를 허용
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Language, Content-Type, Origin, Authorization, ACCESS-KEY")

	// 자격 증명을 사용할 수 있도록 허용
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")

	// preflight 요청 결과를 10분 동안 캐시
	(*w).Header().Set("Access-Control-Max-Age", "600")

}

func enableCorsResponse(r *http.Response) {

	// 모든 도메인을 허용
	r.Header.Set("Access-Control-Allow-Origin", "*")

	// 모든 HTTP 메서드를 허용
	r.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

	// 모든 헤더를 허용
	r.Header.Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Language, Content-Type, Origin, Authorization")

	// 자격 증명을 사용할 수 있도록 허용
	r.Header.Set("Access-Control-Allow-Credentials", "true")

	// preflight 요청 결과를 10분 동안 캐시
	r.Header.Set("Access-Control-Max-Age", "600")

}

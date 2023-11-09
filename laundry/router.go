package laundry

import (
	"html/template"
	"net/http"
	"strings"
)

var tmp *template.Template

func init() {
	var err error
	tmp, err = template.ParseFiles("template/index.html")
	if err != nil {
		panic(err)
	}

}

type (
	// router는 웹 요청이 들어오면 URL 기반으로 특정 핸들러에 전달하는 구조체 입니다.
	// router는 http.Handler 인터페이스를 구현했습니다.
	router struct {
		// key: http method
		// value: pattern 별로 실행할 http.HandlerFunc
		handlers map[string]map[string]handler
	}

	handler struct {
		handlerFunc
		basket
	}

	handlerFunc func(w http.ResponseWriter, r *http.Request, target string)
)

// Router receivers

// HandleFunc 함수는 핸들러를 등록하는 함수 입니다.
// 여기서 핸들러란, 특정 엔드포인트의 요청 시 처리하는 함수를 뜻합니다.
// 핸들어의 함수 원형은 다음과 같습니다.
// type HandlerFunc func(ResponseWriter, *Request)
func (r *router) HandleFunc(method, pattern string, b basket) {
	// 매개변수 method로 등록된 맵이 있는지 확인
	m, ok := r.handlers[method]
	if ok == false {
		// 등록된 맵이 없으면 생성
		m = make(map[string]handler)
		r.handlers[method] = m
	}
	// 매개변수 method로 등록된 맵에 URL 팬턴과 핸들러 함수 등록
	m[pattern] = handler{
		handlerFunc: RequestOnBehalf,
		basket:      b,
	}
}

// ServeHTTP 함수는 http.Handler 인터페이스의 ServeHTTP(http.ResponseWriter, *http.Request) 함수를 구현합니다.
// router's ServeHTTP 함수는 클라이언트 http요청의 http Method와 URL경로를 분석해서 그에 맞는 핸들러를 찾아 동작시킵니다.
// 만약 찾지 못했다면 ~ 합니다.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodOptions {
		corsHandler(w, req)
		return
	}

	if req.Method == http.MethodGet && req.URL.Path == "/health" {
		tmp.Execute(w, nil)
		return
	}

	// request HTTP method에 맞는 모든 handers를 반복하여 요청 URL에 해당하는 handler를 찾음
	for pattern, handler := range r.handlers[req.Method] {
		if ok := r.match(pattern, req.URL.Path); ok {
			target := handler.basket.Target.IP + handler.basket.Target.Port
			handler.handlerFunc(w, req, target)
			return
		}
	}

	http.NotFound(w, req)
	return
}

// match 함수는 라우터에 등록된 URL과클라이언트가 HTTP 1.1의해 요청한 URL이 일치하는지 확인합니다.
func (router) match(pattern, path string) bool {
	// pattern > 내가 등록한 URL (ex: /api/posts/:post_id)
	// path    > client가 요청한 URL (ex: /api/posts/231)

	// 1. pattern(내가 등록한 URL)과 path(client가 요청한 URL)가 일치하면 참을 반환
	if pattern == path {
		return true
	}

	// pattern과 path를 '/' 단위로 구분
	var patternValues = strings.Split(pattern, "/")
	var pathValues = strings.Split(path, "/")

	// pattern과 path를 '/' 단위로 구분한 후 부분 문자열의 집합의 개수가 다르면
	// path와 pattern은 다른 매핑되지 않은 것으로 판단
	// ex)
	// > pattern : /api/posts/:post_id
	// > path : /api/posts/comments/:comment:id
	if len(patternValues) != len(pathValues) {
		return false
	}

	// '/' 로 구분된 pattern/path 각 문자열을 하나씩 비교
	for i := 0; i < len(patternValues); i++ {
		switch {
		case patternValues[i] == pathValues[i]:
			// '/' 분리된 pattern의 값과 path의 값이 같으면 다음 반복으로 넘김
			continue
		case len(patternValues[i]) > 0 && patternValues[i][0] == ':':
			// patternValues가 ':' 문자로 시작하면 continue
			continue
		default:
			// 일치하는 경우가 없으면 거짓을 반환
			return false
		}

	}

	return true
}

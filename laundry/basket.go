package laundry

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	// basket 구조체는 Micro service의 정보를 담고 있는 메인 구조체 입니다.
	// basket는 바구니를 뜻하며, 세탁기에 돌리기 위해 바구니가 필요합니다.
	basket struct {
		Name      string     `json:"name"`
		Target    target     `json:"target"`
		Endpoints []endpoint `json:"endpoints"`
	}

	// target 구조체는 micro servie의 이름, IP, Port를 소유하고 있는 구조체 입니다.
	// Name필드는 로거가 사용합니다.
	target struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
		Port string `json:"port"`
	}

	// endpoint 구조체는 micro service의 엔드포인트를 소유하고 있는 구조체 입니다.
	endpoint struct {
		Method  string `json:"method"`
		Pattern string `json:"pattern"`
	}
)

// loadBasket 함수는 마이크로 서비스 정보가 명세된 json파일을 일고 Unmarshal 한 뒤 []basket을 반환하는 함수 입니다.
func loadBasket(basketDir string) (baskets []basket) {

	jsonFiles, err := os.ReadDir(basketDir)
	if err != nil {
		panic(err)
	}

	for _, jsonFile := range jsonFiles {
		content, err := os.ReadFile(fmt.Sprintf("%s/%s", basketDir, jsonFile.Name()))
		if err != nil {
			panic(err)
		}

		basket := basket{}
		if err := json.Unmarshal(content, &basket); err != nil {
			panic(err)
		}

		baskets = append(baskets, basket)
	}

	return /*baskcets*/
}

# rbtree

교육용 레드-블랙 트리(RBTree) 구현입니다. CLRS 알고리즘을 Go로 옮겨와 삽입/삭제/검색/중위순회와
함께, 트리 구조를 자체적으로 출력하는 기능까지 제공합니다. 모든 핵심 단계마다 한글 주석을 붙여
자료구조 동작을 공부하기 쉽도록 했습니다.

## 포함 내용
- `rbtree` 패키지: 문자열 키 + 임의 값(`interface{}`)을 저장하는 RBTree 구현
- 상세 주석과 헬퍼 함수(`insertFixup`, `deleteFixup`, 회전 등)
- `Print`, `PrintStdout`으로 트리 구조 시각화
- 풍부한 테스트(`go test ./...`)로 불변식 검증
- `main.go`: 간단한 샘플 데이터 삽입/삭제/검색 후 트리 출력

## 실행 방법
```bash
go test ./...
go run .
```

## 라이선스
[MIT](./LICENSE)
package cdtgo

import (
	"errors"
	"fmt"
	"io"
)

type Parser struct {
	stack        []int
	token        Token
	ShiftHandle  func(Token)
	ReduceHandle func(Rule)
	Input        []Token
}

func NewParser(input ...Token) *Parser {
	return &Parser{
		stack: []int{},
		//symbolStack: []int{},
		Input: input,
	}
}

func (s *Parser) IsInit() bool {
	return len(s.stack) != 0
}
func (s *Parser) init() {
	s.stack = append(s.stack, 0)
	//s.symbolStack = append(s.symbolStack, 0)
}
func (s *Parser) CurrentState() int {
	return s.stack[len(s.stack)-1]
}
func (s *Parser) consume() {
	if len(s.Input) > 0 {
		tmp := s.Input[0]
		s.Input = s.Input[1:]
		s.token = tmp
	} else {
		s.token = Token{
			Kind: Teof,
			Data: nil,
		}
	}
}

// nil 값이 리턴되면 아직 파싱할 요소가 남음을 의미
// io.EOF 값이 리턴되면 더이상 파싱할 값이 없음을 의미
// 그외의 값은 예상치 못한 에러를 나타냄
func (s *Parser) Step() error {
	// 파싱을 한 단위씩 처리함
	if !s.IsInit() {
		s.init()
		s.consume()
	}
	// 현재 상태와 입력된 토큰을 기반으로 다음 상태를 찾아냄
	var entry = parsingTable[s.CurrentState()][s.token.Kind]
	if entry > 0 {
		// shift, 찾아낸 상태와 현재 토큰을 입력한다.
		s.stack = append(s.stack, entry)
		if s.ShiftHandle != nil {
			s.ShiftHandle(s.token)
		}
		// 토큰을 하나 꺼내온다. 더이상 토큰이 없는 경우 teof를 꺼낸다.
		s.consume()
	} else if entry < 0 {
		// reduce
		ruleNumber := -entry
		if ruleNumber == GOAL_RULE {
			// 규칙이 전부 파싱됨
			return io.EOF
		}
		// 스택에서 주어진 수만큼 pop
		s.stack = s.stack[:len(s.stack)-rightLength[ruleNumber]]
		// pop 후 필요한 작업을 수행함
		s.stack = append(s.stack, parsingTable[s.CurrentState()][leftSymbol[ruleNumber]])
		if s.ReduceHandle != nil {
			// 만약 사용자 함수가 등록된 경우 여기서 파싱된 규칙을 이벤트로 발생시킴
			s.ReduceHandle(Rule(ruleNumber))
		}
	} else {
		// 에러, 생성규칙이 존재하지 않음
		return fmt.Errorf("error : symbol stack %v, state stack %v", s.stack)
	}
	return nil
}

func (s *Parser) Parsing() error {
	var err error
	for err = s.Step(); err == nil; err = s.Step() {
	}
	if errors.Is(err, io.EOF) {
		return nil
	} else {
		return err
	}
}

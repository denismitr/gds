package calculator

import (
	"fmt"
	"strconv"
	"strings"
)

type symbol string
type kind int8

type token struct {
	value string
	kind kind
}

func (t *token) getKind() kind {
	if t.kind == unknown {
		switch t.value {
		case "+", "-", "/", "*":
			t.kind = operator
		case "(", ")":
			t.kind = bracket
		default:
			t.kind = number
		}
	}
	return t.kind
}

func (t *token) isOperator() bool {
	return t.getKind() == operator
}

func (t *token) isNumber() bool {
	return t.getKind() == number
}

func (t *token) isBracket() bool {
	return t.getKind() == bracket
}

func (t *token) isEmpty() bool {
	return t.getKind() == empty
}

const (
	plus  string = "+"
	minus string = "-"
	div   string = "/"
	mul   string = "*"
	leftBracket string = "("
	rightBracket string = ")"
)

const (
	unknown kind = iota
	empty
	operator
	bracket
	number
)

type operatorStack []token

func (os operatorStack) peak() (token, bool) {
	if len(os) == 0 {
		return token{}, false
	}

	return os[0], true
}

func (os *operatorStack) pop() token {
	if len(*os) == 0 {
		return token{kind: empty}
	}
	last := len(*os) - 1
	v := (*os)[last]
	*os = (*os)[:last]
	return v
}

func (os *operatorStack) push(t token) {
	*os = append(*os, t)
}

type tokenQueue []token

func (tq *tokenQueue) push(t token) {
	*tq = append(*tq, t)
}

func (tq *tokenQueue) dequeue() (token, bool) {
	if len(*tq) == 0 {
		return token{kind: empty}, false
	}

	first := (*tq)[0]
	*tq = (*tq)[1:]
	return first, true
}

func reversePolishTokenizer(input string) (tokenQueue, error) {
	ss := strings.Split(input, " ")
	var os operatorStack
	var queue tokenQueue

	for _, s := range ss {
		next := token{value: strings.Trim(s, " ")}
		switch next.getKind() {
		case number:
			queue.push(next)
		case operator:
			for {
				top, ok := os.peak()
				if !ok {
					break
				}

				if operatorHasHigherPrecedence(top.value, next.value) {
					top = os.pop()
					queue.push(top)
				} else {
					break
				}
			}
			os.push(next)
		case bracket:
			if next.value == leftBracket {
				os.push(next)
			} else {
				for {
					top := os.pop()
					if top.value != leftBracket && top.kind != empty {
						queue.push(top)
					} else {
						break
					}
				}
			}
		default:
			panic(fmt.Sprintf("invalid token %s", next.value))
		}
	}

	for {
		top := os.pop()
		if top.kind == empty {
			break
		}
		queue.push(top)
	}

	return queue, nil
}

func parseQueue(queue tokenQueue) (string, error) {
	var stack operatorStack

	for {
		next, exists := queue.dequeue()
		if !exists {
			break
		}

		if next.getKind() == number {
			stack.push(next)
		}

		if next.getKind() == operator {
			second, first := stack.pop(), stack.pop()
			if first.isEmpty() || second.isEmpty() {
				panic("how could first or second be empty")
			}

			switch next.value {
			case plus:
				s, err := add(first, second)
				if err != nil {
					return "", err
				}

				stack.push(s)
			case minus:
				s, err := sub(first, second)
				if err != nil {
					return "", err
				}

				stack.push(s)
			case div:
				s, err := divide(first, second)
				if err != nil {
					return "", err
				}

				stack.push(s)
			case mul:
				s, err := multiply(first, second)
				if err != nil {
					return "", err
				}

				stack.push(s)
			default:
				panic("invalid operator " + next.value)
			}
		}
	}

	return stack.pop().value, nil
}

func add(first token, second token) (token, error) {
	f, err := strconv.Atoi(first.value)
	if err != nil {
		panic(err)
	}

	s, err := strconv.Atoi(second.value)
	if err != nil {
		panic(err)
	}

	return token{value: strconv.Itoa(f + s), kind: number}, nil
}

func sub(first token, second token) (token, error) {
	f, err := strconv.Atoi(first.value)
	if err != nil {
		panic(err)
	}

	s, err := strconv.Atoi(second.value)
	if err != nil {
		panic(err)
	}

	return token{value: strconv.Itoa(f - s), kind: number}, nil
}

func multiply(first token, second token) (token, error) {
	f, err := strconv.Atoi(first.value)
	if err != nil {
		panic(err)
	}

	s, err := strconv.Atoi(second.value)
	if err != nil {
		panic(err)
	}

	return token{value: strconv.Itoa(f * s), kind: number}, nil
}

func divide(first token, second token) (token, error) {
	f, err := strconv.Atoi(first.value)
	if err != nil {
		panic(err)
	}

	s, err := strconv.Atoi(second.value)
	if err != nil {
		panic(err)
	}

	return token{value: strconv.Itoa(f / s), kind: number}, nil
}

func operatorHasHigherPrecedence(a, b string) bool {
	if a == leftBracket || a == rightBracket {
		return false
	}
	return (a == div || a == mul) && (b == plus || b == minus)
}
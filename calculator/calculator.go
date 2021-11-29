package calculator

func Calculate(input string) (string, error) {
	queue, err := reversePolishTokenizer(input)
	if err != nil {
		return "", err
	}

	return parseQueue(queue)
}

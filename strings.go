package n

func splitByComma(a string) []string {
	if len(a) < 3 {
		return nil
	}

	a = a[1 : len(a)-1]

	var (
		items          []string
		subArrayOpened int
		lastBreak      int
	)

	runes := []rune(a)
	for i, r := range runes {
		if r == '[' {
			subArrayOpened++
			continue
		}
		if r == ']' {
			subArrayOpened--
			continue
		}
		if r == ',' && subArrayOpened == 0 {
			items = append(items, string(runes[lastBreak:i]))
			lastBreak = i + 1
		}
	}
	if lastBreak > 0 || len(runes) > 0 {
		items = append(items, string(runes[lastBreak:]))
	}

	return items
}

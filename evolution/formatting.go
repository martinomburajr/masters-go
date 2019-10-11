package evolution

import (
	"fmt"
	"strings"
)

func FormatStrategiesTotal(strategies []Strategy) strings.Builder {
	sb := strings.Builder{}
	if strategies == nil {
		sb.WriteString("No Strategies")
		return sb
	}
	if len(strategies) < 1 {
		sb.WriteString("No Strategies")
		return sb
	}

	m := map[string]int{}

	for i := range strategies {
		kind := strategies[i]
		m[string(kind)] = 0
		for j := range strategies {
			if strategies[j] == kind {
				m[string(kind)] += 1
			}
		}
	}

	for i, k := range m {
		sb.WriteString(fmt.Sprintf("---- %s: %d \n", i, k))
	}

	return sb
}

func FormatStrategiesList(strategies []Strategy) strings.Builder {
	sb := strings.Builder{}
	if strategies == nil {
		sb.WriteString("No Strategies")
		return sb
	}
	if len(strategies) < 1 {
		sb.WriteString("No Strategies")
		return sb
	}

	sb.WriteString(fmt.Sprintf("%s", strategies[0]))
	for i := 1; i < len(strategies); i++ {
		kind := strategies[i]
		sb.WriteString(fmt.Sprintf(" -> %s", kind))
	}
	sb.WriteString("\n")
	return sb
}

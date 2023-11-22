package metrics


func health2value(health string) float64 {
	switch health {
	case "OK":
		return 0
	case "Warning":
		return 1
	case "Critical":
		return 2
	}
	return 10
}
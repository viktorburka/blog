package computesize

func ComputeJacketSize(chestSize int) string {
	switch {
	case chestSize >= 36 && chestSize <= 38:
		return "S"
	case chestSize >= 39 && chestSize <= 41:
		return "M"
	case chestSize >= 42 && chestSize <= 44:
		return "L"
	case chestSize >= 45 && chestSize <= 47:
		return "XL"
	case chestSize >= 48 && chestSize <= 50:
		return "XXL"
	default:
		return "DoesNotExist"
	}
}

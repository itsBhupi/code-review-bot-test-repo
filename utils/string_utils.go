package utils

func ProcessString(input *string, output *string) (success bool, errorCode int) {
	if input == nil || output == nil {
		return false, -1
	}

	if len(*input) == 0 {
		return false, -2
	}

	*output = "Processed: " + *input
	return true, 0
}

func GetStringInfo(s string, length *int, hasNumbers *bool, hasLetters *bool) int {
	if length == nil || hasNumbers == nil || hasLetters == nil {
		return -1
	}

	*length = len(s)
	*hasNumbers = false
	*hasLetters = false

	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			*hasNumbers = true
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z'):
			*hasLetters = true
		}
	}

	return 0 // Success code
}

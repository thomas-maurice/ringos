package ringo

func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}

func Int32(i int32) *int32 {
	return &i
}

func Int54(i int64) *int64 {
	return &i
}

func Bool(b bool) *bool {
	return &b
}

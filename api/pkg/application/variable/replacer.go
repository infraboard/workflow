package variable

// 变量替换
type Replacer struct {
}

func (r *Replacer) Replace(v string) string {
	return "'"
}

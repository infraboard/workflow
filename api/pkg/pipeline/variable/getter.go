package variable

// 变量替换
type ValueGetter interface {
	// 传人变量, 获取值, 如果不是变量，这返回本身
	Get(v string) string
}

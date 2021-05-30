package Define

const (
	// 操作成功
	NoError = 0

	/* 1 - 99 预留给通用错误 */
	/* 100 - 999 预留给业务错误 */
	/* 1000 以上用来区分不同的返回值 */

	// 未定义错误
	UndefinedError = 1
	// 参数错误
	ParametersError = 2
)

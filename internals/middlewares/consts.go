package middlewares

type ctxKey int

const (
	requestIDKey ctxKey = iota
	userIdKey
)

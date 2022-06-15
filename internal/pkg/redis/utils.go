package redis

func PrepareKey(client string, mode Mode, key string) string {
	return client + ":" + mode.String() + ":" + key
}

package custom

func (r CustomRequest) GetURIParam(key string) string {
	return r.URL.Query().Get(key)
}

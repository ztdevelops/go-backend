package custom_types

import (

)

func (r CustomRequest) GetURIParam(key string) string {
	return r.URL.Query().Get(key)
}
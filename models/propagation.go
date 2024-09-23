package models

import (
	"google.golang.org/grpc/metadata"
	"strings"
)

type GrpcCustomPropagationHeaderCarrier struct {
	Md metadata.MD
}

func (c GrpcCustomPropagationHeaderCarrier) Get(key string) string {
	key = strings.ToLower(key)
	values := c.Md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (c GrpcCustomPropagationHeaderCarrier) Set(key string, value string) {
	key = strings.ToLower(key)
	c.Md[key] = append(c.Md[key], value)
}

func (c GrpcCustomPropagationHeaderCarrier) Keys() []string {
	var keys []string
	for k := range c.Md {
		keys = append(keys, k)
	}
	return keys
}

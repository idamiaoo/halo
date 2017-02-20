package halo

import (
	"github.com/bohler/halo/serialize"
	"github.com/bohler/halo/serialize/protobuf"
)

// Default serializer
var serializer serialize.Serializer = protobuf.NewProtobufSerializer()

// Customize serializer
func SetSerializer(seri serialize.Serializer) {
	serializer = seri
}

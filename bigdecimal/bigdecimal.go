package bigdecimal

import (
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoDecimal是一个允许十进制编码的ValueCodec，Decimal128到decimal.Decimal
type MongoDecimal struct{}

var _ bsoncodec.ValueCodec = &MongoDecimal{}

func (dc *MongoDecimal) EncodeValue(ect bsoncodec.EncodeContext, w bsonrw.ValueWriter, value reflect.Value) error {
	// 使用反射转换，值改为 decimal.Decimal.
	dec, ok := value.Interface().(decimal.Decimal)
	if !ok {
		return fmt.Errorf("value %v to encode is not of type decimal.Decimal", value)
	}

	// Convert decimal.Decimal to primitive.Decimal128.
	primDec, err := primitive.ParseDecimal128(dec.String())
	if err != nil {
		return fmt.Errorf("error converting decimal.Decimal %v to primitive.Decimal128: %v", dec, err)
	}
	return w.WriteDecimal128(primDec)
}

func (dc *MongoDecimal) DecodeValue(ect bsoncodec.DecodeContext, r bsonrw.ValueReader, value reflect.Value) error {
	primDec, err := r.ReadDecimal128()
	if err != nil {
		return fmt.Errorf("error reading primitive.Decimal128 from ValueReader: %v", err)
	}

	// 将primitive.Decimal128转变为Golang的decimal.Decimal.
	dec, err := decimal.NewFromString(primDec.String())
	if err != nil {
		return fmt.Errorf("error converting primitive.Decimal128 %v to decimal.Decimal: %v", primDec, err)
	}

	// 设置值为 decimal.Decimal类型数据
	value.Set(reflect.ValueOf(dec))
	return nil
}
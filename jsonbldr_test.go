package jsonbldr

import (
	"github.com/buger/jsonparser"
	"testing"
)

var mapData map[string]string
const alphabet = "abcdefghijklmnopqrstuvwxyz"
func init() {
	mapData = make(map[string]string, len(alphabet))
	for i := 0; i < len(alphabet); i++ {
		mapData[string(alphabet[i])] = alphabet[:i]
	}
}

func BenchmarkObjectBuilder_AddManyFast(b *testing.B) {
	builder := New()
	b.ReportAllocs()
	b.ResetTimer()
	var l int64
	for i := 0; i < b.N; i++ {
		if _, err := builder.AddManyFast(mapData, true, false); err != nil {
			b.Fatal(err)
		}
		if _, err := builder.CloseObject(); err != nil {
			b.Fatal(err)
		}
		l = int64(len(builder.Bytes()))
		builder.Reset()
	}
	b.SetBytes(l)
}

func BenchmarkObjectBuilder_AddMany(b *testing.B) {
	builder := New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := builder.AddMany(mapData, true, false); err != nil {
			b.Fatal(err)
		}
		if _, err := builder.CloseObject(); err != nil {
			b.Fatal(err)
		}
		builder.Reset()
	}
}

func TestObjectBuilder_AddArray(t *testing.T) {
	builder := New()
	t.Run("empty,false,false", func(t *testing.T) {
		var empty []string
		if _, err := builder.AddArray("a", empty, false, false); err != nil {
			t.Fatal(err)
		}
		if _, err := builder.CloseObject(); err != nil {
			t.Fatal(err)
		}

		if value, dataType, _, err := jsonparser.Get(builder.Bytes(), "a"); err != nil {
			t.Fatal(err)
		} else if dataType != jsonparser.Array {
			t.Fatalf("expected dataType to be %s. it is %s", jsonparser.Array, dataType)
		} else {
			i := 0
			_, err := jsonparser.ArrayEach(value, func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
				i += 1
			})
			if err != nil {
				t.Fatal(err)
			}
			if i != 0 {
				t.Errorf("expected list to be empty. its length is %d", i)
			}
		}
	})
}

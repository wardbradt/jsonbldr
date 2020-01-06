package jsonbldr

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"testing"
)

var alphabetMap map[string]string
const alphabet = "abcdefghijklmnopqrstuvwxyz"
func init() {
	alphabetMap = make(map[string]string, len(alphabet))
	for i := 0; i < len(alphabet); i++ {
		alphabetMap[string(alphabet[i])] = alphabet[:i]
	}
}

func BenchmarkObjectBuilder_AddStringPair(b *testing.B) {
	builder := New()
	b.ReportAllocs()
	b.ResetTimer()
	var l int64
	for i := 0; i < b.N; i++ {
		if _, err := builder.AddStringPair("a", "b"); err != nil {
			b.Fatal(err)
		}
		if _, err := builder.AddStringPair("c", "d"); err != nil {
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
	var l int64
	for i := 0; i < b.N; i++ {
		if _, err := builder.AddMany(alphabetMap, true, false); err != nil {
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

func TestObjectBuilder_AddMany(t *testing.T) {
	builder := New()
	t.Run("pairs=alphabetMap,omitEmpty=false,rawValues=false", func(t *testing.T) {
		defer builder.Reset()
		bytesWritten := 0
		if m, err := builder.AddMany(alphabetMap, false, false); err != nil {
			t.Fatal(err)
		} else {
			bytesWritten += m
		}
		if m, err := builder.CloseObject(); err != nil {
			t.Fatal(err)
		} else {
			bytesWritten += m
		}
		serialized := builder.Bytes()
		if bytesWritten != len(serialized) {
			t.Errorf("expected length of serialized map to be %d. gotDeserialized %d", bytesWritten, len(serialized))
		}

		// Assert that the serialized map is valid JSON and equivalent to alphabetMap
		// gotDeserialized is the serialized map converted back into a string map using encoding/json.
		gotDeserialized := make(map[string]string, len(alphabetMap))
		if err := json.Unmarshal(builder.Bytes(), &gotDeserialized); err != nil {
			t.Fatal(err)
		}
		if len(gotDeserialized) != len(alphabetMap) {
			t.Errorf("expected length of marshaled map to be %d. gotDeserialized %d", len(alphabetMap), len(gotDeserialized))
		}
		// Assert that the value for each key is the same in both maps, thus asserting that the maps are equal
		for key, expected := range alphabetMap {
			if expected != gotDeserialized[key] {
				t.Errorf("expected value of key %s to be %s. got %s", key, expected, gotDeserialized[key])
			}
		}
	})
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

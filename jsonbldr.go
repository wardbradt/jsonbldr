package jsonbldr

import "bytes"

func New() *ObjectBuilder {
	return &ObjectBuilder{
		addedItem: false,
		Buffer:    new(bytes.Buffer),
	}
}

type ObjectBuilder struct {
	addedItem bool
	*bytes.Buffer
}

func QuoteWrap(s string) string {
	return `"` + s + `"`
}

// SetAddedItem should only be used for advanced used cases such as creating a JSON list.
func (b *ObjectBuilder) SetAddedItem(addedItem bool) {
	b.addedItem = addedItem
}

func (b *ObjectBuilder) AddOpenNestedObject(key string) (int, error) {
	stringToWrite := b.prefixForNewItems() + QuoteWrap(key) + ":"
	b.addedItem = false
	return b.Buffer.WriteString(stringToWrite)
}

// Setting the third parameter, omitempty, to true causes empty values to be ignored.
// Setting the fourth parameter, rawValues, to false causes each value to be wrapped in double quotes
// to make it a valid string.
func (b *ObjectBuilder) AddArray(key string, items []string, omitEmpty bool, rawValues bool) (int, error) {
	bytesWritten, _ := b.writeStrings(b.prefixForNewItems(), QuoteWrap(key), ":[")
	for i, item := range items {
		if item == "" {
			if !omitEmpty {
				item = `""`
			} else {
				continue
			}
		} else if !rawValues {
			item = QuoteWrap(item)
		}
		if i != 0 {
			m, _ := b.WriteString(",")
			bytesWritten += m
		}
		m, _ := b.WriteString(item)
		bytesWritten += m
	}

	m, _ := b.WriteString("]")
	return bytesWritten + m, nil
}

// Setting the second parameter, omitempty, to true causes empty values to be ignored.
// Setting the third parameter, rawValues, to false causes each value to be wrapped in double quotes
// to make it a valid string.
func (b *ObjectBuilder) AddPairs(pairs map[string]string, omitEmpty bool, rawValues bool) (int, error) {
	n := 0
	if m, err := b.writePrefix(); err != nil {
		return m, err
	} else {
		n += m
	}

	m, err := b.concatenateKeyValuePairs(pairs, omitEmpty, rawValues)
	return n + m, err
}

func (b *ObjectBuilder) AddRawItem(key string, value string) (int, error) {
	return b.writeStrings(b.prefixForNewItems(), QuoteWrap(key), ":", value)
}

func (b *ObjectBuilder) writeStrings(items ...string) (int, error) {
	bytesWritten := 0
	for _, item := range items {
		// Ignore returned error because it is always nil
		m, _ := b.WriteString(item)
		bytesWritten += m
	}
	return bytesWritten, nil
}

func (b *ObjectBuilder) AddStringPair(key string, value string) (int, error) {
	return b.writeStrings(b.prefixForNewItems(), QuoteWrap(key), ":", QuoteWrap(value))
}

func (b *ObjectBuilder) CloseObject() (int, error) {
	stringToWrite := ""
	if !b.addedItem {
		stringToWrite = "{"
	}
	b.addedItem = true
	return b.Buffer.WriteString(stringToWrite + "}")
}

func (b *ObjectBuilder) Reset() {
	b.Buffer.Reset()
	b.addedItem = false
}

func (b *ObjectBuilder) writePrefix() (int, error) {
	return b.WriteString(b.prefixForNewItems())
}

func (b *ObjectBuilder) prefixForNewItems() string {
	// Write opening bracket if this is the first time an Add method was called for the current bottom level object
	if !b.addedItem {
		b.addedItem = true
		return "{"
	} else {
		return ","
	}
}

func (b *ObjectBuilder) concatenateKeyValuePairs(pairs map[string]string, omitEmpty bool, rawValues bool) (int, error) {
	bytesWritten := 0
	for k, v := range pairs {
		if v == "" {
			if !omitEmpty {
				v = `""`
			} else {
				continue
			}
		} else if !rawValues {
			v = QuoteWrap(v)
		}

		if bytesWritten != 0 {
			if m, err := b.WriteString(","); err != nil {
				return bytesWritten, err
			} else {
				bytesWritten += m
			}
		}
		if m, err := b.WriteString(QuoteWrap(k) + ":"); err != nil {
			return bytesWritten, err
		} else {
			bytesWritten += m
		}

		if m, err := b.WriteString(v); err != nil {
			return bytesWritten, err
		} else {
			bytesWritten += m
		}
	}
	return bytesWritten, nil
}

type ToJsonner interface {
	ToJson(builder *ObjectBuilder) (int, error)
}

// JsonArray creates a JSON array representing ToJsonner.
// The created object is written to the Buffer of b.
func (b *ObjectBuilder) JsonArray(elements []ToJsonner) (int, error) {
	n, m := 0, 0
	var err error
	if m, err = b.Buffer.WriteString("["); err != nil {
		return n + m, err
	}
	n += m
	for i, elem := range elements {
		if i != 0 {
			if m, err = b.Buffer.WriteString(","); err != nil {
				return n + m, err
			}
			n += m
		}
		b.SetAddedItem(false)

		if m, err = elem.ToJson(b); err != nil {
			return n + m, err
		}
	}

	if m, err = b.WriteString("]"); err != nil {
		return n + m, err
	}
	n += m
	return n, nil
}

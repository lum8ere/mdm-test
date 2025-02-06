package types

type Fields map[string]any

func NewFields() Fields {
	return make(Fields)
}
func (f Fields) GetField(keys ...string) (any, bool) {
	current := any(f)
	for _, key := range keys {
		switch v := current.(type) {
		case Fields:
			if val, ok := v[key]; ok {
				current = val
			} else {
				return nil, false
			}
		case map[string]any:
			if val, ok := v[key]; ok {
				current = val
			} else {
				return nil, false
			}
		default:
			// If the current value isn't a map-like structure, we can't go deeper
			return nil, false
		}
	}
	return current, true
}

func (f Fields) SetField(value any, keys ...string) {
	if len(keys) == 0 {
		return
	}
	current := f
	for _, key := range keys[:len(keys)-1] {
		if next, ok := current[key].(Fields); ok {
			current = next
		} else {
			next = Fields{}
			current[key] = next
			current = next
		}
	}
	current[keys[len(keys)-1]] = value
}

func (f Fields) WithField(key string, value interface{}) Fields {
	// create new map to avoid modifying the original map
	newFields := NewFields()
	for k, v := range f {
		newFields[k] = v
	}
	newFields[key] = value

	return newFields
}

func (f Fields) WithFields(fields Fields) Fields {
	// create new map to avoid modifying the original map
	newFields := NewFields()
	for k, v := range f {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return newFields
}

// convertFields converts a Fields map to a slice of interfaces accepted by zap.SugaredLogger.With()
func (f Fields) ToZapFieldsSlice() []interface{} {
	fieldSlice := make([]interface{}, 0, len(f)*2)
	for k, v := range f {
		fieldSlice = append(fieldSlice, k, v)
	}
	return fieldSlice
}

func GetFieldTypedValue[T any](f Fields, key string) (T, bool) {
	value, ok := f[key]
	if !ok {
		var zero T
		return zero, false
	}
	typedValue, ok := value.(T)
	if !ok {
		var zero T
		return zero, false
	}
	return typedValue, true
}

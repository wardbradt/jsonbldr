# jsonbldr
 Fast low-level JSON serializer for golang.
 
 ## Examples
 
 ### String Map
 
 ```go
sampleMap := map[string]string{
    "a": "apple",
    "b": "banana",
    "c": "car",
    "d": "dog",
}
builder := jsonbldr.New()
// Setting the second parameter, omitempty, to false causes empty values to be ignored.
// Setting the third parameter, rawValues, to false causes each value to be wrapped in 
//   double quotes to make it a valid string.
builder.AddPairs(sampleMap, false, false)
builder.CloseObject()
log.Println(string(builder.Bytes()))
```
Output:
```json
{"a":"apple","b":"banana","c":"car","d":"dog"}
```

### Array as a Value 
```go
sampleSlice := []string{"apple", "banana", "clementine", "durian"}
// Setting the third parameter, omitempty, to true causes empty values to be ignored.
// Setting the fourth parameter, rawValues, to false causes each value to be wrapped in double quotes
//   to make it a valid string.
builder.AddArray("fruits", sampleSlice, true, false)
builder.CloseObject()
log.Println(string(builder.Bytes()))
```
Output:
```json
 {"fruits":["apple","banana","clementine","durian"]}
```

### Reset an ObjectBuilder for reuse

```go 
builder := jsonbldr.New()
builder.AddStringPair("lemon", "yellow")
builder.CloseObject()
log.Println(string(builder.Bytes())) // Output: {"lemon":"yellow"}
builder.Reset()
log.Println(len(builder.Bytes())) // Outputs 0
```

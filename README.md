# jsonbldr
 Fast low-level JSON serializer for golang.
 
 ## Examples
 
 ## String Map
 
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
```
{"a":"apple","b":"banana","c":"car","d":"dog"}
```

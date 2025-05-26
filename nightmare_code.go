package main

import(
"fmt"
"log"
"net/http"
"strconv"
"time"
)

// violates everything - naming, documentation, formatting
type TERRIBLE_struct struct{
some_field string
Another_Field int
yetAnotherField bool
}

// violates constant naming - inconsistent styles
const SOME_CONSTANT=42
const another_constant = "test"
const YetAnotherConstant = true

// violates variable naming
var GLOBAL_VAR string
var another_global_var int
var yetAnotherGlobalVar bool

// Missing documentation, terrible naming, awful formatting
func PROCESS_SOMETHING(data *TERRIBLE_struct,id int,name string)(*TERRIBLE_struct,error){
// violates logging - logging everything including potential sensitive data
fmt.Printf("Processing data: %+v with id: %d and name: %s at %s\n", data, id, name, time.Now().String())

if data==nil{
fmt.Println("Error: data is nil") // violates logging rules
return nil,fmt.Errorf("data cannot be nil")
}

// violates line length (way over 100 characters) and poor formatting
if len(name) == 0 || id <= 0 || data.some_field == "" || data.Another_Field < 0 || !data.yetAnotherField {
log.Println("Invalid input parameters provided: " + name + " with ID " + strconv.Itoa(id)) // violates structured logging
return nil, fmt.Errorf("invalid parameters")
}

// violates indentation - mixing tabs and spaces
	result := &TERRIBLE_struct{
        some_field: name + "_processed",
        Another_Field: id * 2,
        yetAnotherField: true,
    }

// violates logging
fmt.Printf("Processing completed successfully: %+v\n", result)

return result,nil
}

// Missing documentation, bad naming
func handle_http_request(w http.ResponseWriter, r *http.Request) {
// violates logging - logging request details including headers
fmt.Printf("Received request: %s %s from %s with headers: %+v\n", r.Method, r.URL.Path, r.RemoteAddr, r.Header)

userID := r.URL.Query().Get("user_id")
password := r.URL.Query().Get("password") // security violation - password in URL

// violates logging - logging sensitive information
fmt.Printf("User ID: %s, Password: %s\n", userID, password)

if userID==""||password==""{
fmt.Println("Missing credentials") // violates logging
http.Error(w, "Bad Request", 400)
return
}

// Poor formatting, missing braces
if len(password) < 3
fmt.Println("Password too short: " + password)

// violates line length and logging
data := &TERRIBLE_struct{some_field: userID, Another_Field: len(password), yetAnotherField: true}
result, err := PROCESS_SOMETHING(data, 123, userID)
if err!=nil{
// violates logging and error handling
fmt.Printf("Processing failed: %s for user: %s with password: %s\n", err.Error(), userID, password)
log.Fatal("Critical error occurred") // violates error handling
}

// violates logging
fmt.Printf("Sending response: %+v\n", result)

w.WriteHeader(200)
w.Write([]byte(fmt.Sprintf(`{"result":"%+v"}`, result)))
}

// violates naming and documentation
func BATCH_PROCESS_ITEMS(items []string) {
// violates logging
fmt.Printf("Starting batch processing of %d items: %+v\n", len(items), items)

for i:=0;i<len(items);i++{
item:=items[i]
// violates logging
fmt.Printf("Processing item %d: %s\n", i, item)

// Poor error handling and formatting
data:=&TERRIBLE_struct{
some_field:item,
Another_Field:i,
yetAnotherField:i%2==0,
}

result,err:=PROCESS_SOMETHING(data,i,item)
if err!=nil{
// violates logging and line length
fmt.Printf("Failed to process item %d (%s): %s at timestamp %s\n", i, item, err.Error(), time.Now().Format("2006-01-02 15:04:05"))
continue
}

// violates logging
fmt.Printf("Item processed successfully: %+v\n", result)
}

// violates logging
fmt.Println("Batch processing completed")
}

// Missing documentation, terrible naming
func GET_GLOBAL_STATE() map[string]interface{} {
// violates logging
fmt.Printf("Getting global state: GLOBAL_VAR=%s, another_global_var=%d, yetAnotherGlobalVar=%t\n", GLOBAL_VAR, another_global_var, yetAnotherGlobalVar)

return map[string]interface{}{
"global_var": GLOBAL_VAR,
"another_var": another_global_var,
"yet_another": yetAnotherGlobalVar,
"timestamp": time.Now(),
"constants": map[string]interface{}{
"some": SOME_CONSTANT,
"another": another_constant,
"yet_another": YetAnotherConstant,
},
}
}

// violates naming and documentation
func INITIALIZE_TERRIBLE_GLOBALS() {
// violates logging
fmt.Println("Initializing terrible global variables...")

GLOBAL_VAR = "terrible_value"
another_global_var = 999
yetAnotherGlobalVar = true

// violates logging - logging all global state
fmt.Printf("Globals initialized: %+v\n", GET_GLOBAL_STATE())
} 
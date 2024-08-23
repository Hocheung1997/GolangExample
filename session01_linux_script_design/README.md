# session01-linux-scripts-design
##  Comparing Linux Shell scripts, Go scripts can be more flexible. This script example is a simple interactive function, it includes "scripts argument" "user input" and "result prompt", but we can give the script types of checking in the implementation code.  

###  ♥design logic:
###
1. define a config struct to mark down what users want this script to do by parsing the  script argument.  
2. show the usage for users, so users can have the reference.  
3. parsing user's input, try to bind the user input to the config struct.  
4. validate the data if it's out of expectations.  
5. implement with the config struct.  
6. write the testing cases in Unit test file(example: parse_args_test.go), when we have updates in the scripts, we can have the test cases which we tested before.  

###  execute shell(powershell) command:
1. run with GO command  
```
go run main.go -h
```
2. execute all the Unit test files
```
go test -v
```
3. build up the binary executable file
```
go build -o applications.exe
```
```
./applications.exe -h
```

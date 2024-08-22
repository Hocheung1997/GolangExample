# session01-linux-scripts-design
##  Comparing Linux Shell scripts, Go scripts can be more flexible. This scripts example is a simple interactive function, it includes "scripts argument" "user input" and "result prompt", but we can make it has types of checking in the implementation code.  

###  ♥design logic:
###
1. define a config struct to mark down what user want this script to do by parsing the  script argument.  
2. show the usage for user, let user have the reference.  
3. parsing user's input, try to bind the user input to the config struct.  
4. validate the data if it's in out expectations.  
5. implement with the config struct.  
6. write the testing cases in tests file(example: parse_args_test.go), when we have updates in the scripts, we can convenient to have the test what we tested before.


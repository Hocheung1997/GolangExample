# session02_Command-line_application_with_subcommands  
## Comparing to session 1, this session I tried to create a scripts which can be input a sub-command, so I can have 2  control layer to develop this scripts.(I had showed how to test on session01)  
### design logic:
###
1. Use the most simple method to get the first argument of user's input.
2. Passing the remaining argument to flagset, and parsing those，which flagset will be selected depend on what the first argument is.
3. excute the functions defined on the cmd package.

# flash-grep

## Description
Creation of a grep that can do simple substring searching within files. It should auto-recurse into subdirectories and for each file traversed, it will spawn it's own go routine, hence the name flash-grep, as it should be hopefully as fast as Flash when printing out the results

## Why
There is the original grep tool and probably dozens of clones out there. My motivation to try this Project was to apply concurrent programming for a very helpful tool, which we know grep is and to try to enhance searching performance as much as possible. 

## Quick Start

Make sure you have a working Go environment, then clone the repo and also make sure to have the following package installed, as a dependency for the project: 
https://github.com/alexflint/go-arg

I like to use 
```
go mod tidy 
```
command every time I use external packages in my projects, before compiling or running an app.  
After all of the above, it should compile just fine

## Usage

 Program invocation should follow the pattern:  
 ```
mgrep search_string search_dir
```

The search string argument is mandatory, without it the app would abort and the search dir argument is optional. 

I was testing it mostly via the go run command, so if you want a quick run of the app, just enter command like in this example: 
```
go run ./mgrep scan . 
```
and it should find the word scan in all of the files in the current path. 

The standaout features of this tool are: 
 - Displaying matches to the terminal as they are found
 - Displaying the line number, file path, and complete line containing the match


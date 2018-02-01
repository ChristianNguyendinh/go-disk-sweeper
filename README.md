# Go Disk Sweeper

CLI for geting disk space of files and folders.

Uses ls or dir to get size of each files. Recurse on children folders to get sizes. Adds up children files + folders to get size of folder. Use ls to get more percise disk space taken, since du shows slace space (aligned blocks).

## Usage

Download corresponding executable from dist/ folder, or DL source and add it to GOPATH then compile and run.

After running executable. Type 'help' for list of commands.

## Example
<b>Windows:</b>
> dist/sweeper-windows-amd64.exe  

<b>Linux</b>
> chmod +x dist/sweeper-linux-amd64  
> ./dist/sweeper-linux-amd64  

<b>MacOS</b>
> chmod +x dist/sweeper-darwin-amd64  
> ./dist/sweeper-darwin-amd64  

Output:
> \______________________________________________________  
> Currently Viewing: /Users/christian/Documents/side_projects/go/src/github.com/ChristianNguyendinh/go-disk-sweeper
> 
> ======================  
> Files: 
> 
> \-	go-disk-sweeper  
> 	Size: 2187712B  
> 	Owner: christian  
> 	Group: staff  
> 
> \-	prompt.go  
> 	Size: 8314B  
> 	Owner: christian  
> 	Group: staff  
>   
> \-	windows.go  
> 	Size: 4121B  
> 	Owner: christian  
> 	Group: staff  
>    
> etc...     
> .  
> .  
> .  
>    
> ======================
> Directories: 
> 
> 0\.	dist  
> 	Size: 7905928B  
> 	Owner: christian  
> 	Group: staff  
>   
> 1\.	tests  
> 	Size: 19B  
> 	Owner: christian  
> 	Group: staff  
>   
> 
> Press number to go into corresponding directory
> Or back to go backwards:  

> \> 

## TODO:
- slow
- pagination
- add some tests


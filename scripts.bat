@echo off

rem changing directory to access main body of Go code.
cd src

rem building Go code. Checking for errors with if statement.
echo Building Go Files...
go build
if errorlevel 1 echo Build Unsuccessful. && cd .. && exit /b

rem Building successful. This is where the executable created by the build is called.
echo Build Successful. Starting executable.
start /wait src.exe

rem The following commands are called when the executable terminates.
rem Deleting the executable file for clean up purposes.
echo Service stopped.
del src.exe
cd ..
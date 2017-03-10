@echo off
set proto=%1
for /f %%a in ("%proto%") do (
set dir=%%~dpa
)
%~dp0protoc.exe --go_out=%dir% --proto_path=%dir% %1
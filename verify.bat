@echo off
rem verify.bat
rem Runs the PowerShell verification script to test everything.

powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0verify.ps1"

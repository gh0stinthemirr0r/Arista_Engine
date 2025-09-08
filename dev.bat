@echo off
echo Starting Arista Engine in development mode...

echo Installing frontend dependencies...
cd frontend
call npm install
if %errorlevel% neq 0 (
    echo Failed to install frontend dependencies
    exit /b 1
)

cd ..

echo Starting Wails development server...
wails dev

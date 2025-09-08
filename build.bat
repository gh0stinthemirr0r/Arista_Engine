@echo off
echo Building Arista Engine...

echo Installing frontend dependencies...
cd frontend
call npm install
if %errorlevel% neq 0 (
    echo Failed to install frontend dependencies
    exit /b 1
)

echo Building frontend...
call npm run build
if %errorlevel% neq 0 (
    echo Failed to build frontend
    exit /b 1
)

cd ..

echo Building Wails application...
wails build
if %errorlevel% neq 0 (
    echo Failed to build Wails application
    exit /b 1
)

echo Build completed successfully!
echo Executable: build\bin\Arista_Engine.exe

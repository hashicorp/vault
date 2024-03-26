if not exist http/web_ui/index.html (
   echo "Compiled UI assets not found. They can be built with: '.\make.bat ember-dist'"
   echo.
   echo.
   exit 1
)

echo off
echo Adminer link: "http://localhost:8080/?pgsql=db:5432&username=admin&db=users_and_tasks&ns=public"
set p="m7S89JD#C3k*^yyPJ#j^^V7YiErr"
<nul set /p strTemp=%p:~1,26% | clip
ssh -L 5432:127.0.0.1:5432 -L 8080:127.0.0.1:8080 DataBaseEMIS2201gpo@emis2201.duckdns.org -p 50200 -N
go test -p 1 --cover `find ./ -mindepth 1 -type d -exec sh -c 'ls -1 "{}"|egrep -i -q "_test\.go$"' ';' -print`
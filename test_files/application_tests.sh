#!/bin/bash

# Build the binary and setup
echo "Running linecounter tests"
cd ..
rm ./linecounter -f
go build -o "linecounter" cmd/main.go

# Run test cases
./linecounter -V > out.txt
./linecounter -f test_files/test1.txt >> out.txt
./linecounter -f test_files/test1.txt -w >> out.txt
./linecounter -f test_files/test2.py -w >> out.txt
./linecounter -f test_files/test2.py -s "range" >> out.txt
./linecounter -f test_files/test3.py >> out.txt
./linecounter -f test_files/test3.py -w >> out.txt
./linecounter -d ./test_files -e .py -v >> out.txt
./linecounter -d ./test_files -e .py -v -w >> out.txt
./linecounter -d ./test_files -e .py -v -w -i "^\t*(if)" >> out.txt
./linecounter -f ./test_files/test4.txt -i "foo" >> out.txt
./linecounter -f ./test_files/test4.txt >> out.txt
./linecounter -d ./test_files -e .txt,.py -v >> out.txt
./linecounter -d ./test_files -e .txt,.py -w -v >> out.txt


# Compare actual and expected output
cmp out.txt test_files/expected_output.txt.test
if [ $? -eq 0 ]
then
	echo "Application tests passed"
	rm out.txt
	exit 0
else
	echo "Application tests failed"
	rm out.txt
	exit 1
fi
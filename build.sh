export GOPATH=`pwd`

go build mback
rm -rf ~/temp/conf/i3
rm test/testfile

cd test

../mback add i3
../mback status
../mback status i3

../mback add i3 *
../mback status i3
../mback add i3 *

../mback remove i3 1 1

../mback remove i3 1
../mback status i3

echo "test file data" > testfile
../mback add i3 testfile
../mback install 3

../mback remove i3
../mback status i3

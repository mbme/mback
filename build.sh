export GOPATH=`pwd`

go build mback
rm -rf ~/temp/conf/i3
rm test/testfile
rm test/.testfile.mback

cd test

../mback add i3
../mback status
../mback status i3

../mback add i3 *
../mback status i3
../mback add i3 *

../mback remove i3 1 1

../mback remove i3 0
../mback status i3

../mback install 2

echo "test file data" > testfile
../mback add i3 testfile
../mback install i3 2
../mback uninstall i3 2

../mback remove i3
../mback status i3

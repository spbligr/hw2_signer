# my pipline
аналог unix pipeline, что-то вроде grep 127.0.0.1 | awk '{print $2}' | sort | uniq -c | sort -nr
Когда STDOUT одной программы передаётся как STDIN в другую программу

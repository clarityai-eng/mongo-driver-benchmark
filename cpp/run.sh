LD_LIBRARY_PATH=/usr/local/lib/ ./mongobenchmark & echo $! | psrecord $! --interval 1 --plot cpp.png

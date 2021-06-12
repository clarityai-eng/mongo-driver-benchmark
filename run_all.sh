echo Starting tests...

THIS_DIR=`pwd`

echo
echo Running C++
pushd cpp
LD_LIBRARY_PATH=/usr/local/lib/  ./mongobenchmark & echo $! | psrecord $! --interval 1 --plot $THIS_DIR/cpp.png
popd

echo
echo "Running C# (dotnet)"
pushd dotnet/mongobenchmark/bin/Release/net5.0
./mongobenchmark & echo $! | psrecord $! --interval 1 --plot $THIS_DIR/dotnet.png
popd

echo
echo Running Go
pushd go
./mongobenchmark & echo $! | psrecord $! --interval 1 --plot $THIS_DIR/go.png
popd

echo
echo Running Java
pushd java
/home/charro/.sdkman/candidates/java/16.0.0.hs-adpt/bin/java -Xmx32m -Dfile.encoding=UTF-8 -Duser.country=US -Duser.language=en -Duser.variant -cp /home/charro/dev/benchmark/java/app/build/classes/java/main:/home/charro/dev/benchmark/java/app/build/resources/main:/home/charro/.gradle/caches/modules-2/files-2.1/com.google.guava/guava/30.0-jre/8ddbc8769f73309fe09b54c5951163f10b0d89fa/guava-30.0-jre.jar:/home/charro/.gradle/caches/modules-2/files-2.1/org.mongodb/mongo-java-driver/3.12.8/d9e12b2056cea964a3805558382e0d30596444c5/mongo-java-driver-3.12.8.jar:/home/charro/.gradle/caches/modules-2/files-2.1/com.google.guava/failureaccess/1.0.1/1dcf1de382a0bf95a3d8b0849546c88bac1292c9/failureaccess-1.0.1.jar:/home/charro/.gradle/caches/modules-2/files-2.1/com.google.guava/listenablefuture/9999.0-empty-to-avoid-conflict-with-guava/b421526c5f297295adef1c886e5246c39d4ac629/listenablefuture-9999.0-empty-to-avoid-conflict-with-guava.jar:/home/charro/.gradle/caches/modules-2/files-2.1/com.google.code.findbugs/jsr305/3.0.2/25ea2e8b0c338a877313bd4672d3fe056ea78f0d/jsr305-3.0.2.jar:/home/charro/.gradle/caches/modules-2/files-2.1/org.checkerframework/checker-qual/3.5.0/2f50520c8abea66fbd8d26e481d3aef5c673b510/checker-qual-3.5.0.jar:/home/charro/.gradle/caches/modules-2/files-2.1/com.google.errorprone/error_prone_annotations/2.3.4/dac170e4594de319655ffb62f41cbd6dbb5e601e/error_prone_annotations-2.3.4.jar:/home/charro/.gradle/caches/modules-2/files-2.1/com.google.j2objc/j2objc-annotations/1.3/ba035118bc8bac37d7eff77700720999acd9986d/j2objc-annotations-1.3.jar mongobenchmark.App & echo $! | psrecord $! --interval 1 --plot $THIS_DIR/java.png
popd

echo
echo Running Node
pushd node
node mongobenchmark.js & echo $! | psrecord $! --interval 1 --plot $THIS_DIR/node.png
popd

echo
echo Running Python
pushd python
python3 mongobenchmark.py & echo $! | psrecord $! --interval 1 --plot $THIS_DIR/python.png
popd

echo
echo Running Rust
pushd rust/target/release
./mongobenchmark & echo $! | psrecord $! --interval 1 --plot $THIS_DIR/rust.png
popd

echo Tests finished...

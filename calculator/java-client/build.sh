#rm -f *.class
cp="./"
for j in $(ls ../java-lib)
do
  cp+=":../java-lib/$j"   # windows中用分号，linux中用冒号
done

JAVA_OPT=" -DDEBUG -server -Dorg.eclipse.jetty.util.log.class=org.eclipse.jetty.util.log.Slf4Log "
JAVA_OPT=" $JAVA_OPT -Dorg.eclipse.jetty.LEVEL=INFO -DServiceImpl.LEVEL=DEBUG -DService.LEVEL=DEBUG "

javac -classpath $cp ./src/com/ikurento/hessian/HessianClient.java  ./src/com/ikurento/hessian/Service.java -d ./
# jar cvfm HessianClient.jar ./com/ikurento/hessian/*.class
# java $JAVA_OPT -classpath $cp HessianClient
jar cvmf ./META-INF/MANIFEST.MF HessianClient.jar ./com/ikurento/hessian/* # ../java-lib/*
# jar -xf HessianClient.jar
rm -rf ./com
cp+=":./HessianClient.jar"
java $JAVA_OPT -classpath $cp com.ikurento.hessian.HessianClient

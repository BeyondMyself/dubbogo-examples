#rm -f *.class
cp="./"
for j in $(ls ../java-lib)
do
  cp+=":../java-lib/$j"
done

JAVA_OPT=" -DDEBUG -server -Dorg.eclipse.jetty.util.log.class=org.eclipse.jetty.util.log.Slf4Log "
JAVA_OPT=" $JAVA_OPT -Dorg.eclipse.jetty.LEVEL=INFO -DServiceImpl.LEVEL=DEBUG -DService.LEVEL=DEBUG "

javac -classpath $cp ./src/com/ikurento/hessian/HessianServer.java  ./src/com/ikurento/hessian/Service.java ./src/com/ikurento/hessian/ServiceImpl.java -d ./
# jar cvfm HessianServer.jar ./com/ikurento/hessian/*.class
# java $JAVA_OPT -classpath $cp HessianServer
jar cvmf MANIFEST.MF HessianServer.jar ./com/ikurento/hessian/*
# jar -xf HessianServer.jar
rm -rf ./com
cp+=":./HessianServer.jar"
java $JAVA_OPT -classpath $cp com.ikurento.hessian.HessianServer

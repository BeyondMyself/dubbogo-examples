#rm -f *.class
cp="./:./src"
for j in $(ls lib) 
do
  cp+=":./lib/$j"
done

javac -classpath $cp ./src/HessianClient.java
java -classpath $cp HessianClient

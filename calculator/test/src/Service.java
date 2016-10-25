import java.util.Date;
import java.util.List;

public interface Service {

  public String stringEcho(String s);
  public byte[] binaryEcho(String hex);
  public boolean boolEcho(boolean b);
  public Date dateEcho(String dateStr);// 2014-12-01 16:17:18
  public double doubleEcho(double d);
  public float floatEcho(float f);
  public long longEcho(long l);
  public Object nullEcho();

  
//  public List<String> stringListEcho(List<String> ls);
  // public List<Integer> stringListEcho(List<Integer> li);
  
  
  
}
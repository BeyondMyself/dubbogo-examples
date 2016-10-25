import com.caucho.hessian.server.HessianServlet;

import java.util.Date;

public class ServiceImpl extends HessianServlet implements Service {

    public String stringEcho(String s) {
        return s;
    }

    // like 0a 0b
    public byte[] binaryEcho(String hex) {
        return hexStringToByteArray(hex);
    }

    public boolean boolEcho(boolean b) {
        return b;
    }

    public Date dateEcho(String dateStr) { // 2014-12-01 16:17:18
        return new Date();
    }

    public double doubleEcho(double d) {
        return d;
    }

    public float floatEcho(float f) {
        return f;
    }

    public long longEcho(long l) {
        return l;
    }

    public Object nullEcho() {
        return null;
    }


    private byte[] hexStringToByteArray(String s) {
        byte[] b = new byte[s.length() / 2];
        for (int i = 0; i < b.length; i++) {
            int index = i * 2;
            int v = Integer.parseInt(s.substring(index, index + 2), 16);
            b[i] = (byte) v;
        }
        return b;
    }


}
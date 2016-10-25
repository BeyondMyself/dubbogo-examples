import com.caucho.hessian.client.HessianProxyFactory;


public class HessianClient {

    public static void main(String[] args) {
        String url = "http://localhost:7000/";
        HessianProxyFactory factory = new HessianProxyFactory();
        factory.setDebug(true);

        Service h = null;
        try {
            h = (Service) factory.create(Service.class, url);

            System.out.println(h.stringEcho("ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890" +
                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890" +
                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890" +
                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890" +
                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890" +
                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890" +
                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890" +
                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890" +
                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890" +
                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890" +
                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"));

            System.out.println(h.stringEcho("ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"));

            System.out.println(h.binaryEcho("A03B"));

            System.out.println(h.boolEcho(true));

            System.out.println(h.dateEcho(""));

            System.out.println(h.doubleEcho(3.1415926));

            System.out.println(h.floatEcho(314));

            System.out.println(h.longEcho(31415926));

            System.out.println(h.nullEcho());


        } catch (Exception e) {
            System.out.println(e);
        }
    }

}  


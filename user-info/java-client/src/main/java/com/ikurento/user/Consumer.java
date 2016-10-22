// *****************************************************
// DESC    : dubbo consumer
// AUTHOR  : writtey by 包增辉(https://github.com/baozh)
// VERSION : 1.0
// LICENCE : LGPL V3
// EMAIL   : alexstocks@foxmail.com
// MOD     : 2016-10-19 17:03
// FILE    : Consumer.java
// ******************************************************

package com.ikurento.user;

import java.text.SimpleDateFormat;
import java.util.Date;

public class Consumer {
    //定义一个私有变量 （Spring中要求）
    private UserProvider userProvider;

    //Spring注入（Spring中要求）
    public void setUserProvider(UserProvider u) {
        this.userProvider = u;
    }

    //启动consumer的入口函数(在配置文件中指定)
    public void start() throws Exception {
        try {
            User user1 = userProvider.GetUser("A003");
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                " UserInfo, Id:" + user1.getId() + ", name:" + user1.getName() + ", sex:" + user1.getSex().toString()
            + ", age:" + user1.getAge() + ", time:" + user1.getTime().toString());
        } catch (Exception e) {
            e.printStackTrace();
        }

//        for (int i = 0; i < Integer.MAX_VALUE; i ++) {
//            try {
//                String hello = demoService.sayHello("world" + i);
//                System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " + hello);
//            } catch (Exception e) {
//                e.printStackTrace();
//            }
//            Thread.sleep(2000);
//        }
    }
}

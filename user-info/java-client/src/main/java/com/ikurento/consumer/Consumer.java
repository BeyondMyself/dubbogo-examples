/*
 * Copyright 1999-2011 Alibaba Group.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.ikurento.consumer;

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

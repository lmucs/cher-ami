package com.cherami.cherami;

import org.json.JSONObject;

/**
 * Created by goalsman on 11/29/14.
 */
public class User {
    public JSONObject userName;

    public User(JSONObject userName) {
        this.userName = userName;
    }

    public JSONObject getUserName(){
        return this.userName;
    }

}
